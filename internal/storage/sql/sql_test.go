// +build integration

package sql_test

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal/storagetest"
	"log"
	"testing"
)

// TestMain will exec each test, one by one
func TestMain(m *testing.M) {
	err := godotenv.Load("../../../config/.env.test")
	if err != nil {
		log.Printf("no .env.test file exists: %v\n", err)
	}

	// exec test and this returns an exit code to pass to os
	m.Run()
}

// TestNew tests sql storage
func TestNew(t *testing.T) {
	store, err := storagetest.SetConfig()
	assert.Nil(t, err)
	assert.NotNil(t, store)
	_, e := store.DB().Query(context.TODO(), "SELECT * FROM aggregated_trips")
	if e != nil {
		t.Error(e)
	}
}
