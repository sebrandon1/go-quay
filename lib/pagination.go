package lib

// PageResult describes one page returned to Paginate.
type PageResult struct {
	// HasMore reports whether the API has another page.
	HasMore bool
	// Next is the opaque cursor for token-based listings. Page-based listings
	// leave it empty and rely on the page number.
	Next string
}

// Paginate walks a paginated API until it reports no more pages or maxPages
// pages have been read. maxPages <= 0 means no limit. The callback receives the
// 1-based page number and the previous page's token; token-based APIs ignore the
// page number.
//
// It returns whether maxPages stopped a walk that still had more pages.
func Paginate(maxPages int, fetch func(page int, token string) (PageResult, error)) (bool, error) {
	page, token, hasMore := 1, "", false
	for fetched := 0; ; fetched++ {
		if maxPages > 0 && fetched >= maxPages {
			return hasMore, nil
		}

		result, err := fetch(page, token)
		if err != nil {
			return false, err
		}

		hasMore = result.HasMore
		if !hasMore {
			return false, nil
		}
		page++
		token = result.Next
	}
}
