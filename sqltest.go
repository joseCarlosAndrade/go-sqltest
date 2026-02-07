package sqltest

import (
	"context"
	// "github.com/orasis-holding/pricing-go-swiss-army-lib/cpool"
	// "github.com/orasis-holding/pricing-go-swiss-army-lib/nexus"
	// mysqlcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type TestOption func(*TestInstance)

// WithNexusConfig enables nexus mock
func WithNexus(companyId string) TestOption {
	return func(t *TestInstance) {
		t.shouldMockNexus = true
		t.nexusCompanyID = companyId
	}
}

// WithCpool
func WithCpool() TestOption {
	return func(t *TestInstance) {
		t.shouldReturnCpool = true
	}
}

// WithPricing instantiates pricing db with the selected migration version
func WithPricing() TestOption {
	return func(t *TestInstance) {
		t.shouldMockPricing = true
	}
}

// WithMigrationVersion sets the migration version used to instantiate the pricing db. If not set, this will use the latest
func WithMigrationVersion(version uint) TestOption {
	return func(t *TestInstance) {
		t.migrationVersion = version
	}
}

func WithPopulateData(data ...PopulateTable) TestOption {
	return func(t *TestInstance) {
		t.populateData = append(t.populateData, data...)
	}
}

func NewTest(ctx context.Context, configs ...TestOption) (*TestInstance, error) {
	instance := TestInstance{
		migrationVersion:  0,     // default: use latest
		shouldMockNexus:   false, // default: explicit false
		shouldReturnCpool: false,
		shouldMockPricing: false,
		populateData: make([]PopulateTable, 0),
	}

	for _, config := range configs {
		config(&instance)
	}

	return &instance, nil
}
