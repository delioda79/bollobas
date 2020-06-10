package main

import (
	"context"
	"fmt"
	"github.com/beatlabs/patron/client/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	softDelete = "soft"
	hardDelete = "hard"
)

var tables = []string{
	"aggregated_trips",
	"operator_stats",
	"traffic_incidents",
}

type config struct {
	DBUsername         string
	DBPassword         string
	DBName             string
	DBRead             string
	DBWrite            string
	DBPort             string
	DataCleanUpEnabled bool
}

type store struct {
	db *sql.DB
}

// In order this job to work:
// Locally:
//  - you need to pass the type of the deletion (hard || soft)
//  - you need to pass the PATH of the env file to load
// Kubernetes / Docker:
//  - you need to pass the type of the deletion (hard || soft)
//  - you need to add the env vars to the running container
func main() {
	ctx := context.Background()

	deleteType := os.Args[1]
	if len(os.Args) == 3 {
		if err := godotenv.Load(os.Args[2]); err != nil {
			log.Printf("cannot open given .env file: %v", err)
		}
	}

	log.Printf("initializing %s deleting job...", deleteType)

	cfg := setupConfig()
	if *cfg == (config{}) {
		log.Fatalf("failed to fetch configuration")
	}

	if !cfg.DataCleanUpEnabled {
		log.Println("functionality is disabled, exiting...")
		os.Exit(0)
	}

	// setup SQL
	store := setupSQL(cfg)
	defer store.db.Close(ctx)

	// truncate tables || soft delete data
	switch deleteType {
	case softDelete:
		if err := store.softDelete(ctx, tables...); err != nil {
			log.Fatal(err)
		}
	case hardDelete:
		if err := store.hardDelete(ctx, tables...); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Not valid delete type")
	}

	log.Println("all data are successfully removed, exiting...")
	os.Exit(0)
}

func setupConfig() *config {
	cfg := &config{}
	cfg.DBUsername = os.Getenv("MYSQL_USERNAME")
	cfg.DBPassword = os.Getenv("MYSQL_PASS")
	cfg.DBName = os.Getenv("MYSQL_DB")
	cfg.DBRead = os.Getenv("MYSQL_READ")
	cfg.DBWrite = os.Getenv("MYSQL_WRITE")
	cfg.DBPort = os.Getenv("MYSQL_PORT")
	cfg.DataCleanUpEnabled = false
	if os.Getenv("DATA_CLEANUP_ENABLED") == "true" {
		cfg.DataCleanUpEnabled = true
	}

	return cfg
}

func setupSQL(cfg *config) *store {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBWrite,
		cfg.DBPort,
		cfg.DBName,
	)

	conn, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	return &store{db: conn}
}

// hardDelete truncates the tables
func (s *store) hardDelete(ctx context.Context, tables ...string) error {
	q := "TRUNCATE TABLE %s"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction")
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

// softDelete soft deletes each months data
func (s *store) softDelete(ctx context.Context, tables ...string) error {
	q := "UPDATE %s " +
		"SET deleted_at = NOW() " +
		"WHERE MONTH(date) = MONTH(NOW() - INTERVAL 1 MONTH)"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction")
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
