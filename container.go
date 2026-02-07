package sqltest

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/orasis-holding/pricing-go-swiss-army-lib/cpool"
	"github.com/orasis-holding/pricing-go-swiss-army-lib/nexus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type PopulateTable struct {
	Name    string
	Columns []string
	Data    [][]string
}

type TestInstance struct {
	connString        string
	migrationVersion  uint
	container         *mysql.MySQLContainer
	shouldMockNexus   bool
	nexusCompanyID    string
	nexusInstance     *nexus.Nexus
	shouldReturnCpool bool
	cpoolInstance     *cpool.Factory
	shouldMockPricing bool

	populateData []PopulateTable
}

func NewPopulate(tableName string, columns ...string) *PopulateTable {
	// columns := make([]string, len(columns))
	return &PopulateTable{
		Name:    tableName,
		Columns: columns,
		Data:    make([][]string, 0),
	}
}

func (p *PopulateTable) Insert(values ...string) *PopulateTable {
	p.Data = append(p.Data, values)
	return p
}

// SetupContainer setups a mysql container with the defaul schema in schema/schema.sql. setps up the instance
func (ti *TestInstance) SetupContainer(ctx context.Context, t *testing.T) error {

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("pricingtest"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		mysql.WithScripts(filepath.Join("schema", "schema.sql")))

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			t.Logf("failed to termiante container: %s", err.Error())
		}
	})

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

	return nil
}

// func (ti *TestInstance)

func (ti *TestInstance) SetupTest(ctx context.Context, t *testing.T) (*nexus.Nexus, error) {
	if ti.container == nil {
		otelzap.L().ErrorContext(ctx, "nexus setup failed: container not initialized")

		return nil, fmt.Errorf("container not initialized")
	}

	// TODO: should populate the nexus db inside the container with mock data
	connString := ti.connString // todo: insert this in nexus

	nn, err := nexus.NewNexusInstance(ctx, connString, "secret", true)
	if err != nil {
		t.Logf("could not initialize nexus instance: %s", err.Error())
		return nil, err
	}

	ti.nexusInstance = nn

	// returns the nexus instance instantiated

	return ti.nexusInstance, nil
}
