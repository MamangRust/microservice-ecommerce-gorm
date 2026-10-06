package server

import (
	"time"
)

// Config defines the common configuration for all gRPC services
type Config struct {
	ServiceName          string
	ServiceVersion       string
	Environment          string
	Port                 int
	OtelEndpoint         string
	OtelSamplingFraction float64

	// DBCluster is the env prefix used to resolve this service's bounded-context
	// PostgreSQL connection (e.g. "DB_SALES" -> DB_SALES_HOST, DB_SALES_NAME).
	// It must be one of the pkg/database cluster constants; an empty or
	// incomplete prefix fails fast at startup rather than falling back to a
	// shared database.
	DBCluster string

	// MigrationPath is the directory (relative to the service workdir) that
	// holds this service's goose migrations. When set, migrations are applied
	// automatically at startup.
	MigrationPath string
}

// Default constants for gRPC server
const (
	DefaultMaxConcurrentConn = 1024
	DefaultWindowSize        = 16 * 1024 * 1024
	DefaultKeepaliveTime     = 20 * time.Second
	DefaultKeepaliveTimeout  = 5 * time.Second
	DefaultMinKeepaliveTime  = 5 * time.Second

	MonitoringInterval     = 30 * time.Second
	CleanupInterval        = 120 * time.Second
	CacheRefCountThreshold = 500

	ShutdownTimeout = 30 * time.Second

	RedisDialTimeout  = 5 * time.Second
	RedisReadTimeout  = 3 * time.Second
	RedisWriteTimeout = 3 * time.Second
	RedisPoolSize     = 10
	RedisMinIdleConns = 3
)
