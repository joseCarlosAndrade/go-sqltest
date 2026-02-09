package test

import (
	"context"
	"testing"

	sqltest "github.com/joseCarlosAndrade/go-sqltest"
)

func test()  {
	ctx := context.Background()
	t := &testing.T{}

	populateProduct := sqltest.NewPopulate("product", "id", "description", "price").
			Insert("10", "COCA COLA 2L", "11.50").
			Insert("12", "TODDY 200G", "7.80")

	populateHierarchy := sqltest.NewPopulate("hierarchy", "id", "type", "parent").
			Insert("123", "1", "555634")

	instance, err := sqltest.NewSQLTest(ctx, 
		sqltest.WithNexus("123"),
		sqltest.WithMigrationVersion(9),
		sqltest.WithPricing(),
		sqltest.WithSeedStrategy(sqltest.PopulateCustom),
		sqltest.WithPopulateData(*populateProduct, *populateHierarchy))

	if err != nil {
		panic(err)
	}

	// sets up with proper
	if err := instance.SetupTest(ctx, t); err != nil {
		panic(err)
	}
	defer instance.CleanUp(ctx, t)

	nexusInstance, err := instance.GetNexusInstance(ctx)
	if err != nil {
		panic(err)
	}

	connString, err := instance.GetConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	 _ = nexusInstance
	 _ = connString
}