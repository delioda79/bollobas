// +build integration

package sql_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"github.com/taxibeat/bollobas/internal/storagetest"
	"testing"
	"time"
)

func TestGetAllOperatorStats(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at  := sql.NewOperatorStatsRepository(context.Background(), st)
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
		Date: time.Now(),
		OperatorID: 1,
		Gender:3,
		AgeRange:"AgeRange1",
	}

	if err := r.Add(ctx, a); err != nil {
		return err
	}
	a = &internal.OperatorStats{
		Date: time.Now(),
		OperatorID: 2,
		Gender:1,
		AgeRange:"AgeRange2",
	}
	return r.Add(context.Background(), a)
}
