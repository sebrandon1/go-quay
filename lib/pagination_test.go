package lib

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestPaginateStopsWhenNoMorePages(t *testing.T) {
	pages := []PageResult{{HasMore: true}, {HasMore: true}, {HasMore: false}}
	calls := 0

	truncated, err := Paginate(0, func(page int, _ string) (PageResult, error) {
		if page != calls+1 {
			t.Fatalf("page = %d, want %d", page, calls+1)
		}
		calls++
		return pages[page-1], nil
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
	if truncated {
		t.Fatal("truncated = true, want false")
	}
}

func TestPaginateHonorsMaxPages(t *testing.T) {
	calls := 0

	truncated, err := Paginate(2, func(int, string) (PageResult, error) {
		calls++
		return PageResult{HasMore: true}, nil
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
	if !truncated {
		t.Fatal("truncated = false, want true when more pages remain")
	}
}

func TestPaginatePassesPageAndToken(t *testing.T) {
	tokens := []string{"", "cursor-1", "cursor-2"}
	var got []struct {
		page  int
		token string
	}

	_, err := Paginate(0, func(page int, token string) (PageResult, error) {
		got = append(got, struct {
			page  int
			token string
		}{page, token})
		if page < len(tokens) {
			return PageResult{HasMore: true, Next: tokens[page]}, nil
		}
		return PageResult{HasMore: false}, nil
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}

	want := []struct {
		page  int
		token string
	}{{1, ""}, {2, "cursor-1"}, {3, "cursor-2"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("calls = %#v, want %#v", got, want)
	}
}

func TestPaginatePropagatesError(t *testing.T) {
	wantErr := errors.New("page fetch failed")
	calls := 0

	_, err := Paginate(0, func(page int, _ string) (PageResult, error) {
		calls++
		if page == 2 {
			return PageResult{}, wantErr
		}
		return PageResult{HasMore: true}, nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func TestPaginateSinglePageNeedsNoCap(t *testing.T) {
	calls := 0

	truncated, err := Paginate(5, func(int, string) (PageResult, error) {
		calls++
		return PageResult{}, nil
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if truncated {
		t.Fatal("truncated = true, want false")
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestPaginateCapsUnboundedWalk(t *testing.T) {
	old := maxPaginationPages
	maxPaginationPages = 3
	t.Cleanup(func() { maxPaginationPages = old })

	calls := 0
	_, err := Paginate(0, func(int, string) (PageResult, error) {
		calls++
		return PageResult{HasMore: true}, nil
	})
	if err == nil || !strings.Contains(err.Error(), "exceeded 3 pages") {
		t.Fatalf("error = %v, want page-ceiling error", err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
}

func TestPaginateRejectsRepeatedCursor(t *testing.T) {
	calls := 0
	_, err := Paginate(0, func(int, string) (PageResult, error) {
		calls++
		return PageResult{HasMore: true, Next: "same-cursor"}, nil
	})
	if err == nil || !strings.Contains(err.Error(), "did not advance") {
		t.Fatalf("error = %v, want non-advancing-cursor error", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func TestPaginateAllowsDistinctCursors(t *testing.T) {
	next := []string{"page-two", "page-three"}
	calls := 0

	truncated, err := Paginate(0, func(_ int, token string) (PageResult, error) {
		if calls > 0 && token != next[calls-1] {
			t.Fatalf("token = %q, want %q", token, next[calls-1])
		}
		calls++
		if calls > len(next) {
			return PageResult{}, nil
		}
		return PageResult{HasMore: true, Next: next[calls-1]}, nil
	})
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if truncated {
		t.Fatal("truncated = true, want false")
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
}
