package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/dotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

const dialect = "pgx"

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	root  = flags.String("root", ".", "repo root containing service/<svc>/database/migration")
	svc   = flags.String("service", "", "restrict the command to one service (required for 'create')")
)

// contextSpec pairs a bounded-context cluster (its DB_<CONTEXT>_* env prefix)
// with the services whose goose migration directories target that context's
// database. One context = one physical PostgreSQL instance (fronted by its own
// PgBouncer), so every service here shares one goose_db_version table.
type contextSpec struct {
	cluster  string
	services []string
}

// contexts lists the six bounded-context databases in dependency order. It
// mirrors the service → context map in SUMMARY_MONGODB.md §2 and the
// DB_<CONTEXT> prefixes defined in pkg/database/names.go.
var contexts = []contextSpec{
	{database.IdentityCluster, []string{"auth", "user", "role"}},
	{database.MerchantCluster, []string{"merchant", "merchant_award", "merchant_business", "merchant_detail", "merchant_policy"}},
	{database.CatalogCluster, []string{"category", "product"}},
	{database.SalesCluster, []string{"order", "order_item", "transaction"}},
	{database.ExperienceCluster, []string{"cart", "shipping_address", "banner", "slider", "review", "review_detail"}},
	{database.EmailCluster, []string{"email"}},
}

func main() {
	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}
	command := args[0]

	// Authoring a migration is a filesystem-only operation that must land in
	// exactly one service's directory — never looped across contexts.
	if command == "create" {
		if *svc == "" {
			log.Fatalf("'create' requires -service <name> so the migration lands in one service's migration dir")
		}
		if err := goose.RunContext(context.Background(), command, nil, migrationDir(*svc), args[1:]...); err != nil {
			log.Fatalf("Failed to create migration for %s: %v", *svc, err)
		}
		return
	}

	if err := dotenv.Viper(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	matched := false
	for _, c := range contexts {
		services := c.services
		if *svc != "" {
			if !contains(services, *svc) {
				continue
			}
			services = []string{*svc}
		}
		matched = true

		if err := runContext(c.cluster, services, command, args[1:]...); err != nil {
			log.Fatalf("Migration failed for %s: %v", c.cluster, err)
		}
	}

	if !matched {
		log.Fatalf("Unknown service %q: no bounded context owns it", *svc)
	}
}

// runContext applies command to the context's database over a single
// connection. The context's per-service migration directories are staged into
// one temp dir first so goose applies them in global version order: goose
// refuses to run a migration whose version is below the context's current max,
// and services sharing a context have interleaved version ranges (e.g. in the
// identity context, user's 20250212095428 is lower than auth's 20250212095807).
func runContext(cluster string, services []string, command string, extra ...string) error {
	host, port, dbname, user, password, err := connectionParams(cluster)
	if err != nil {
		return err
	}

	dir, err := stageMigrations(cluster, services)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)

	db, err := goose.OpenDBWithDriver(dialect, connStr)
	if err != nil {
		return fmt.Errorf("open %s (%s): %w", dbname, cluster, err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("[%s] error closing database: %v", cluster, cerr)
		}
	}()

	log.Printf("[%s] goose %s (%d services)", cluster, command, len(services))
	if err := goose.RunContext(context.Background(), command, db, dir, extra...); err != nil {
		return fmt.Errorf("[%s] goose %s: %w", cluster, command, err)
	}
	return nil
}

// stageMigrations copies every migration file owned by the context's services
// into a single temp directory. Files are keyed by base name (which is the
// goose version), so the byte-identical outbox_events migration that order and
// transaction both ship collapses to one file instead of tripping goose's
// duplicate-version check.
func stageMigrations(cluster string, services []string) (string, error) {
	dir, err := os.MkdirTemp("", "migrate-"+cluster+"-")
	if err != nil {
		return "", fmt.Errorf("create staging dir: %w", err)
	}

	for _, service := range services {
		src := migrationDir(service)
		matches, err := filepath.Glob(filepath.Join(src, "*.sql"))
		if err != nil {
			return "", fmt.Errorf("[%s] glob %s: %w", cluster, src, err)
		}
		if len(matches) == 0 {
			return "", fmt.Errorf("[%s] no migrations found under %s", cluster, src)
		}
		for _, m := range matches {
			data, err := os.ReadFile(m)
			if err != nil {
				return "", fmt.Errorf("[%s] read %s: %w", cluster, m, err)
			}
			if err := os.WriteFile(filepath.Join(dir, filepath.Base(m)), data, 0o644); err != nil {
				return "", fmt.Errorf("[%s] stage %s: %w", cluster, filepath.Base(m), err)
			}
		}
	}
	return dir, nil
}

// connectionParams reads the mandatory per-context connection keys. Host, port
// and dbname are required: the six databases are separate instances, so there
// is no generic DB_HOST/DB_PORT/DB_NAME fallback. Only the credentials fall
// back to the base DB_USERNAME/DB_PASSWORD keys.
func connectionParams(cluster string) (host, port, dbname, user, password string, err error) {
	host = viper.GetString(cluster + "_HOST")
	port = viper.GetString(cluster + "_PORT")
	dbname = viper.GetString(cluster + "_NAME")
	user = viper.GetString(cluster + "_USERNAME")
	password = viper.GetString(cluster + "_PASSWORD")

	if host == "" || port == "" || dbname == "" {
		return "", "", "", "", "", fmt.Errorf(
			"%s_HOST, %s_PORT and %s_NAME must be set (no generic DB_HOST/DB_PORT/DB_NAME fallback)",
			cluster, cluster, cluster)
	}
	if user == "" {
		user = viper.GetString("DB_USERNAME")
	}
	if password == "" {
		password = viper.GetString("DB_PASSWORD")
	}
	return host, port, dbname, user, password, nil
}

func migrationDir(service string) string {
	return filepath.Join(*root, "service", service, "database", "migration")
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Runs goose against the six bounded-context PostgreSQL databases (identity,
merchant, catalog, sales, experience, email) in dependency order, reading the
DB_<CONTEXT>_* keys from the environment. Pass -service to target one service.
Examples:
    migrate up
    migrate status
    migrate -service auth create add_users_email sql
`

	usageCommands = `
Commands:
    up                   Migrate every context to the most recent version available
    up-by-one            Migrate every context up by 1
    up-to VERSION        Migrate every context to a specific VERSION
    down                 Roll back every context by 1
    down-to VERSION      Roll back every context to a specific VERSION
    redo                 Re-run the latest migration on every context
    reset                Roll back all migrations on every context
    status               Dump the migration status for every context
    version              Print the current version of every context
    create NAME [sql|go] Create a migration in -service's directory
    fix                  Apply sequential ordering to every context`
)
