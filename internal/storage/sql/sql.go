package sql

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/beatlabs/patron/client/sql"
	"github.com/beatlabs/patron/log"
	"github.com/taxibeat/bollobas/internal/config"

	// this is the mysqld driver, this comment is needed by the smart linter...
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

// RemoveDataInTable deletes everything in the table
func (s *Store) RemoveDataInTable(ctx context.Context, tables ...string) error {
	q := "TRUNCATE TABLE %s"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.New("cannot begin transaction")
	}
	for _, t := range tables {
		if _, err := tx.Exec(ctx, fmt.Sprintf(q, t)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("unable to clear table %s: %v", t, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable to commit transaction: %v", err)
	}

	return nil
}

// DB returns the db connection pool
func (s *Store) DB() *sql.DB {
	return s.db
}
