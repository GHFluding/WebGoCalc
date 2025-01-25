package integration_test

import (
	"context"
	"fmt"
	"log"
	"test/internal/config"
	"test/internal/database/postgres"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestDBTable(t *testing.T) {
	// config load
	cfg := config.MustLoad()
	log.Printf("Config loaded: %+v\n", cfg)

	// connect to db
	dbpool, err := postgres.Connect(*cfg)
	if err != nil {
		assert.Error(t, err, "Failed to connect to the database")
		t.FailNow()
	}
	defer dbpool.Close()
	log.Println("Connected to the database successfully")

	// setting migration
	migrationPath := "file:///test/migrations" // Проверьте, что путь корректен
	log.Printf("Initializing migrations from: %s\n", migrationPath)
	m, err := migrate.New(
		migrationPath,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Storage.User,
			cfg.Storage.Password,
			"webgocalc_postgres",
			cfg.Storage.Port,
			cfg.Storage.DBName,
		),
	)
	if err != nil {
		assert.Error(t, err, "Failed to initialize migrations")
		t.FailNow()
	}
	defer func() {
		if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
			log.Printf("Migration close errors: sourceErr=%v, dbErr=%v", sourceErr, dbErr)
		}
	}()
	log.Println("Migration initialized successfully")

	// use migration
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		assert.Error(t, err, "Failed to apply migrations")
		t.FailNow()
	}
	log.Println("Migrations applied successfully!")

	// Check table
	tables := []string{"students", "calendar"}
	for _, table := range tables {
		t.Run("Check table "+table, func(t *testing.T) {
			log.Printf("Checking table existence: %s\n", table)
			exists, err := checkTableExists(context.Background(), dbpool, table)
			if err != nil {
				log.Printf("Error checking table %s: %v\n", table, err)
			}
			assert.NoError(t, err, "Error occurred while checking table existence")
			assert.True(t, exists, "Table does not exist: "+table)
		})
	}
}

// Func for checking tables
func checkTableExists(ctx context.Context, dbpool *pgxpool.Pool, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = $1
		)
	`
	log.Printf("Executing query to check table: %s\nQuery: %s\n", tableName, query)
	err := dbpool.QueryRow(ctx, query, tableName).Scan(&exists)
	if err != nil {
		log.Printf("Error executing query for table %s: %v\n", tableName, err)
	}
	return exists, err
}
