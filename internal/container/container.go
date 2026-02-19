package container

import (
	"context"
	"testing"

	// "github.com/testcontainers/testcontainers-go"
	// "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type AccessConfig struct {
	DBName string
	Username string
	Password string
}

type Wrapper interface {
	Config(ctx context.Context, c *AccessConfig, schemaPath string) error 
	Start(ctx context.Context, t*testing.T) error
	Cleanup(ctx context.Context, t *testing.T) error
	Ping(ctx context.Context) error
	GetConnString(ctx context.Context) (string, error)
}