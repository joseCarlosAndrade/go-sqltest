package storage

import (
	"context"

	"github.com/joseCarlosAndrade/go-sqltest/pkg/populate"
)


type Storage interface {
	// Connect
	Connect(ctx context.Context, connString string) error
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
	InsertRows(ctx context.Context, data *populate.PopulateTable) error
}