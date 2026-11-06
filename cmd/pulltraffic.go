package cmd

import (
	"context"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

const (
	pullTrafficWorkers            = 6
	pullTrafficSampleRepoLimit    = 10
	pullTrafficTopTags            = 5
	pullTrafficRepositoryPageSize = 100
	pullTrafficMaxRepoPages       = 20
	pullTrafficKindPull           = "pull"
	pullTrafficKindPush           = "push"
	pullTrafficDateLayout         = "2006-01-02"
	pullTrafficTopConsumers       = 5
	pullTrafficRetentionDays      = 30
	pullTrafficAnonymousPuller    = "anonymous"
	pullTrafficConsumerTimeFormat = "2006-01-02 15:04"
	pullTrafficHeaderRepository   = "REPOSITORY"
	pullTrafficHeaderPulls        = "PULLS"
)

var (
	pullTrafficNamespace string
	pullTrafficStartDate string
	pullTrafficEndDate   string
	pullTrafficTop       int
	pullTrafficMaxPages  int
	pullTrafficByTag     bool
	pullTrafficConsumers bool
)

// reportMode selects the breakdown a report renders and samples.
type reportMode int

const (
	modeRanking reportMode = iota
	modeTags
	modeConsumers
)

// reportModeFromFlags picks the mode and rejects the ambiguous combination.
func reportModeFromFlags(byTag, consumers bool) (reportMode, error) {
	switch {
	case byTag && consumers:
		return modeRanking, fmt.Errorf("--by-tag and --consumers are mutually exclusive")
	case consumers:
		return modeConsumers, nil
	case byTag:
		return modeTags, nil
	default:
		return modeRanking, nil
	}
}

// pullTrafficOptions configures one pull-traffic report.
type pullTrafficOptions struct {
	Namespace string
	StartDate string
	EndDate   string
	Top       int
	MaxPages  int
	Mode      reportMode
}

// pullTrafficTag is one tag and how often it was pulled in the sampled logs.
type pullTrafficTag struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// pullTrafficConsumer is one identity pulling a repository, usually a cluster
// authenticating with its own pull secret.
type pullTrafficConsumer struct {
	Puller    string `json:"puller"`
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Pulls     int    `json:"pulls"`
	LastSeen  string `json:"last_seen,omitempty"`
}

// reportRowStatus summarizes how complete one repository row is.
type reportRowStatus string

const (
	rowStatusOK        reportRowStatus = "ok"
	rowStatusPartial   reportRowStatus = "partial"
	rowStatusTruncated reportRowStatus = "truncated"
	rowStatusError     reportRowStatus = "error"
)

// pullTrafficRepo is one repository's traffic within the requested range.
type pullTrafficRepo struct {
	Repository         string                `json:"repository"`
	Pulls              int                   `json:"pulls"`
	Pushes             int                   `json:"pushes"`
	SharePercent       float64               `json:"share_percent"`
	Status             reportRowStatus       `json:"status,omitempty"`
	TopTags            []pullTrafficTag      `json:"top_tags,omitempty"`
	UntaggedPulls      int                   `json:"untagged_pulls,omitempty"`
	TagsTruncated      bool                  `json:"tags_truncated,omitempty"`
	Consumers          []pullTrafficConsumer `json:"consumers,omitempty"`
	ConsumersTruncated bool                  `json:"consumers_truncated,omitempty"`
	Error              string                `json:"error,omitempty"`

	firstEvent time.Time
	lastEvent  time.Time
}

// setStatus derives the row status. Error wins over truncation, and truncation
// wins over a complete row.
func (r *pullTrafficRepo) setStatus() {
	switch {
	case r.Error != "":
		r.Status = rowStatusError
	case r.TagsTruncated || r.ConsumersTruncated:
		r.Status = rowStatusTruncated
	default:
		r.Status = rowStatusOK
	}
}

// pullTrafficReport is the machine-readable form of `list pull-traffic`.
type pullTrafficReport struct {
	Namespace    string            `json:"namespace"`
	StartDate    string            `json:"start_date"`
	EndDate      string            `json:"end_date"`
	CoveredStart string            `json:"covered_start,omitempty"`
	CoveredEnd   string            `json:"covered_end,omitempty"`
	TotalPulls   int               `json:"total_pulls"`
	TotalPushes  int               `json:"total_pushes"`
	Repositories []pullTrafficRepo `json:"repositories"`
}

var pullTrafficCmd = &cobra.Command{
	Use:   "pull-traffic",
	Short: "Rank repositories by image pulls",
	Long: `Rank every repository in a namespace by the number of image pulls recorded
in a date range.

Quay keeps aggregated usage logs for roughly 30 days and rejects ranges longer
than 30 days. Longer ranges are split into consecutive 30-day windows and
summed, so windows outside the retention period simply contribute nothing.
Without --startdate/--enddate the last 30 days are used.

Dates accept MM/DD/YYYY or YYYY-MM-DD. Output is a table unless -O json or
-O yaml is given.`,
	Example: `  go-quay list pull-traffic -n myorg --token "$QUAY_TOKEN"
  go-quay list pull-traffic -n myorg --top 5 --by-tag
  go-quay list pull-traffic -n myorg --top 5 --consumers`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		// A ranking is meant to be read; JSON stays one flag away.
		if !cmd.Flags().Changed(outputFlagName) {
			outputFormat = outputTable
		}

		now := time.Now()
		start, err := resolveReportDate(pullTrafficStartDate, now)
		if err != nil {
			return err
		}
		end, err := resolveReportDate(pullTrafficEndDate, now)
		if err != nil {
			return err
		}
		if start == "" && end == "" {
			start, end = defaultReportWindow(now)
		}

		mode, err := reportModeFromFlags(pullTrafficByTag, pullTrafficConsumers)
		if err != nil {
			return err
		}

		opts := pullTrafficOptions{
			Namespace: pullTrafficNamespace,
			StartDate: start,
			EndDate:   end,
			Top:       pullTrafficTop,
			MaxPages:  pullTrafficMaxPages,
			Mode:      mode,
		}

		report, err := collectPullTraffic(cmd.Context(), client, opts)
		if err != nil {
			return err
		}
		return printPullTraffic(report, opts)
	},
}

