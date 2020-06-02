// +build integration

package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"github.com/taxibeat/bollobas/internal/storagetest"
)

func TestGetAllOperatorStats(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewOperatorStatsRepository(st)
	err = populateOperatorStatsTable(at)
	assert.Nil(t, err)

	rr, err := at.GetAll(context.Background(), internal.DateFilter{}, internal.Pagination{})
	assert.Nil(t, err)

	assert.Len(t, rr, 2)
	assert.Equal(t, int64(2), rr[0].ID)
	assert.Equal(t, "asd", rr[0].OperatorID)
	assert.Equal(t, "AgeRange2", rr[0].AgeRange)
	assert.Equal(t, 1, rr[0].Gender)
	assert.Equal(t, int64(1), rr[1].ID)
	assert.Equal(t, "ass", rr[1].OperatorID)
	assert.Equal(t, "AgeRange1", rr[1].AgeRange)
	assert.Equal(t, 3, rr[1].Gender)
}

func TestFilteredStatsQuery(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewOperatorStatsRepository(st)
	err = populateOperatorStatsTable(at)
	assert.Nil(t, err)

	f := func(ctx context.Context, filter internal.DateFilter) (interface{}, error) {

		return at.GetAll(ctx, filter, internal.Pagination{})
	}

	storagetest.TestFilteredQuery(t, f)
}

func populateOperatorStatsTable(r *sql.OperatorStatsRepo) error {
	ctx := context.Background()
	r.DB().Exec(ctx, "TRUNCATE operator_stats")
	id1 := "asd"
	id2 := "qwe"
	age1 := "AgeRange1"
	age2 := "AgeRange2"
	gend1 := 3
	gend2 := 1

	a := &internal.OperatorStats{
		Date:       time.Now().AddDate(0, -1, 0),
		OperatorID: &id1,
		Gender:     &gend1,
		AgeRange:   &age1,
	}

	if err := r.Add(ctx, a); err != nil {
		return err
	}
	a = &internal.OperatorStats{
		Date:       time.Now().AddDate(0, -1, 0).Add(time.Hour),
		OperatorID: &id2,
		Gender:     &gend2,
		AgeRange:   &age2,
	}
	return r.Add(context.Background(), a)
}
