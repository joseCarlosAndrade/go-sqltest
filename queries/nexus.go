package queries

const (
	// company connection data
	GetCompaniesDBConnection = `SELECT id, db_name, db_write_connection, db_read_connection FROM companies`

	// global config data
	GetMySQLConfigs = `SELECT config_key, config_value FROM configs`
)