func init() {
	listCmd.AddCommand(pullTrafficCmd)
	pullTrafficCmd.Flags().StringVarP(&pullTrafficNamespace, "namespace", "n", appCfg.Namespace, "Repository namespace (default: config file)")
	pullTrafficCmd.Flags().StringVarP(&pullTrafficStartDate, "startdate", "s", "", "Start date (MM/DD/YYYY or YYYY-MM-DD)")
	pullTrafficCmd.Flags().StringVarP(&pullTrafficEndDate, "enddate", "e", "", "End date (MM/DD/YYYY or YYYY-MM-DD)")
	pullTrafficCmd.Flags().IntVar(&pullTrafficTop, "top", 0, "Only show the top N repositories (0 shows all)")
	pullTrafficCmd.Flags().BoolVar(&pullTrafficByTag, "by-tag", false, "Also rank tags, sampling individual logs")
	pullTrafficCmd.Flags().BoolVar(&pullTrafficConsumers, "consumers", false, "Also show who is pulling, sampling individual logs")
	pullTrafficCmd.Flags().IntVar(&pullTrafficMaxPages, "max-pages", 25, "Maximum log pages per repository when using --by-tag")
}

// collectPullTraffic builds the report. Per-repository failures are recorded on
// the row instead of aborting the run; the command only fails when no
// repository could be read at all.
func collectPullTraffic(ctx context.Context, client *lib.Client, opts pullTrafficOptions) (*pullTrafficReport, error) {
	if opts.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if _, _, err := lib.ValidateLogDateRange(opts.StartDate, opts.EndDate); err != nil {
		return nil, err
	}

	repositories, err := listPullTrafficRepositories(ctx, client, opts.Namespace)
	if err != nil {
		return nil, err
	}

	report := &pullTrafficReport{
		Namespace:    opts.Namespace,
		StartDate:    opts.StartDate,
		EndDate:      opts.EndDate,
		Repositories: pullTrafficRows(ctx, client, opts, repositories),
	}
	if start, end := coveredWindow(report.Repositories); !start.IsZero() {
		report.CoveredStart = start.Format(pullTrafficDateLayout)
		report.CoveredEnd = end.Format(pullTrafficDateLayout)
	}
	if err := summarizePullTraffic(report, opts); err != nil {
		return nil, err
	}
	if opts.Mode != modeRanking {
		collectPullTrafficDetails(ctx, client, opts, report)
	}
	return report, nil
}

