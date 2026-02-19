package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joseCarlosAndrade/go-sqltest/internal/storage"
	"github.com/joseCarlosAndrade/go-sqltest/pkg/populate"
)

type MySQLStorage struct {
	connString string
	client     *sql.DB
	status     int // todo: iota for: not connected, connected, closed, etc 
}

var _ storage.Storage = (*MySQLStorage)(nil)

func NewStorage(ctx context.Context) *MySQLStorage {

	return &MySQLStorage{}
}

func (s *MySQLStorage) Connect(ctx context.Context, connString string) error {
	s.connString = connString

	client, err := sql.Open("mysql", connString)
	if err != nil {
		return err
	}
	if err := client.PingContext(ctx); err != nil {
		_ = client.Close()
		return err
	}

	s.client = client
	s.status = 1

	return nil
}

func (s *MySQLStorage)Ping(ctx context.Context) error {
	if s.status == 0 {
		return fmt.Errorf("mysql instance is not connected")
	}

	if s.client == nil {
		return fmt.Errorf("mysql instance is not initialized")
	}

	return s.client.PingContext(ctx)
}

// Close closes the underlying database connection.
func (s *MySQLStorage) Close(ctx context.Context) error {
	s.status = -1
	return s.client.Close()
}


func (s *MySQLStorage) InsertRows(ctx context.Context, data *populate.PopulateTable) error {
	query, args := buildInsertQuery(data)
	if query == "" {
		return nil
	}
	_, err := s.client.ExecContext(ctx, query, args...)
	return err
}
