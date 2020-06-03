package internal

import "time"

// DateFilter includes teh details for filtering by dates
type DateFilter struct {
	From *time.Time
	To   *time.Time
}

// Pagination includes the pagination parameters
type Pagination struct {
	Offset int
	Limit  int
}
