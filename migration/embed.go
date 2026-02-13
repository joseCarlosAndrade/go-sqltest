package migration


import "embed"

//go:embed pricing-client-sql-schema/migrations/*.sql
var MigrationFS embed.FS