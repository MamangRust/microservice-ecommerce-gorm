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

// NewGormClient connects to the base database using the generic DB_* keys
// and returns a *gorm.DB instance.
func NewGormClient(logger logger.LoggerInterface) (*gorm.DB, error) {
	return NewGormClientWithPrefix(logger, "DB")
}

// NewGormClientWithPrefix connects to the database configured via the given
// prefix keys (e.g. DB_ORDER_HOST, DB_ORDER_NAME) with fallback to the base
// DB_* keys. Each microservice uses its own prefix so it talks exclusively
// to its own PostgreSQL instance.
func NewGormClientWithPrefix(logger logger.LoggerInterface, prefix string) (*gorm.DB, error) {
	if prefix == "" {
		prefix = "DB"
	}

	dbDriver := viper.GetString(fmt.Sprintf("%s_DRIVER", prefix))
	if dbDriver == "" {
		dbDriver = viper.GetString("DB_DRIVER")
	}

	if dbDriver != "postgres" && dbDriver != "pgx" {
		logger.Error("gorm postgres driver only supports PostgreSQL", zap.String("DB_DRIVER", dbDriver))
		return nil, fmt.Errorf("gorm postgres driver only supports PostgreSQL, got: %s", dbDriver)
	}

	hostKey := fmt.Sprintf("%s_HOST", prefix)
	portKey := fmt.Sprintf("%s_PORT", prefix)
	userKey := fmt.Sprintf("%s_USERNAME", prefix)
	nameKey := fmt.Sprintf("%s_NAME", prefix)
	passKey := fmt.Sprintf("%s_PASSWORD", prefix)

	host := viper.GetString(hostKey)
	if host == "" {
		host = viper.GetString("DB_HOST")
	}
	port := viper.GetString(portKey)
	if port == "" {
		port = viper.GetString("DB_PORT")
	}
	user := viper.GetString(userKey)
	if user == "" {
		user = viper.GetString("DB_USERNAME")
	}
	dbname := viper.GetString(nameKey)
	if dbname == "" {
		dbname = viper.GetString("DB_NAME")
	}
	password := viper.GetString(passKey)
	if password == "" {
		password = viper.GetString("DB_PASSWORD")
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
