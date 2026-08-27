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