// pullTrafficRows fetches every repository's aggregate logs concurrently.
func pullTrafficRows(ctx context.Context, client *lib.Client, opts pullTrafficOptions, repositories []string) []pullTrafficRepo {
	rows := make([]pullTrafficRepo, len(repositories))
	sem := make(chan struct{}, pullTrafficWorkers)
	var wg sync.WaitGroup
	for i, name := range repositories {
		wg.Add(1)
		go func(i int, name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			rows[i] = pullTrafficForRepo(ctx, client, opts, name)
		}(i, name)
	}
	wg.Wait()
	return rows
}

// summarizePullTraffic fills in totals and shares, orders the rows by pulls, and
// applies the --top limit.
func summarizePullTraffic(report *pullTrafficReport, opts pullTrafficOptions) error {
	failed := 0
	for _, row := range report.Repositories {
		report.TotalPulls += row.Pulls
		report.TotalPushes += row.Pushes
		if row.Error != "" {
			failed++
		}
	}
	if len(report.Repositories) > 0 && failed == len(report.Repositories) {
		return fmt.Errorf("no repository logs could be read: %s", report.Repositories[0].Error)
	}
	for i := range report.Repositories {
		if report.TotalPulls > 0 {
			share := float64(report.Repositories[i].Pulls) / float64(report.TotalPulls) * 100
			report.Repositories[i].SharePercent = math.Round(share*100) / 100
		}
	}

	sort.SliceStable(report.Repositories, func(i, j int) bool {
		if report.Repositories[i].Pulls != report.Repositories[j].Pulls {
			return report.Repositories[i].Pulls > report.Repositories[j].Pulls
		}
		return report.Repositories[i].Repository < report.Repositories[j].Repository
	})
	if opts.Top > 0 && opts.Top < len(report.Repositories) {
		report.Repositories = report.Repositories[:opts.Top]
	}
	return nil
}

func pullTrafficForRepo(ctx context.Context, client *lib.Client, opts pullTrafficOptions, name string) pullTrafficRepo {
	row := pullTrafficRepo{Repository: name}
	logs, err := client.GetAggregatedLogsRange(ctx, opts.Namespace, name, opts.StartDate, opts.EndDate)
	if err != nil {
		row.Error = err.Error()
		row.setStatus()
		return row
	}
	for _, entry := range logs.Aggregated {
		if when, ok := parseLogDatetime(entry.Datetime); ok {
			if row.firstEvent.IsZero() || when.Before(row.firstEvent) {
				row.firstEvent = when
			}
			if row.lastEvent.IsZero() || when.After(row.lastEvent) {
				row.lastEvent = when
			}
		}
		switch {
		case strings.HasPrefix(entry.Kind, pullTrafficKindPull):
			row.Pulls += entry.Count
		case strings.HasPrefix(entry.Kind, pullTrafficKindPush):
			row.Pushes += entry.Count
		}
	}
	row.setStatus()
	return row
}

// parseLogDatetime parses the RFC1123Z timestamps Quay uses for log entries.
func parseLogDatetime(value string) (time.Time, bool) {
	when, err := time.Parse(time.RFC1123Z, value)
	if err != nil {
		return time.Time{}, false
	}
	return when, true
}

// coveredWindow returns the span Quay actually had data for. Usage logs are
// retained for roughly 30 days, so a longer request can succeed while covering
// far less time than asked.
func coveredWindow(rows []pullTrafficRepo) (time.Time, time.Time) {
	var start, end time.Time
	for _, row := range rows {
		if row.firstEvent.IsZero() {
			continue
		}
		if start.IsZero() || row.firstEvent.Before(start) {
			start = row.firstEvent
		}
		if end.IsZero() || row.lastEvent.After(end) {
			end = row.lastEvent
		}
	}
	return start, end
}

