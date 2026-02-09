package sqltest

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/joseCarlosAndrade/go-sqltest/storage"
	"github.com/orasis-holding/pricing-go-swiss-army-lib/nexus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type PopulateTable struct {
	Name    string
	Columns []string
	Data    [][]any
}

type SeedStrategy uint

const (
	DBName string = "pricingtest"
)

const (
	// PopulateEmpty wont seed the pricing database at all. usefull for testing insertion operations
	PopulateEmpty SeedStrategy = iota

	// PopulateDefault will seed the pricing database with the default populate script defined at: TODO
	PopulateDefault

	// PopulateCustom will seed the pricing database with the custom populate data configured with
	//   sqltest.WithPopulateData()
	PopulateCustom
)

type TestInstance struct {
	connString string // connString to mysql instance
	// connStringPricing string

	migrationVersion uint

	container *mysql.MySQLContainer

	shouldMockNexus bool
	nexusCompanyID  string
	nexusConfigData []PopulateTable // populates the config table in nexus (use only when needed)
	nexusInstance   *nexus.Nexus
	// shouldReturnCpool bool
	// cpoolInstance     *cpool.Factory

	shouldMockPricing   bool
	seedStategy         SeedStrategy
	pricingPopulateData []PopulateTable
	storage *storage.Storage
}

func NewPopulate(tableName string, columns ...string) *PopulateTable {
	// columns := make([]string, len(columns))
	return &PopulateTable{
		Name:    tableName,
		Columns: columns,
		Data:    make([][]any, 0),
	}
}

func (p *PopulateTable) Insert(values ...any) *PopulateTable { // test: maybe values ...any ?
	p.Data = append(p.Data, values)
	return p
}

func (ti *TestInstance) SetupTest(ctx context.Context, t *testing.T) error {
	err := ti.setupContainer(ctx, t)
	if err != nil {
		return fmt.Errorf("setup test failed: %w", err)
	}

	// generate a client instance
	err = ti.setupStorageConnection(ctx, t )
	if err != nil {
		return fmt.Errorf("could not connect to container: %w", err)
	}

	if ti.shouldMockNexus {
		// init nexus instance
		nn, err := ti.initNexusInstance(ctx, t)
		if err != nil {
			return fmt.Errorf("init nexus instance failed: %w", err)
		}

		ti.nexusInstance = nn

		// populate nexus
		if err := ti.seedNexus(ctx, t); err != nil {
			return fmt.Errorf("seed nexus failed: %w", err)
		}

	}

	if ti.shouldMockPricing {
		// populate pricing

		// init pricing schema
		err := ti.initPricing(ctx, t)
		if err != nil {
			return fmt.Errorf("init pricing instance failed: %w", err)
		}

		// seed pricing data
		if err := ti.seedPricing(ctx, t); err != nil {
			return fmt.Errorf("seed pricing failed: %w", err)
		}
	}

	return nil
}

func (ti * TestInstance) CleanUp(ctx context.Context, t * testing.T) {
	if err := testcontainers.TerminateContainer(ti.container); err !=nil {
		t.Logf("failed to terminate container: %s", err.Error())
	}

	if ti.storage != nil {
		if err := ti.storage.Close(); err != nil {
			t.Logf("failed to cleanup myqsl connection: %s", err.Error())
		}
	}
}

// SetupContainer setups a mysql container with the defaul schema in schema/
func (ti *TestInstance) setupContainer(ctx context.Context, t *testing.T) error {
	schema := "default.sql"

	if ti.shouldMockNexus {
		schema = "schema-nexus.sql"
	}

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase(DBName),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		mysql.WithScripts(filepath.Join("schema", schema)))


	if err != nil {
		t.Logf("failed to start container: %s", err.Error())
		return err
	}

	t.Logf("container successfully initialized")

	ti.container = mysqlContainer
	ti.connString, err = mysqlContainer.ConnectionString(ctx)
	if err != nil {
		t.Logf("could not get mysql connection string: %s", err.Error())
		return err
	}

	time.Sleep(2 * time.Second) // sleep a bit to give proper container init

	return nil
}

func (ti *TestInstance) setupStorageConnection(ctx context.Context, t *testing.T) error {
	storage, err := storage.NewStorage(ctx, ti.connString)
	if err != nil {
		return err
	}

	ti.storage = storage

	return nil
}

func (ti *TestInstance) initNexusInstance(ctx context.Context, t *testing.T) (*nexus.Nexus, error) {
	if ti.container == nil {
		return nil, fmt.Errorf("container not initialized")
	}

	connString := ti.connString // todo: insert this in nexus

	nn, err := nexus.NewNexusInstance(ctx, connString, "secret", true)
	if err != nil {
		t.Logf("could not initialize nexus instance: %s", err.Error())
		return nil, err
	}

	// returns the nexus instance instantiated

	return nn, nil
}

func (ti *TestInstance) seedNexus(ctx context.Context, t *testing.T) error {
	if err := ti.storage.Ping(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	if err := ti.insertNexusCompanyID(ctx); err != nil {

		return err
	}

	if err := ti.populateNexusConfig(ctx, t); err != nil {
		return err
	}

	return nil
}

func (ti *TestInstance) insertNexusCompanyID(ctx context.Context) error {
	return ti.storage.InsertNexusCompany(ctx, ti.nexusCompanyID, ti.connString, DBName)
}

func (ti *TestInstance) populateNexusConfig(ctx context.Context, t *testing.T) error {
	return nil
}

func (ti *TestInstance) initPricing(ctx context.Context, t *testing.T) error {
	// CHECK MIGRATION VERSION
	// 	ANALYSE LATEST ONE IF INVALID

	// CHECK CONTAINER REACHABILITY

	// apply everuthing

	return nil
}

func (ti *TestInstance) seedPricing(ctx context.Context, t *testing.T) error {
	// CHECK CONTAINER REACHABILITY

	// CHECK IF USER PASSED POULATION DATA
	//	IF NOT, USE DEFAULT FOR THIS MIGRATION

	// apply everyuthing

	return nil
}

func (ti *TestInstance) GetNexusInstance(ctx context.Context) (*nexus.Nexus, error) {
	if !ti.shouldMockNexus {
		return nil, fmt.Errorf("nexus mock is set to false")
	}

	if ti.container == nil {
		return nil, fmt.Errorf("container was not initialized")
	}

	if ti.nexusInstance == nil {
		return nil, fmt.Errorf("nexus was not initialized")
	}

	return ti.nexusInstance, nil
}

func (ti *TestInstance) GetConnectionString(ctx context.Context) (string, error) {
	if ti.container == nil {
		return "", fmt.Errorf("container was not initialized")
	}

	if !ti.shouldMockNexus && !ti.shouldMockPricing {
		return "", fmt.Errorf("no mock was initialized")
	}

	if ti.connString == "" {
		return "", fmt.Errorf("initialization was not successful")
	}

	return ti.connString, nil
}

// func (ti* TestInstance) SetupNexus(ctx context.Context, t *testing.T) (*nexus.Nexus, error) {
// 	if !ti.shouldMockNexus {
// 		return nil, fmt.Errorf("nexus mock is set to false")
// 	}

// 	if ti.container == nil {
// 		return nil, fmt.Errorf("container not initialized")
// 	}

// 	nn, err := ti.initNexus(ctx, t)
// 	if err != nil {
// 		return nil, fmt.Errorf("nexus setup failed: %w", err)
// 	}

// 	return nn, nil
// }
