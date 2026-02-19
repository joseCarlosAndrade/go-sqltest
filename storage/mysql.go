package storage

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)


// Storage holds a MySQL client and exposes Ping and custom SQL operations.
type Storage struct {
	client *sql.DB
}


// Connect opens a connection to MySQL and returns a Wrapper holding the client.
func NewStorage(ctx context.Context, connString string) (*Storage, error) {
	client, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	if err := client.PingContext(ctx); err != nil {
		_ = client.Close()
		return nil, err
	}
	return &Storage{client: client}, nil
}



// Client returns the underlying *sql.DB for direct use.
func (s *Storage) Client() *sql.DB {
	return s.client
}

// Ping verifies the connection to the database is still alive.
func (s *Storage) Ping(ctx context.Context) error {
	return s.client.PingContext(ctx)
}

// Close closes the underlying database connection.
func (s *Storage) Close() error {
	return s.client.Close()
}

func (s *Storage) InsertNexusCompany(ctx context.Context, id, connString, dbname string) error {
	stmt, err := s.client.PrepareContext(ctx, InsertCompaniesDBConnection)
	if err != nil {
		return err
	}
	
	defer stmt.Close()

	_, err = stmt.Exec(id, dbname, connString, connString)
	if err != nil {
		return err
	}
	
	return nil
}



func (s *Storage) InsertNexusConfig(ctx context.Context, data *SeedConfig) error {
	query, args := buildInsertQuery(data)
	if query == "" {
		return nil
	}
	_, err := s.client.ExecContext(ctx, query, args...)
	return err
}

func (s *Storage) InsertPricingData(ctx context.Context, data *SeedConfig) error {
	query, args := buildInsertQuery(data)
	if query == "" {
		return nil
	}
	_, err := s.client.ExecContext(ctx, query, args...)
	return err
}