// collectPullTrafficDetails samples individual log entries for the repositories
// being displayed and records tag and consumer breakdowns. Individual logs are
// far more expensive than aggregated ones, so the sample is capped by
// --max-pages. One sample feeds every enabled breakdown.
func collectPullTrafficDetails(ctx context.Context, client *lib.Client, opts pullTrafficOptions, report *pullTrafficReport) {
	limit := len(report.Repositories)
	if opts.Top <= 0 && limit > pullTrafficSampleRepoLimit {
		limit = pullTrafficSampleRepoLimit
	}

	// Rows past the sample cap still have aggregate data but no breakdown.
	for i := limit; i < len(report.Repositories); i++ {
		if report.Repositories[i].Pulls > 0 {
			report.Repositories[i].Status = rowStatusPartial
		}
	}

	for i := 0; i < limit; i++ {
		row := &report.Repositories[i]
		if row.Pulls == 0 {
			continue
		}
		entries, truncated, err := client.GetAllLogs(ctx, opts.Namespace, row.Repository, opts.StartDate, opts.EndDate, opts.MaxPages)
		if err != nil {
			row.Error = err.Error()
			row.setStatus()
			continue
		}
		if opts.Mode == modeTags {
			row.TopTags, row.UntaggedPulls = rankPullTrafficTags(entries)
			row.TagsTruncated = truncated
		}
		if opts.Mode == modeConsumers {
			row.Consumers = rankPullTrafficConsumers(entries)
			row.ConsumersTruncated = truncated
		}
		row.setStatus()
	}
}

// rankPullTrafficConsumers counts pulls per puller identity and IP, which is
// what separates one cluster pulling repeatedly from many clusters doing it.
func rankPullTrafficConsumers(entries []lib.LogEntry) []pullTrafficConsumer {
	type consumerKey struct{ puller, ip string }

	aggregated := map[consumerKey]*pullTrafficConsumer{}
	lastSeen := map[consumerKey]time.Time{}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Kind, pullTrafficKindPull) {
			continue
		}
		puller := entry.Performer.Name
		if puller == "" {
			puller = pullTrafficAnonymousPuller
		}
		key := consumerKey{puller: puller, ip: entry.IP}
		consumer, ok := aggregated[key]
		if !ok {
			consumer = &pullTrafficConsumer{
				Puller:    puller,
				IP:        entry.IP,
				UserAgent: entry.Metadata.UserAgent,
			}
			aggregated[key] = consumer
		}
		consumer.Pulls++
		if when, ok := parseLogDatetime(entry.Datetime); ok && when.After(lastSeen[key]) {
			lastSeen[key] = when
		}
	}

	consumers := make([]pullTrafficConsumer, 0, len(aggregated))
	for key, consumer := range aggregated {
		if !lastSeen[key].IsZero() {
			consumer.LastSeen = lastSeen[key].Format(pullTrafficConsumerTimeFormat)
		}
		consumers = append(consumers, *consumer)
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
	if len(consumers) > pullTrafficTopConsumers {
		consumers = consumers[:pullTrafficTopConsumers]
	}
	return consumers
}

// rankPullTrafficTags counts pulls per tag and returns the busiest tags.
func rankPullTrafficTags(entries []lib.LogEntry) ([]pullTrafficTag, int) {
	counts := map[string]int{}
	untagged := 0
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Kind, pullTrafficKindPull) {
			continue
		}
		if entry.Metadata.Tag == "" {
			untagged++
			continue
		}
		counts[entry.Metadata.Tag]++
	}

	tags := make([]pullTrafficTag, 0, len(counts))
	for tag, count := range counts {
		tags = append(tags, pullTrafficTag{Tag: tag, Count: count})
	}
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].Count != tags[j].Count {
			return tags[i].Count > tags[j].Count
		}
		return tags[i].Tag < tags[j].Tag
	})
	if len(tags) > pullTrafficTopTags {
		tags = tags[:pullTrafficTopTags]
	}
	return tags, untagged
}

// listPullTrafficRepositories returns every repository name in the namespace.
func listPullTrafficRepositories(ctx context.Context, client *lib.Client, namespace string) ([]string, error) {
	names := []string{}
	_, err := lib.Paginate(pullTrafficMaxRepoPages, func(page int, _ string) (lib.PageResult, error) {
		list, err := client.ListRepositories(ctx, namespace, false, false, false, page, pullTrafficRepositoryPageSize)
		if err != nil {
			return lib.PageResult{}, fmt.Errorf("listing repositories: %w", err)
		}

		for _, repository := range list.Repositories {
			names = append(names, repository.Name)
		}
		return lib.PageResult{HasMore: list.HasAdditional}, nil
	})
	if err != nil {
		return nil, err
	}
	if len(names) > 0 {
		return names, nil
	}

	// Quay answers an unauthenticated namespace listing with 200 and an empty
	// result, which looks exactly like an empty namespace. Public repositories
	// are visible without a token, so if any exist the caller's token was not
	// accepted and reporting zero pull traffic would be wrong.
	public, err := client.ListRepositories(ctx, namespace, true, false, false, 1, pullTrafficRepositoryPageSize)
	if err != nil {
		return nil, fmt.Errorf("listing public repositories: %w", err)
	}
	if len(public.Repositories) > 0 {
		return nil, fmt.Errorf(
			"no repositories visible in namespace %q, but %d public repositories exist: the token was not accepted (check --token or $QUAY_TOKEN)",
			namespace, len(public.Repositories))
	}
	return names, nil
}

