// +build integration

package sql_test

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"testing"
	"time"

	"github.com/taxibeat/bollobas/internal/storagetest"
)

func TestGetAllOperatorStats(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewOperatorStatsRepository(st)
	err = populateOperatorStatsTable(at)
	assert.Nil(t, err)

	rr, err := at.GetAll(context.Background(), internal.DateFilter{})
	assert.Nil(t, err)

	assert.Len(t, rr, 2)
	assert.Equal(t, int64(2), rr[0].ID)
	assert.Equal(t, 2, rr[0].OperatorID)
	assert.Equal(t, "AgeRange2", rr[0].AgeRange)
	assert.Equal(t, 1, rr[0].Gender)
	assert.Equal(t, int64(1), rr[1].ID)
	assert.Equal(t, 1, rr[1].OperatorID)
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

		return at.GetAll(ctx, filter)
	}

	storagetest.TestFilteredQuery(t, f)
}

func populateOperatorStatsTable(r *sql.OperatorStatsRepo) error {
	ctx := context.Background()
	r.DB().Exec(ctx, "TRUNCATE operator_stats")
	a := &internal.OperatorStats{
		Date:       time.Now(),
		OperatorID: "ass",
		Gender:     3,
		AgeRange:   "AgeRange1",
	}

	if err := r.Add(ctx, a); err != nil {
		return err
	}
	a = &internal.OperatorStats{
		Date:       time.Now().Add(time.Hour),
		OperatorID: "asd",
		Gender:     1,
		AgeRange:   "AgeRange2",
	}
	return r.Add(context.Background(), a)
}
