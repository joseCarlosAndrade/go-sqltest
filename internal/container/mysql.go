package container

import (
	// "github.com/testcontainers/testcontainers-go"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// MySQLContainer implements Wrapper
type MySQLContainer struct {
	container       *mysql.MySQLContainer
	connString      string
	schemaDirectory string
	config * AccessConfig
}

var _ Wrapper = (*MySQLContainer)(nil)

func NewMySQLContainer() *MySQLContainer {
	return &MySQLContainer{}
}

func (m *MySQLContainer) Config(ctx context.Context, c *AccessConfig, schemaPath string) error {
	m.config = c
	m.schemaDirectory = schemaPath

	return nil
} 

func (m *MySQLContainer) Start(ctx context.Context, t*testing.T) error {
	// var mysqlContainer *mysql.MySQLContainer
	var err error

	params := make([]testcontainers.ContainerCustomizer, 0)

	if m.schemaDirectory != "" {
		params = append(params, mysql.WithScripts(m.schemaDirectory))
	}

	// todo: properly allow passwordless instances
	params = append(params,
		 mysql.WithDatabase(m.config.DBName),
		 mysql.WithUsername(m.config.Username),
		 mysql.WithPassword(m.config.Password))


	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		params...)
	
	if err != nil {
		t.Logf("failed to start container: %s", err.Error())
		return err
	}

	time.Sleep(2 * time.Second) // sleep a bit to give proper container init

	t.Logf("container successfully initialized")

	m.container = mysqlContainer
	m.connString, err = mysqlContainer.ConnectionString(ctx)
	if err != nil {
		t.Logf("could not get mysql connection string: %s", err.Error())
	}

	return nil
}

func (m *MySQLContainer) Cleanup(ctx context.Context, t *testing.T) error {

	if err := testcontainers.TerminateContainer(m.container); err != nil {
		t.Logf("failed to terminate container: %s", err.Error())

		return err
	}
	
	t.Logf("container successfully terminated")

	return nil
}

func (m *MySQLContainer) Ping(ctx context.Context) error {
	if m.container != nil {
		if m.container.IsRunning() {
			return nil
		} 

		return fmt.Errorf("container is stopped")
	}

	return fmt.Errorf("contaier not initialized")
}

func (m *MySQLContainer) GetConnString(ctx context.Context) (string, error) {
	if m.container != nil && m.connString != "" {
		return m.connString, nil
	}

	return "", fmt.Errorf("contaier not initialized")
}