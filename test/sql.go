package test

import (
	"context"

	sqltest "github.com/joseCarlosAndrade/go-sqltest"
)

func test() {
	ctx := context.Background()

	populateProduct := sqltest.NewPopulate("product", "id", "description", "price").
			Insert("10", "COCA COLA 2L", "11.50").
			Insert("12", "TODDY 200G", "7.80")

	populateHierarchy := sqltest.NewPopulate("hierarchy", "id", "type", "parent").
			Insert("123", "1", "555634")

	sqltest.NewSQLTest(ctx, 
		sqltest.WithNexus("123"),
		sqltest.WithMigrationVersion(9),
		sqltest.WithPricing(),
		sqltest.WithPopulateData(*populateProduct, *populateHierarchy))

		// todo: maybe use the populate like i did in readme.
		// instead of directly passing populate product, we pass a function to WithPopulateData
		// this function calls everything.
		// this allows me to pass custom functions and put my populate inside,
		// or use predefined functions

		// WithPopulateData(sqltest.PopulateDefault)
		// WithPopulateData(sqltest.PopulateEmpty)
		// WithPopulateData(sqltest.PopulateCustom(populateProduct, populateHierarchy))
}