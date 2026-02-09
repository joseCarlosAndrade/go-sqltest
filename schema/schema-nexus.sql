-- CREATE DATABASE IF NOT EXISTS test;

-- CREATE DATABASE IF NOT EXISTS nexus;
-- USE nexus;

CREATE TABLE companies (
	id char(36) NOT NULL,
	db_name varchar(255) NOT NULL,
	db_write_connection text NULL,
	db_read_connection text NULL,
	CONSTRAINT companies_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX companies_ix1 ON companies (db_name);

-- Create the configs table.
CREATE TABLE configs (
	config_key varchar(255) NOT NULL,
	config_value varchar(255) NOT NULL,
	CONSTRAINT configs_pkey PRIMARY KEY (config_key)
);

-- PALCEHOLDER DATA
INSERT INTO companies (id, db_name, db_write_connection, db_read_connection) 
VALUES
('xxxxxxx', 'placeholder', 'root:password@tcp(localhost:3306)/pricing', 'root:password@tcp(localhost:3306)/pricing');

INSERT INTO configs (config_key, config_value)
VALUES
('CURRENTMYSQLHOST', 'localhost:3360');  