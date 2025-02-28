package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(runningEnv string) (*pgxpool.Pool, error) {
	log.Println("Opening connection to database...")

	var connStr string
	if runningEnv == "docker" {
		connStr = os.Getenv("DB_URL")
	} else {
		connStr = os.Getenv("DB_URL_LOCALHOST")
	}

	if connStr == "" {
		return nil, fmt.Errorf("DB_URL environment variable not set")
	}

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("connected to database")

	if err := migrateDB(dbpool); err != nil {
		return nil, err
	}

	return dbpool, nil
}

func migrateDB(dbpool *pgxpool.Pool) error {
	db, err := sql.Open("pgx", dbpool.Config().ConnString())
	if err != nil {
		return fmt.Errorf("opening sql.DB from pgxpool: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("creating postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("creating migration instance: %w", err)
	}

	defer m.Close()

	// Run migrations up to the latest version
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("running migrations: %w", err)
		}
		log.Println("No migrations to apply.")
	} else {
		log.Println("Migrations completed successfully.")
	}

	return nil
}
