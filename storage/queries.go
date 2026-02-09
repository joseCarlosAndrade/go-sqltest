package storage

const (
	// company connection data
	GetCompaniesDBConnection = `SELECT id, db_name, db_write_connection, db_read_connection FROM companies`

	InsertCompaniesDBConnection = `INSERT INTO companies (id, db_name, db_write_connection, db_read_connection) VALUES (?, ?, ?, ?)`

	// global config data
	GetMySQLConfigs = `SELECT config_key, config_value FROM configs`
)