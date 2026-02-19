package sqltest

import (
	"context"

	"github.com/joseCarlosAndrade/go-sqltest/internal/container"
	"github.com/joseCarlosAndrade/go-sqltest/pkg/populate"
)

type TestOption func(*TestInstance)

type DBType uint

const (
	MySQL    = iota
	Postgres // not supported yet

)

// WithSeedStrategy sets the strategy for seeding the pricing instance of mysql container
//
// Strategies accepted so far:
//
//	PopulateEmpty - leaves the mysql with no inserted data. Useful when testing only insertion operations
//	PopulateDefault - will seed the database based on a default script for the configured migration version defined in TODO
//	PopulateCustom - will seed the database using the custom options configuired int WithPopulateData()
func WithSeedStrategy(strategy SeedStrategy) TestOption {
	return func(t *TestInstance) {
		t.seedStategy = strategy
	}
}

// WithPopulateData sets a custom seed data to populate the pricing mysql container instance
//
// Note: should only be used when:
//
//	WithSeedStrategy() is also set to PopulateCustom: (TODO - NOT IMPLEMENTED YET)
func WithPopulateData(data ...populate.PopulateTable) TestOption {
	return func(t *TestInstance) {
		t.populateData = append(t.populateData, data...)
	}
}

// WithPopulateScripts sets a path to some script.sql responsible for the custom seeding
//
// Note: should only be used when:
//
//	WithSeedStrategy() is also set to PopulateScript
func WithPopulateScript(script string) TestOption {
	return func(t *TestInstance) {
		t.populateScript = script
	}
}

// WithSchemaScript sets a script.sql to be executed when the container is instantiated. It's responsible for setting up the database schema
// if not set, the database will start empty
func WithSchemaScript(script string) TestOption {
	return func(t *TestInstance) {
		t.schemaScript = script
	}
}

// WithCredentials sets custom credentials. if not set, default credentials will be used: 
//   db, root, password
func WithCredentials(dbName, username, password string) TestOption {
	return func(t *TestInstance) {
		t.config.DBName = dbName
		t.config.Username = username
		t.config.Password = password
	}
}

func NewSQLTest(ctx context.Context, driver DBType, configs ...TestOption) (*TestInstance, error) {
	instance := TestInstance{
		populateData: make([]populate.PopulateTable, 0),
		seedStategy:  PopulateEmpty,
		driver:       driver,
		config: &container.AccessConfig{
			DBName:   "db",
			Username: "root",
			Password: "password",
		},
	}

	for _, config := range configs {
		config(&instance)
	}

	return &instance, nil
}
