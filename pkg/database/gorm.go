package database

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewGormClientWithPrefix connects to the database configured via the given
// bounded-context prefix (e.g. DB_SALES_HOST, DB_SALES_NAME). Each service
// passes its context's prefix so it talks exclusively to that context's
// PostgreSQL instance (fronted by its own PgBouncer).
//
// Host, port and dbname are mandatory per context: there is deliberately no
// generic DB_HOST/DB_PORT/DB_NAME fallback, which would silently collapse
// every bounded context onto one shared database. Only the credentials fall
// back to the base DB_USERNAME/DB_PASSWORD keys.
func NewGormClientWithPrefix(logger logger.LoggerInterface, prefix string) (*gorm.DB, error) {
	if prefix == "" {
		return nil, fmt.Errorf("database cluster prefix must not be empty")
	}

	dbDriver := viper.GetString(fmt.Sprintf("%s_DRIVER", prefix))
	if dbDriver == "" {
		dbDriver = viper.GetString("DB_DRIVER")
	}

	if dbDriver != "postgres" && dbDriver != "pgx" {
		logger.Error("gorm postgres driver only supports PostgreSQL", zap.String("DB_DRIVER", dbDriver))
		return nil, fmt.Errorf("gorm postgres driver only supports PostgreSQL, got: %s", dbDriver)
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
		err := fmt.Errorf("%s_HOST, %s_PORT and %s_NAME must be set (no generic DB_HOST/DB_PORT/DB_NAME fallback)", prefix, prefix, prefix)
		logger.Error("Incomplete database cluster configuration",
			zap.String("prefix", prefix),
			zap.Error(err),
		)
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password,
	)

	maxOpenConns := viper.GetInt("DB_MAX_OPEN_CONNS")
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}

	maxIdleConns := viper.GetInt("DB_MIN_IDLE_CONNS")
	if maxIdleConns <= 0 {
		maxIdleConns = 50
	}

	connMaxLifetime := viper.GetDuration("DB_CONN_MAX_LIFETIME")
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour
	}

	connMaxIdleTime := viper.GetDuration("DB_CONN_MAX_IDLE_TIME")
	if connMaxIdleTime == 0 {
		connMaxIdleTime = 30 * time.Minute
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})
	if err != nil {
		logger.Error("Failed to connect to database via GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database via GORM: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Error("Failed to get underlying sql.DB from GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Error("Failed to ping database via GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database via GORM: %w", err)
	}

	logger.Debug("GORM database connection established successfully",
		zap.String("prefix", prefix),
		zap.String("dbname", dbname),
		zap.Int("MaxOpenConns", maxOpenConns),
		zap.Int("MaxIdleConns", maxIdleConns),
		zap.Duration("ConnMaxLifetime", connMaxLifetime),
		zap.Duration("ConnMaxIdleTime", connMaxIdleTime),
	)

	return gormDB, nil
}
