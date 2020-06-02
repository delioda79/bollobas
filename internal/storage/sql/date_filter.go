package sql

import (
	"fmt"
	"github.com/taxibeat/bollobas/internal"
	"time"
)

// DateFilter filters a sql query by dates
type DateFilter struct {
	*internal.DateFilter
}

// Filter performs filtering
func (df DateFilter) Filter(query string) (string, []interface{}) {
	var params []interface{}
	f := ""

	fq, ft := df.from()
	if fq != "" {
		f += fq
		params = append(params, *ft)
	}

	tq, tt := df.to()
	if tq != "" {
		f += tq
		params = append(params, *tt)
	}
	return fmt.Sprintf(query, f), params
}

func (df DateFilter) from() (string, *time.Time) {
	if df.DateFilter.From == nil {
		return "", nil
	}
	return fmt.Sprintf(" AND date >= ?"), df.DateFilter.From
}

func (df DateFilter) to() (string, *time.Time) {
	if df.DateFilter.To == nil {
		return "", nil
	}
	return fmt.Sprintf(" AND date <= ?"), df.DateFilter.To
}
