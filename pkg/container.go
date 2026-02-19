package sqltest

import (
	"context"
	"fmt"
	"testing"

	"github.com/joseCarlosAndrade/go-sqltest/internal/container"
	"github.com/joseCarlosAndrade/go-sqltest/internal/storage"
	"github.com/joseCarlosAndrade/go-sqltest/pkg/populate"
)

type SeedStrategy uint

const (
	// PopulateEmpty wont seed the database at all. usefull for testing insertion operations
	// default option when WithPopulateData is not set
	PopulateEmpty SeedStrategy = iota

	// PopulateCustom will seed the database with the custom populate data configured with
	//   sqltest.WithPopulateData()
	PopulateCustom

	// PopulateScript will seed the database with the custom data defined in the script
	//   sqltest.WithPopulateScript()
	PopulateScript
)

type TestInstance struct {
	connString string // connString to database instance

	driver    DBType
	container container.Wrapper
	config *container.AccessConfig

	seedStategy    SeedStrategy
	populateData   []populate.PopulateTable
	populateScript string

	schemaScript string
	storage      storage.Storage
}

func (ti *TestInstance) SetupTest(ctx context.Context, t *testing.T) error {
	// initializes the container
	err := ti.setupContainer(ctx, t)
	if err != nil {
		return fmt.Errorf("setup test failed: %w", err)
	}

	// generate a client instance
	err = ti.setupStorageConnection(ctx, t)
	if err != nil {
		return fmt.Errorf("could not connect to container: %w", err)
	}

	return nil
}

func (ti *TestInstance) CleanUp(ctx context.Context, t *testing.T) {
	
	// terminate storage
	if ti.storage != nil {
		if err := ti.storage.Close(ctx); err != nil {
			t.Logf("failed to cleanup storage connection: %s", err.Error())
		}
	}

	// terminate container
	if ti.container != nil {
		if err := ti.container.Cleanup(ctx, t); err != nil {
			t.Logf("failed to cleanup container: %s", err.Error())
		}
	}
}

// SetupContainer setups a sql container with the selected driver
func (ti *TestInstance) setupContainer(ctx context.Context, t *testing.T) error {
	switch ti.driver {
	case MySQL:
		return ti.setupMySQL(ctx, t)

	case Postgres:
		return fmt.Errorf("postgres driver not implemented yet")
	default:
		return fmt.Errorf("invalid driver type")
	}

}

func (ti *TestInstance) setupMySQL(ctx context.Context, t *testing.T) error {
	// set mysql container

	container := container.NewMySQLContainer()
	container.Config(ctx, ti.config, ti.schemaScript)
	if err := container.Start(ctx, t); err != nil {
		return fmt.Errorf("error starting container: %w", err)
	}

	// set mysql storage TODO
	return nil
}

func (ti *TestInstance) setupStorageConnection(ctx context.Context, t *testing.T) error {
	// TODO: GET CONFIG TO SEE WHICH DIALECT IS USED

	return nil
}

func (ti *TestInstance) GetConnectionString(ctx context.Context) (string, error) {
	return ti.container.GetConnString(ctx)
}
