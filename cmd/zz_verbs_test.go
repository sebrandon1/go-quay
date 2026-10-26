package cmd

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const testTrueValue = "true"

func subcommand(parent *cobra.Command, name string) *cobra.Command {
	for _, c := range parent.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

func TestVerbParentsRegistered(t *testing.T) {
	for _, name := range []string{cmdCreate, cmdDelete, cmdUpdate, cmdList, cmdInfo} {
		if subcommand(rootCmd, name) == nil {
			t.Errorf("root command missing %q", name)
		}
	}
}

func TestVerbLeafIsDistinctCommand(t *testing.T) {
	leaf := subcommand(createCmd, cmdRepository)
	if leaf == nil {
		t.Fatal("create repository not registered")
	}
	if leaf == repoCreateCmd {
		t.Fatal("verb leaf must not reuse the get-tree command pointer")
	}
}

func TestCreateRepositoryHasResourceFlags(t *testing.T) {
	leaf := subcommand(createCmd, cmdRepository)
	if leaf == nil {
		t.Fatal("create repository not registered")
	}
	for _, name := range []string{testNamespaceFlagName, cmdRepository, "visibility", "description"} {
		if leaf.Flags().Lookup(name) == nil {
			t.Errorf("create repository missing --%s flag", name)
		}
	}
}

func TestListRepositoryDoesNotRequireRepositoryFlag(t *testing.T) {
	leaf := subcommand(listCmd, cmdRepository)
	if leaf == nil {
		t.Fatal("list repository not registered")
	}
	if leaf.Flags().Lookup("namespace") == nil {
		t.Fatal("list repository missing --namespace")
	}
	f := leaf.Flags().Lookup("repository")
	if f == nil {
		return
	}
	if _, required := f.Annotations[cobra.BashCompOneRequiredFlag]; required {
		t.Fatal("list repository must not mark --repository required")
	}
}

func TestGetMutationsAreDeprecated(t *testing.T) {
	if repoCreateCmd.Deprecated == "" {
		t.Fatal("get repository create should be deprecated")
	}
	if !strings.Contains(repoCreateCmd.Deprecated, cliName+" "+cmdCreate+" "+cmdRepository) {
		t.Errorf("Deprecated = %q, want pointer to create repository", repoCreateCmd.Deprecated)
	}
	if repoListCmd.Deprecated != "" {
		t.Errorf("get repository list should not be deprecated, got %q", repoListCmd.Deprecated)
	}
	leaf := subcommand(createCmd, cmdRepository)
	if leaf != nil && leaf.Deprecated != "" {
		t.Errorf("verb-first create repository should not be deprecated, got %q", leaf.Deprecated)
	}
}

func TestCopyFlagSetSkipsHelpAndExisting(t *testing.T) {
	src := pflag.NewFlagSet("src", pflag.ContinueOnError)
	src.String("help", "", "help")
	src.String("foo", "from-src", "foo")
	src.String("bar", "bar-val", "bar")

	dst := pflag.NewFlagSet("dst", pflag.ContinueOnError)
	dst.String("foo", "keep-me", "foo")

	copyFlagSet(dst, src)

	if dst.Lookup("help") != nil {
		t.Fatal("copyFlagSet should skip help")
	}
	if got := dst.Lookup("foo").DefValue; got != "keep-me" {
		t.Errorf("existing foo flag overwritten, DefValue = %q", got)
	}
	if dst.Lookup("bar") == nil {
		t.Fatal("expected bar to be copied")
	}
}

func TestCopyOneFlagClonesAnnotations(t *testing.T) {
	src := pflag.NewFlagSet("src", pflag.ContinueOnError)
	src.String("ns", "", "namespace")
	if err := src.SetAnnotation("ns", cobra.BashCompOneRequiredFlag, []string{testTrueValue}); err != nil {
		t.Fatal(err)
	}

	dst := pflag.NewFlagSet("dst", pflag.ContinueOnError)
	copyOneFlag(dst, src.Lookup("ns"))

	srcAnn := src.Lookup("ns").Annotations[cobra.BashCompOneRequiredFlag]
	srcAnn[0] = "mutated"
	got := dst.Lookup("ns").Annotations[cobra.BashCompOneRequiredFlag]
	if len(got) != 1 || got[0] != testTrueValue {
		t.Errorf("dst annotations = %v, want [true] (independent of src)", got)
	}
}

func TestVerbLeafNilFlagSet(t *testing.T) {
	dst := pflag.NewFlagSet("dst", pflag.ContinueOnError)
	copyFlagSet(dst, nil)
	var n int
	dst.VisitAll(func(*pflag.Flag) { n++ })
	if n != 0 {
		t.Errorf("copyFlagSet(nil) added %d flags, want 0", n)
	}
}

// TestVerbLeafCopiesPreRunE compares against the source rather than asserting a
// non-nil hook, so it keeps passing if the source command later gains or drops
// its PreRunE.
func TestVerbLeafCopiesPreRunE(t *testing.T) {
	tests := []struct {
		parent *cobra.Command
		use    string
		src    *cobra.Command
	}{
		{deleteCmd, cmdTag, tagDeleteCmd},
		{updateCmd, cmdTag, tagUpdateCmd},
		{infoCmd, cmdTag, tagInfoCmd},
		{infoCmd, cmdRepository, repoInfoCmd},
	}

	for _, test := range tests {
		leaf := subcommand(test.parent, test.use)
		if leaf == nil {
			t.Fatalf("%s %s not registered", test.parent.Use, test.use)
		}
		if funcPointer(leaf.PreRunE) != funcPointer(test.src.PreRunE) {
			t.Errorf("%s %s: PreRunE differs from its source, so local validation behaves differently on the verb-first path", test.parent.Use, leaf.Use)
		}
	}

	if funcPointer(tagInfoCmd.PreRunE) == 0 {
		t.Fatal("tag info no longer has a PreRunE; this test no longer proves hook copying")
	}
}

func funcPointer(fn any) uintptr {
	if reflect.ValueOf(fn).IsNil() {
		return 0
	}
	return reflect.ValueOf(fn).Pointer()
}

// TestVerbLeafExamplesMatchCommandPath guards the example prefix that addVerbs
// builds from the parent verb and the leaf name.
func TestVerbLeafExamplesMatchCommandPath(t *testing.T) {
	var checked int
	for _, parent := range []*cobra.Command{createCmd, deleteCmd, updateCmd, listCmd, infoCmd} {
		for _, leaf := range parent.Commands() {
			if leaf.Example == "" {
				continue
			}
			checked++
			want := "  " + cliName + " " + parent.Use + " " + leaf.Use + " "
			if !strings.HasPrefix(leaf.Example, want) {
				t.Errorf("%s %s example = %q, want prefix %q", parent.Use, leaf.Use, leaf.Example, want)
			}
			if strings.Contains(leaf.Example, cliName+" "+cmdGet) {
				t.Errorf("%s %s example advertises the legacy get path: %q", parent.Use, leaf.Use, leaf.Example)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no verb leaf carries an example; the check proved nothing")
	}
}

// TestVerbLeafUsageShowsExample inspects UsageString rather than executing
// --help: Cobra's help flag keeps its parsed value, so a command that once saw
// --help returns ErrHelp on every later execution in the same process.
func TestVerbLeafUsageShowsExample(t *testing.T) {
	tests := []struct {
		verb     string
		resource string
		want     string
	}{
		{cmdInfo, cmdTag, `go-quay info tag --namespace`},
		{cmdInfo, cmdRepository, `go-quay info repository --namespace`},
	}
	for _, test := range tests {
		t.Run(test.verb+" "+test.resource, func(t *testing.T) {
			leaf := subcommand(subcommand(rootCmd, test.verb), test.resource)
			if leaf == nil {
				t.Fatalf("%s %s not registered", test.verb, test.resource)
			}
			got := leaf.UsageString()
			if !strings.Contains(got, "Examples:") || !strings.Contains(got, test.want) {
				t.Errorf("usage missing %q:\n%s", test.want, got)
			}
			if strings.Contains(got, cliName+" "+cmdGet) {
				t.Errorf("verb-first usage advertises the legacy get path:\n%s", got)
			}
		})
	}
}

func TestVerbLeafValidatesTagBeforeBuildingClient(t *testing.T) {
	resetRootFlags(t)
	rootCmd.SetArgs([]string{
		cmdInfo, cmdTag, "-n", "ns", "-r", "repo",
		testTokenFlag, testTokenValue, testQuayURLFlag, testUnreachableURL,
	})
	rootCmd.SetErr(io.Discard)

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected tag validation to fail")
	}
	if err.Error() != "tag is required; provide --tag" {
		t.Errorf("error = %q, want the local validation %q (not a client error)", err.Error(), "tag is required; provide --tag")
	}
}

func TestVerbLeafDoesNotCopyRootFlags(t *testing.T) {
	leaf := subcommand(createCmd, cmdRepository)
	if leaf == nil {
		t.Fatal("create repository not registered")
	}
	for _, name := range []string{"token", "quay-url", "output"} {
		f := leaf.Flags().Lookup(name)
		rootF := rootCmd.PersistentFlags().Lookup(name)
		if f != nil && f != rootF {
			t.Errorf("verb leaf copied root --%s as a distinct flag", name)
		}
		if leaf.PersistentFlags().Lookup(name) != nil {
			t.Errorf("verb leaf persistent flags should not include --%s", name)
		}
	}
}
