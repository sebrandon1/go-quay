package lib

import "time"

// LatestTagTimestamp returns the most recent timestamp across all tags,
// using the StartTs unix epoch field (more reliable than parsing LastModified).
func (r *RepositoryWithTags) LatestTagTimestamp() time.Time {
	var latest int64
	for _, tag := range r.Tags.Tags {
		if tag.StartTs > latest {
			latest = tag.StartTs
		}
	}
	if latest == 0 {
		return time.Time{}
	}
	return time.Unix(latest, 0)
}

// TagsTruncated reports whether Tags holds only the first page of results.
//
// GetRepository fetches a single page of tags, so any count derived from Tags is
// a lower bound once this returns true. Quay returns tags newest-first, so
// LatestTagTimestamp is still the newest tag even for a truncated list; call
// ListAllTags if an exact count is required.
func (r *RepositoryWithTags) TagsTruncated() bool {
	return r.Tags.HasAdditional
}

// ParseQuayTimestamp parses a Quay.io timestamp string.
// Quay uses RFC 2822 / RFC 1123Z format with a fallback to RFC 1123.
func ParseQuayTimestamp(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC1123Z, s)
	if err != nil {
		t, err = time.Parse(time.RFC1123, s)
		if err != nil {
			return time.Time{}
		}
	}
	return t
}

func CalculateDaysSince(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}
