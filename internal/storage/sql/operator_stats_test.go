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

	rr, err := at.GetAll(context.Background())
	assert.Nil(t, err)

	assert.Len(t, rr, 2)
	assert.Equal(t, 2, rr[0].ID)
	assert.Equal(t, 2, rr[0].OperatorID)
	assert.Equal(t, "AgeRange2", rr[0].AgeRange)
	assert.Equal(t, 1, rr[0].Gender)
	assert.Equal(t, 1, rr[1].ID)
	assert.Equal(t, 1, rr[0].OperatorID)
	assert.Equal(t, "AgeRange1", rr[0].AgeRange)
	assert.Equal(t, 3, rr[0].Gender)
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
		Date:       time.Now(),
		OperatorID: "asd",
		Gender:     1,
		AgeRange:   "AgeRange2",
	}
	return r.Add(context.Background(), a)
}
