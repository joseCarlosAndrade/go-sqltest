package sqltest

import (
	"context"
	
	
	// "github.com/orasis-holding/pricing-go-swiss-army-lib/cpool"
	// "github.com/orasis-holding/pricing-go-swiss-army-lib/nexus"
	// mysqlcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type TestOption func(*TestInstance)

// WithNexusConfig enables nexus mocking
//
// This initializes the nexus schema in the mysql container, but using the nexus testing mode, which
// skips all rsa decryption operations. 
//
// The companyId passed as parameter will be inserted in the companies table
// and a nexus Instance will be returned in the setup functions.
//
// Querying for that companyID will result in a connection to the mysql container pricing instance
func WithNexus(companyId string) TestOption {
	return func(t *TestInstance) {
		t.shouldMockNexus = true
		t.nexusCompanyID = companyId
	}
}

// WithNexusConfigData provides options to populate the config table in nexus
// 
// Use only when needed
func WithNexusConfigData(configs ...PopulateTable) TestOption {
	return func(t *TestInstance) {
		t.nexusConfigData = append(t.nexusConfigData, configs...)
	}
}

// WithCpool
// func WithCpool() TestOption {
// 	return func(t *TestInstance) {
// 		t.shouldReturnCpool = true
// 	}
// }

// WithPricing instantiates pricing db 
//
// The seed data used will be default for the latest 
// schema unless specified otherwise  the options
//   WithMigrationVersion() and WithPopulateData()
func WithPricing() TestOption {
	return func(t *TestInstance) {
		t.shouldMockPricing = true
	}
}

// WithMigrationVersion sets the migration version used to instantiate the pricing db. 
// 
// If not set, this will use the latest available
func WithMigrationVersion(version uint) TestOption {
	return func(t *TestInstance) {
		t.migrationVersion = version
	}
}

// WithPopulateData sets a custom seed data to populate the pricing mysql container instance
//
// If not set, this will seed default data specified at: (TODO - NOT IMPLEMENTED YET)
func WithPopulateData(data ...PopulateTable) TestOption {
	return func(t *TestInstance) {
		t.pricingPopulateData = append(t.pricingPopulateData, data...)
	}
}

func NewSQLTest(ctx context.Context, configs ...TestOption) (*TestInstance, error) {
	instance := TestInstance{
		migrationVersion:    0,     // default: use latest
		shouldMockNexus:     false, // default: explicit false
		shouldMockPricing:   false,
		pricingPopulateData: make([]PopulateTable, 0),
		nexusConfigData:     make([]PopulateTable, 0),
	}

	for _, config := range configs {
		config(&instance)
	}

	return &instance, nil
}
