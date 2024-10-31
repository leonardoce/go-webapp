package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/viper"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate execute the migrations on the connected database
func Migrate(ctx context.Context) error {
	migrationsSource, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("while loading embedded migrations: %w", err)
	}

	db, err := sql.Open("pgx", viper.GetString("connection-string"))
	if err != nil {
		return fmt.Errorf("while creating database connection pool: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error while closing the connection pool", err)
		}
	}()

	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("while creating database connection: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error while closing the connection", err)
		}
	}()

	pgDriver, err := postgres.WithConnection(
		ctx,
		conn,
		&postgres.Config{})
	if err != nil {
		return fmt.Errorf("while creating database migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"embedded",
		migrationsSource,
		"pg",
		pgDriver,
	)
	if err != nil {
		return fmt.Errorf("while starting migrations: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("while executing migrations: %w", err)
	}

	return nil
}
