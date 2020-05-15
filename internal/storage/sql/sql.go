package sql

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/beatlabs/patron/client/sql"
	"github.com/beatlabs/patron/log"
	"github.com/taxibeat/bollobas/internal/config"

	// this is the mysqld driver, thsi comment is needed by the smart linter...
	_ "github.com/go-sql-driver/mysql"
)

var (
	// ErrNoRecordUpdated an error when no record is updated
	ErrNoRecordUpdated = errors.New("no record updated")
)

// Store stores any sql related info and functionality
type Store struct {
	db *sql.DB
}

// New initializes an sql configuration
func New(cfg *config.Configuration) (*Store, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUsername.Get(),
		cfg.DBPassword.Get(),
		cfg.DBWriteHost.Get(),
		cfg.DBPort.Get(),
		cfg.DBName.Get(),
	)

	conn, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	return &Store{db: conn}, nil
}

// Close the connection pool
func (s *Store) Close() {
	if err := s.db.Close(context.Background()); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

// DB returns the db connection pool
func (s *Store) DB() *sql.DB {
	return s.db
}
