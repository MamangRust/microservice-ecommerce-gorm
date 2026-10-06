package database

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// RunMigrations executes database migrations using goose against the DB
// configured via the given bounded-context prefix (e.g. "DB_SALES").
// Host, port and dbname are mandatory per context — there is deliberately no
// generic DB_HOST/DB_PORT/DB_NAME fallback, so a misconfigured service fails
// fast instead of migrating the wrong instance. Only the credentials fall back
// to the base DB_USERNAME/DB_PASSWORD keys.
// path: directory containing migration files.
func RunMigrations(log logger.LoggerInterface, prefix, path string) error {
	if prefix == "" {
		return fmt.Errorf("database cluster prefix must not be empty")
	}

	host := viper.GetString(fmt.Sprintf("%s_HOST", prefix))
	port := viper.GetString(fmt.Sprintf("%s_PORT", prefix))
	dbname := viper.GetString(fmt.Sprintf("%s_NAME", prefix))

	user := viper.GetString(fmt.Sprintf("%s_USERNAME", prefix))
	if user == "" {
		user = viper.GetString("DB_USERNAME")
	}
	password := viper.GetString(fmt.Sprintf("%s_PASSWORD", prefix))
	if password == "" {
		password = viper.GetString("DB_PASSWORD")
	}

	if host == "" || port == "" || dbname == "" {
		return fmt.Errorf("%s_HOST, %s_PORT and %s_NAME must be set (no generic DB_HOST/DB_PORT/DB_NAME fallback)", prefix, prefix, prefix)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password,
	)

	db, err := goose.OpenDBWithDriver("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database for migrations: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Failed to close database after migrations", zap.Error(err))
		}
	}()

	log.Info("Running database migrations",
		zap.String("path", path),
		zap.String("dbname", dbname),
		zap.String("prefix", prefix),
	)

	if err := goose.RunContext(context.Background(), "up", db, path); err != nil {
		return fmt.Errorf("migration 'up' failed: %w", err)
	}

	log.Info("Database migrations completed successfully")
	return nil
}
