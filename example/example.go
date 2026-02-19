package test

import (
	"context"
	"testing"

	sqltest "github.com/joseCarlosAndrade/go-sqltest/pkg"
	"github.com/joseCarlosAndrade/go-sqltest/pkg/populate"
)

func test()  {
	ctx := context.Background()
	t := &testing.T{}

	populateProduct := populate.New("product", "id", "description", "price").
			Insert("10", "COCA COLA 2L", "11.50").
			Insert("12", "TODDY 200G", "7.80")

	populateUser := populate.New("id", "role", "email", "age").
			Insert("123", "manager", "user@mail", "26")

	instance, err := sqltest.NewSQLTest(ctx, 
		sqltest.MySQL,
		sqltest.WithSchemaScript("schema.sql"),
		sqltest.WithCredentials("customdb", "myuser", "password"),
		sqltest.WithSeedStrategy(sqltest.PopulateCustom),
		sqltest.WithPopulateData(*populateProduct, *populateUser))

	if err != nil {
		panic(err)
	}

	// sets up 
	if err := instance.SetupTest(ctx, t); err != nil {
		panic(err)
	}
	defer instance.CleanUp(ctx, t)

	connString, err := instance.GetConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	 _ = connString
}