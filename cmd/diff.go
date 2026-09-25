package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	diffNamespace  string
	diffRepository string
	diffTagA       string
	diffTagB       string
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare Quay.io resources",
}

var diffTagCmd = &cobra.Command{
	Use:     "tag",
	Short:   "Compare two tags' labels and security vulnerabilities",
	Example: "  go-quay diff tag -n myorg -r myrepo --tag-a previous --tag-b candidate",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		diff, err := compareTags(cmd.Context(), client, diffNamespace, diffRepository, diffTagA, diffTagB)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Tag diff for %s/%s:%s -> %s\n", diffNamespace, diffRepository, diffTagA, diffTagB)
		return printJSON(diff)
	},
}

type tagDiffSide struct {
	Name           string `json:"name"`
	ManifestDigest string `json:"manifest_digest"`
	SecurityStatus string `json:"security_status"`
}

type tagDiffResult struct {
	TagA                   tagDiffSide         `json:"tag_a"`
	TagB                   tagDiffSide         `json:"tag_b"`
	AddedLabels            []lib.ManifestLabel `json:"added_labels"`
	RemovedLabels          []lib.ManifestLabel `json:"removed_labels"`
	AddedVulnerabilities   []string            `json:"added_vulnerabilities"`
	RemovedVulnerabilities []string            `json:"removed_vulnerabilities"`
}

type labelDiffKey struct {
	key        string
	value      string
	sourceType string
	mediaType  string
}

func compareTags(ctx context.Context, client *lib.Client, namespace, repository, tagAName, tagBName string) (*tagDiffResult, error) {
	tagA, err := client.GetTag(ctx, namespace, repository, tagAName)
	if err != nil {
		return nil, fmt.Errorf("getting tag %q: %w", tagAName, err)
	}
	if tagA.ManifestDigest == "" {
		return nil, fmt.Errorf("tag %q did not return a manifest digest", tagAName)
	}

	tagB, err := client.GetTag(ctx, namespace, repository, tagBName)
	if err != nil {
		return nil, fmt.Errorf("getting tag %q: %w", tagBName, err)
	}
	if tagB.ManifestDigest == "" {
		return nil, fmt.Errorf("tag %q did not return a manifest digest", tagBName)
	}

	labelsA, err := client.GetManifestLabels(ctx, namespace, repository, tagA.ManifestDigest)
	if err != nil {
		return nil, fmt.Errorf("getting labels for tag %q: %w", tagAName, err)
	}
	securityA, err := client.GetManifestSecurity(ctx, namespace, repository, tagA.ManifestDigest, true)
	if err != nil {
		return nil, fmt.Errorf("getting security scan for tag %q: %w", tagAName, err)
	}
	labelsB, err := client.GetManifestLabels(ctx, namespace, repository, tagB.ManifestDigest)
	if err != nil {
		return nil, fmt.Errorf("getting labels for tag %q: %w", tagBName, err)
	}
	securityB, err := client.GetManifestSecurity(ctx, namespace, repository, tagB.ManifestDigest, true)
	if err != nil {
		return nil, fmt.Errorf("getting security scan for tag %q: %w", tagBName, err)
	}

	addedLabels, removedLabels := diffManifestLabels(labelsA.Labels, labelsB.Labels)
	addedVulnerabilities, removedVulnerabilities := diffVulnerabilityIDs(
		vulnerabilityIDs(securityA), vulnerabilityIDs(securityB),
	)

	return &tagDiffResult{
		TagA:                   tagDiffSide{Name: tagAName, ManifestDigest: tagA.ManifestDigest, SecurityStatus: securityA.Status},
		TagB:                   tagDiffSide{Name: tagBName, ManifestDigest: tagB.ManifestDigest, SecurityStatus: securityB.Status},
		AddedLabels:            addedLabels,
		RemovedLabels:          removedLabels,
		AddedVulnerabilities:   addedVulnerabilities,
		RemovedVulnerabilities: removedVulnerabilities,
	}, nil
}

func diffManifestLabels(tagA, tagB []lib.ManifestLabel) ([]lib.ManifestLabel, []lib.ManifestLabel) {
	labelsA := make(map[labelDiffKey]lib.ManifestLabel, len(tagA))
	labelsB := make(map[labelDiffKey]lib.ManifestLabel, len(tagB))
	for _, label := range tagA {
		labelsA[labelDiffKey{label.Key, label.Value, label.SourceType, label.MediaType}] = label
	}
	for _, label := range tagB {
		labelsB[labelDiffKey{label.Key, label.Value, label.SourceType, label.MediaType}] = label
	}

	added := make([]lib.ManifestLabel, 0)
	removed := make([]lib.ManifestLabel, 0)
	for key, label := range labelsB {
		if _, found := labelsA[key]; !found {
			added = append(added, label)
		}
	}
	for key, label := range labelsA {
		if _, found := labelsB[key]; !found {
			removed = append(removed, label)
		}
	}
	sort.Slice(added, func(i, j int) bool { return labelSortKey(added[i]) < labelSortKey(added[j]) })
	sort.Slice(removed, func(i, j int) bool { return labelSortKey(removed[i]) < labelSortKey(removed[j]) })
	return added, removed
}

func labelSortKey(label lib.ManifestLabel) string {
	return strings.Join([]string{label.Key, label.Value, label.SourceType, label.MediaType}, "\x00")
}

func vulnerabilityIDs(scan *lib.SecurityScan) []string {
	ids := make(map[string]struct{})
	if scan != nil && scan.Data != nil && scan.Data.Layer != nil {
		for _, feature := range scan.Data.Layer.Features {
			for _, vulnerability := range feature.Vulnerabilities {
				if id := strings.TrimSpace(vulnerability.Name); id != "" {
					ids[id] = struct{}{}
				}
			}
		}
	}

	result := make([]string, 0, len(ids))
	for id := range ids {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}

func diffVulnerabilityIDs(tagA, tagB []string) ([]string, []string) {
	idsA := make(map[string]struct{}, len(tagA))
	idsB := make(map[string]struct{}, len(tagB))
	for _, id := range tagA {
		idsA[id] = struct{}{}
	}
	for _, id := range tagB {
		idsB[id] = struct{}{}
	}

	added := make([]string, 0)
	removed := make([]string, 0)
	for id := range idsB {
		if _, found := idsA[id]; !found {
			added = append(added, id)
		}
	}
	for id := range idsA {
		if _, found := idsB[id]; !found {
			removed = append(removed, id)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	return added, removed
}

func init() {
	diffCmd.AddCommand(diffTagCmd)
	diffTagCmd.Flags().StringVarP(&diffNamespace, "namespace", "n", appCfg.Namespace, "Repository namespace")
	diffTagCmd.Flags().StringVarP(&diffRepository, "repository", "r", "", "Repository name")
	diffTagCmd.Flags().StringVar(&diffTagA, "tag-a", "", "Starting tag to compare")
	diffTagCmd.Flags().StringVar(&diffTagB, "tag-b", "", "Tag to compare against tag A")
	if appCfg.Namespace == "" {
		_ = diffTagCmd.MarkFlagRequired("namespace")
	}
	_ = diffTagCmd.MarkFlagRequired("repository")
	_ = diffTagCmd.MarkFlagRequired("tag-a")
	_ = diffTagCmd.MarkFlagRequired("tag-b")
}
