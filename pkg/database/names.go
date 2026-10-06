package database

// Bounded-context PostgreSQL cluster prefixes. Each service sets
// server.Config.DBCluster to the prefix of the context that owns its tables so
// it talks exclusively to that context's PostgreSQL instance (fronted by its
// own PgBouncer). The values match the DB_<CONTEXT>_* keys defined in .env,
// docker-compose and the Kubernetes ConfigMap.
const (
	IdentityCluster   = "DB_IDENTITY"
	CatalogCluster    = "DB_CATALOG"
	MerchantCluster   = "DB_MERCHANT"
	SalesCluster      = "DB_SALES"
	ExperienceCluster = "DB_EXPERIENCE"
	EmailCluster      = "DB_EMAIL"
)

// Logical database names carried by each context instance (POSTGRES_DB /
// DB_<CONTEXT>_NAME). Exposed here so env, seeder and tests share one source
// of truth.
const (
	IdentityDB   = "ec_identity"
	CatalogDB    = "ec_catalog"
	MerchantDB   = "ec_merchant"
	SalesDB      = "ec_sales"
	ExperienceDB = "ec_experience"
	EmailDB      = "ec_email"
)
