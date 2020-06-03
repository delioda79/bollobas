package sql

import (
	"fmt"
	"math"
	"time"

	"github.com/taxibeat/bollobas/internal"
)

// AllFilter filters a sql Query by dates and pagination params
type AllFilter struct {
	internal.DateFilter
	internal.Pagination
}

// FilterDate performs filtering
func (af *AllFilter) FilterDate(query string) (string, []interface{}) {
	var params []interface{}
	f := ""

	fq, ft := af.from()
	if fq != "" {
		f += fq
		params = append(params, *ft)
	}

	tq, tt := af.to()
	if tq != "" {
		f += tq
		params = append(params, *tt)
	}
	return fmt.Sprintf(query, f), params
}

func (af *AllFilter) from() (string, *time.Time) {
	if af.DateFilter.From == nil {
		return "", nil
	}
	return fmt.Sprintf(" AND date >= ?"), af.DateFilter.From
}

func (af *AllFilter) to() (string, *time.Time) {
	if af.DateFilter.To == nil {
		return "", nil
	}
	return fmt.Sprintf(" AND date <= ?"), af.DateFilter.To
}

// Paginate performs pagination
func (af *AllFilter) Paginate() []interface{} {
	var params []interface{}

	ft := af.offset()
	params = append(params, ft)

	tt := af.limit()
	params = append(params, tt)

	return params
}

func (af *AllFilter) offset() int {
	return af.Pagination.Offset
}

func (af *AllFilter) limit() int {
	if af.Pagination.Limit == 0 {
		return math.MaxInt32
	}

	return af.Pagination.Limit
}
