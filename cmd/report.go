package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

const (
	reportDefaultRepositories = 10
	reportDefaultLimit        = 20
	// reportAPIDateLayout is the date format the Quay log API expects.
	reportAPIDateLayout = "01/02/2006"
)

// resolveReportDate accepts an explicit date or a relative token such as
// "30d", which resolves to that many days before today.
func resolveReportDate(value string, now time.Time) (string, error) {
	if !strings.HasSuffix(value, "d") {
		return value, nil
	}
	days, err := strconv.Atoi(strings.TrimSuffix(value, "d"))
	if err != nil || days < 0 {
		return "", fmt.Errorf("invalid relative date %q: use a number of days like 30d", value)
	}
	return now.UTC().AddDate(0, 0, -days).Format(reportAPIDateLayout), nil
}

// reportWindowLabel renders the requested range, noting the range Quay actually
// covered when retention was shorter.
func reportWindowLabel(start, end, coveredStart, coveredEnd string) string {
	if coveredStart == "" {
		return fmt.Sprintf("%s to %s", start, end)
	}
	return fmt.Sprintf("%s to %s, requested %s to %s", coveredStart, coveredEnd, start, end)
}

// defaultReportWindow returns the most recent span Quay retains, formatted the
// way the API expects.
func defaultReportWindow(now time.Time) (string, string) {
	end := now.UTC()
	start := end.AddDate(0, 0, -pullTrafficRetentionDays+1)
	return start.Format(reportAPIDateLayout), end.Format(reportAPIDateLayout)
}

var (
	reportNamespace string
	reportStartDate string
	reportEndDate   string
	reportRepos     int
	reportLimit     int
	reportMaxPages  int
)

