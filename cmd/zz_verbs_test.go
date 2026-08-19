package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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
	for _, name := range []string{"namespace", "repository", "visibility", "description"} {
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
	if err := src.SetAnnotation("ns", cobra.BashCompOneRequiredFlag, []string{"true"}); err != nil {
		t.Fatal(err)
	}

	dst := pflag.NewFlagSet("dst", pflag.ContinueOnError)
	copyOneFlag(dst, src.Lookup("ns"))

	srcAnn := src.Lookup("ns").Annotations[cobra.BashCompOneRequiredFlag]
	srcAnn[0] = "mutated"
	got := dst.Lookup("ns").Annotations[cobra.BashCompOneRequiredFlag]
	if len(got) != 1 || got[0] != "true" {
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