func printPullTraffic(report *pullTrafficReport, opts pullTrafficOptions) error {
	if outputFormat != outputTable && outputFormat != outputCSV {
		return printJSON(report)
	}
	headers, rows := pullTrafficTable(report, opts)
	if outputFormat == outputTable {
		printPullTrafficSummary(report)
	}
	return printRows(headers, rows)
}

// pullTrafficTable selects the column layout that matches the requested
// breakdown. Table and CSV output share it.
func pullTrafficTable(report *pullTrafficReport, opts pullTrafficOptions) ([]string, [][]string) {
	switch opts.Mode {
	case modeConsumers:
		return pullTrafficConsumerTable(report)
	case modeTags:
		return pullTrafficTagTable(report)
	default:
		return pullTrafficRankingTable(report)
	}
}

func printPullTrafficSummary(report *pullTrafficReport) {
	fmt.Fprintf(os.Stderr, "Pull traffic for '%s' (%s): %d pulls, %d pushes\n",
		report.Namespace, reportWindowLabel(report.StartDate, report.EndDate, report.CoveredStart, report.CoveredEnd),
		report.TotalPulls, report.TotalPushes)
	if counts := report.statusCounts(); counts[rowStatusPartial]+counts[rowStatusTruncated]+counts[rowStatusError] > 0 {
		fmt.Fprintf(os.Stderr, "Rows: %d partial, %d truncated, %d errors\n",
			counts[rowStatusPartial], counts[rowStatusTruncated], counts[rowStatusError])
	}
}

// statusCounts tallies rows by their derived status.
func (r *pullTrafficReport) statusCounts() map[reportRowStatus]int {
	counts := map[reportRowStatus]int{}
	for _, row := range r.Repositories {
		counts[row.Status]++
	}
	return counts
}

func pullTrafficRankingTable(report *pullTrafficReport) ([]string, [][]string) {
	rows := make([][]string, 0, len(report.Repositories))
	for _, repository := range report.Repositories {
		rows = append(rows, []string{
			repository.Repository,
			strconv.Itoa(repository.Pulls),
			strconv.Itoa(repository.Pushes),
			fmt.Sprintf("%.1f%%", repository.SharePercent),
		})
	}
	return []string{pullTrafficHeaderRepository, pullTrafficHeaderPulls, "PUSHES", "SHARE"}, rows
}

func pullTrafficTagTable(report *pullTrafficReport) ([]string, [][]string) {
	rows := make([][]string, 0, len(report.Repositories))
	for _, repository := range report.Repositories {
		tags := make([]string, 0, len(repository.TopTags))
		for _, tag := range repository.TopTags {
			tags = append(tags, fmt.Sprintf("%s (%d)", tag.Tag, tag.Count))
		}
		rows = append(rows, []string{
			repository.Repository,
			strconv.Itoa(repository.Pulls),
			strings.Join(tags, ", "),
		})
	}
	return []string{pullTrafficHeaderRepository, pullTrafficHeaderPulls, "TAGS"}, rows
}

func pullTrafficConsumerTable(report *pullTrafficReport) ([]string, [][]string) {
	rows := make([][]string, 0, len(report.Repositories))
	for _, repository := range report.Repositories {
		for _, consumer := range repository.Consumers {
			rows = append(rows, []string{
				repository.Repository,
				consumer.Puller,
				consumer.IP,
				strconv.Itoa(consumer.Pulls),
				consumer.LastSeen,
			})
		}
	}
	return []string{pullTrafficHeaderRepository, "PULLER", "IP", pullTrafficHeaderPulls, "LAST SEEN"}, rows
}