// reportConsumer is one puller aggregated across every sampled repository.
type reportConsumer struct {
	Puller       string `json:"puller"`
	IP           string `json:"ip,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	Pulls        int    `json:"pulls"`
	Repositories int    `json:"repositories"`
	LastSeen     string `json:"last_seen,omitempty"`
}

// reportConsumersReport is the machine-readable form of `report consumers`.
type reportConsumersReport struct {
	Namespace           string           `json:"namespace"`
	StartDate           string           `json:"start_date"`
	EndDate             string           `json:"end_date"`
	CoveredStart        string           `json:"covered_start,omitempty"`
	CoveredEnd          string           `json:"covered_end,omitempty"`
	RepositoriesSampled int              `json:"repositories_sampled"`
	TotalPulls          int              `json:"total_pulls"`
	Consumers           []reportConsumer `json:"consumers"`
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Reporting and analytics commands",
}

var reportConsumersCmd = &cobra.Command{
	Use:   "consumers",
	Short: "Show who is pulling images in a namespace",
	Long: `Summarize pull traffic by puller identity and source IP.

Individual log entries are sampled from the busiest repositories (--repos) and
grouped by puller and IP, which is what distinguishes one cluster pulling in a
loop from many clusters.

Without --startdate/--enddate the last 30 days are used, matching Quay's log
retention. Dates accept MM/DD/YYYY, YYYY-MM-DD, or a relative value like 7d.`,
	Example: `  go-quay report consumers -n myorg --token "$QUAY_TOKEN"
  go-quay report consumers -n myorg --repos 5 --limit 10 -O csv --token "$QUAY_TOKEN"`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		// A consumer summary is meant to be read; JSON stays one flag away.
		if !cmd.Flags().Changed(outputFlagName) {
			outputFormat = outputTable
		}

		now := time.Now()
		start, err := resolveReportDate(reportStartDate, now)
		if err != nil {
			return err
		}
		end, err := resolveReportDate(reportEndDate, now)
		if err != nil {
			return err
		}
		if start == "" && end == "" {
			start, end = defaultReportWindow(now)
		}

		report, err := collectReportConsumers(cmd.Context(), client, reportNamespace, start, end,
			reportRepos, reportMaxPages, reportLimit)
		if err != nil {
			return err
		}
		return printReportConsumers(report)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.AddCommand(reportConsumersCmd)

	reportConsumersCmd.Flags().StringVarP(&reportNamespace, "namespace", "n", appCfg.Namespace, "Repository namespace (default: config file)")
	reportConsumersCmd.Flags().StringVarP(&reportStartDate, "startdate", "s", "", "Start date (MM/DD/YYYY, YYYY-MM-DD, or 30d)")
	reportConsumersCmd.Flags().StringVarP(&reportEndDate, "enddate", "e", "", "End date (MM/DD/YYYY, YYYY-MM-DD, or 30d)")
	reportConsumersCmd.Flags().IntVar(&reportRepos, "repos", reportDefaultRepositories, "Number of busiest repositories to sample")
	reportConsumersCmd.Flags().IntVar(&reportLimit, "limit", reportDefaultLimit, "Maximum consumers to show (0 shows all)")
	reportConsumersCmd.Flags().IntVar(&reportMaxPages, "max-pages", 25, "Maximum log pages per repository")
}

// collectReportConsumers samples the busiest repositories and merges the same
// puller and IP across repositories into one row.
func collectReportConsumers(ctx context.Context, client *lib.Client, namespace, start, end string,
	repositories, maxPages, limit int) (*reportConsumersReport, error) {
	base, err := collectPullTraffic(ctx, client, pullTrafficOptions{
		Namespace: namespace,
		StartDate: start,
		EndDate:   end,
		Top:       repositories,
		MaxPages:  maxPages,
		Mode:      modeConsumers,
	})
	if err != nil {
		return nil, err
	}

	type consumerKey struct{ puller, ip string }
	merged := map[consumerKey]*reportConsumer{}
	for _, repository := range base.Repositories {
		for _, consumer := range repository.Consumers {
			key := consumerKey{puller: consumer.Puller, ip: consumer.IP}
			entry, ok := merged[key]
			if !ok {
				entry = &reportConsumer{
					Puller:    consumer.Puller,
					IP:        consumer.IP,
					UserAgent: consumer.UserAgent,
				}
				merged[key] = entry
			}
			entry.Pulls += consumer.Pulls
			entry.Repositories++
			// Last seen values use a sortable fixed-width layout.
			if consumer.LastSeen > entry.LastSeen {
				entry.LastSeen = consumer.LastSeen
			}
		}
	}

	consumers := make([]reportConsumer, 0, len(merged))
	for _, entry := range merged {
		consumers = append(consumers, *entry)
	}
	sort.Slice(consumers, func(i, j int) bool {
		if consumers[i].Pulls != consumers[j].Pulls {
			return consumers[i].Pulls > consumers[j].Pulls
		}
		if consumers[i].Puller != consumers[j].Puller {
			return consumers[i].Puller < consumers[j].Puller
		}
		return consumers[i].IP < consumers[j].IP
	})
	if limit > 0 && len(consumers) > limit {
		consumers = consumers[:limit]
	}

	return &reportConsumersReport{
		Namespace:           base.Namespace,
		StartDate:           base.StartDate,
		EndDate:             base.EndDate,
		CoveredStart:        base.CoveredStart,
		CoveredEnd:          base.CoveredEnd,
		RepositoriesSampled: len(base.Repositories),
		TotalPulls:          base.TotalPulls,
		Consumers:           consumers,
	}, nil
}

func printReportConsumers(report *reportConsumersReport) error {
	if outputFormat != outputTable && outputFormat != outputCSV {
		return printJSON(report)
	}

	headers := []string{"PULLER", "IP", "PULLS", "REPOS", "LAST SEEN"}
	rows := make([][]string, 0, len(report.Consumers))
	for _, consumer := range report.Consumers {
		rows = append(rows, []string{
			consumer.Puller,
			consumer.IP,
			strconv.Itoa(consumer.Pulls),
			strconv.Itoa(consumer.Repositories),
			consumer.LastSeen,
		})
	}

	if outputFormat == outputTable {
		fmt.Fprintf(os.Stderr, "Consumers of '%s' (%s) across %d sampled repositories: %d pulls\n",
			report.Namespace, reportWindowLabel(report.StartDate, report.EndDate, report.CoveredStart, report.CoveredEnd),
			report.RepositoriesSampled, report.TotalPulls)
	}
	return printRows(headers, rows)
}
