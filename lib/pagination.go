package lib

import "fmt"

// maxPaginationPages caps a walk that was not given a limit of its own, so a
// server that never reports the last page cannot loop forever. It is a
// variable so tests can lower it.
var maxPaginationPages = 100000

// PageResult describes one page returned to Paginate.
type PageResult struct {
	// HasMore reports whether the API has another page.
	HasMore bool
	// Next is the opaque cursor for token-based listings. Page-based listings
	// leave it empty and rely on the page number.
	Next string
}

// Paginate walks a paginated API until it reports no more pages or maxPages
// pages have been read. maxPages <= 0 means no limit, apart from the safety
// ceiling in maxPaginationPages. The callback receives the 1-based page number
// and the previous page's token; token-based APIs ignore the page number.
//
// It returns whether maxPages stopped a walk that still had more pages.
func Paginate(maxPages int, fetch func(page int, token string) (PageResult, error)) (bool, error) {
	page, token, hasMore := 1, "", false
	for fetched := 0; ; fetched++ {
		if maxPages > 0 && fetched >= maxPages {
			return hasMore, nil
		}
		if maxPages <= 0 && page > maxPaginationPages {
			return false, fmt.Errorf("pagination exceeded %d pages; aborting", maxPaginationPages)
		}

		result, err := fetch(page, token)
		if err != nil {
			return false, err
		}

		hasMore = result.HasMore
		if !hasMore {
			return false, nil
		}
		// A walk with no caller-supplied limit must still terminate: a cursor
		// that comes back unchanged would repeat the same page forever. A
		// capped walk is already bounded, so it keeps its existing behavior.
		if maxPages <= 0 && result.Next != "" && result.Next == token {
			return false, fmt.Errorf("pagination cursor %q did not advance; aborting to avoid an infinite loop", result.Next)
		}
		page++
		token = result.Next
	}
}
