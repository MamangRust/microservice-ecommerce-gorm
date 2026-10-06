package tests

import (
	"context"
	"time"

	chDriver "github.com/ClickHouse/clickhouse-go/v2"
	pkgclickhouse "github.com/MamangRust/microservice-ecommerce-pkg/clickhouse"
	"github.com/spf13/viper"
	tcclickhouse "github.com/testcontainers/testcontainers-go/modules/clickhouse"
)

// SetupClickHouse starts a throwaway ClickHouse container, publishes its
// connection settings through viper (the keys pkg/clickhouse reads) and makes
// sure CLICKHOUSE_DATABASE exists before any suite opens a client. Suites that
// read from ClickHouse call it right after BaseTestSuite.SetupSuite; the
// container is terminated by TearDownSuite.
func (s *BaseTestSuite) SetupClickHouse() {
	s.Require().Nil(s.ts.ClickHouseContainer, "clickhouse is already running for this suite")

	const (
		username = "testuser"
		password = "testpass"
		database = "testdb"
	)

	container, err := tcclickhouse.Run(s.Ctx,
		"clickhouse/clickhouse-server:24-alpine",
		tcclickhouse.WithDatabase(database),
		tcclickhouse.WithUsername(username),
		tcclickhouse.WithPassword(password),
	)
	s.Require().NoError(err, "failed to start clickhouse container")
	s.ts.ClickHouseContainer = container

	host, err := container.ConnectionHost(s.Ctx)
	s.Require().NoError(err, "failed to resolve clickhouse host")

	viper.Set("CLICKHOUSE_ADDR", host)
	viper.Set("CLICKHOUSE_HOST", host)
	viper.Set("CLICKHOUSE_PORT", "9000")
	viper.Set("CLICKHOUSE_DATABASE", database)
	viper.Set("CLICKHOUSE_USERNAME", username)
	viper.Set("CLICKHOUSE_PASSWORD", password)

	// NewClient pings with CLICKHOUSE_DATABASE as the default database, so the
	// database has to be reachable (and therefore exist) before the suite opens
	// its connection.
	s.Require().NoError(waitClickHouseReady(s.Ctx, host, username, password, database, 30*time.Second))
	s.Require().NoError(pkgclickhouse.EnsureDatabase(s.Log), "failed to create clickhouse database")
}

// waitClickHouseReady polls the native protocol until the container answers a
// trivial ping, so a slow container start does not surface as a confusing
// failure from pkg/clickhouse.NewClient.
func waitClickHouseReady(ctx context.Context, addr, username, password, database string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	var lastErr error
	for time.Now().Before(deadline) {
		conn, err := chDriver.Open(&chDriver.Options{
			Addr: []string{addr},
			Auth: chDriver.Auth{
				Database: database,
				Username: username,
				Password: password,
			},
			DialTimeout: 5 * time.Second,
		})
		if err == nil {
			lastErr = conn.Ping(ctx)
			conn.Close()
			if lastErr == nil {
				return nil
			}
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}

	return lastErr
}
