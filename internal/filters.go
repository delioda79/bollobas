package internal

import (
	"fmt"
	"math"
	"time"
)

const minTimeRangeAllowed = 60 //minutes

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

// Validate ensures that the date range used is within the margins of
// the previous month and that the min period retrieved is 1 hour.
func (df *DateFilter) Validate() error {
	if df.From == nil || df.To == nil {
		return nil
	}

	today := time.Now()
	lastMonth := today.AddDate(0, -1, 0).Month()

	if df.From.Month() != lastMonth || df.To.Month() != lastMonth ||
		df.From.Year() != today.Year() || df.To.Year() != today.Year() {
		return fmt.Errorf("invalid date range requested")
	}

	if df.To.Sub(*df.From) < (time.Duration(minTimeRangeAllowed) * time.Minute) {
		return fmt.Errorf("invalid time range requested")
	}

	return nil
}

// GetOffset returns Offset
func (pg *Pagination) GetOffset() int {
	return pg.Offset
}

// GetLimit returns Limit
func (pg *Pagination) GetLimit() int {
	if pg.Limit == 0 {
		return math.MaxInt32
	}

	return pg.Limit
}
