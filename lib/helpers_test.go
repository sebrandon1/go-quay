package lib

import (
	"testing"
	"time"
)

func TestLatestTagTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		tags     []Tag
		wantZero bool
		wantUnix int64
	}{
		{
			name:     "no tags",
			tags:     nil,
			wantZero: true,
		},
		{
			name: "single tag",
			tags: []Tag{
				{Name: testTagNameLatest, StartTs: 1735689600},
			},
			wantUnix: 1735689600,
		},
		{
			name: "multiple tags picks latest",
			tags: []Tag{
				{Name: "v1", StartTs: 1704067200},
				{Name: "v2", StartTs: 1752566400},
				{Name: "v3", StartTs: 1748736000},
			},
			wantUnix: 1752566400,
		},
		{
			name: "zero StartTs ignored",
			tags: []Tag{
				{Name: "zero", StartTs: 0},
				{Name: "good", StartTs: 1735689600},
			},
			wantUnix: 1735689600,
		},
		{
			name:     "all zero StartTs",
			tags:     []Tag{{Name: "a", StartTs: 0}},
			wantZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RepositoryWithTags{
				Tags: RepositoryTags{Tags: tt.tags},
			}
			got := r.LatestTagTimestamp()
			if tt.wantZero {
				if !got.IsZero() {
					t.Errorf("expected zero time, got %v", got)
				}
				return
			}
			want := time.Unix(tt.wantUnix, 0)
			if !got.Equal(want) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	}
}

func TestParseQuayTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantZero bool
	}{
		{"RFC1123Z", "Mon, 02 Jan 2006 15:04:05 -0700", false},
		{"RFC1123", "Mon, 02 Jan 2006 15:04:05 MST", false},
		{"empty", "", true},
		{"garbage", "not-a-timestamp", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseQuayTimestamp(tt.input)
			if tt.wantZero && !got.IsZero() {
				t.Errorf("expected zero time, got %v", got)
			}
			if !tt.wantZero && got.IsZero() {
				t.Error("expected non-zero time, got zero")
			}
		})
	}
}

func TestCalculateDaysSince(t *testing.T) {
	threeDaysAgo := time.Now().Add(-3 * 24 * time.Hour)
	days := CalculateDaysSince(threeDaysAgo)
	if days != 3 {
		t.Errorf("expected 3 days, got %d", days)
	}
}

func TestCalculateDaysSince_ZeroTime(t *testing.T) {
	days := CalculateDaysSince(time.Time{})
	if days <= 0 {
		t.Errorf("expected positive days for zero time, got %d", days)
	}
}
