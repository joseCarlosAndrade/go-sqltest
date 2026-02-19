package mysql

import (
	"github.com/joseCarlosAndrade/go-sqltest/pkg/populate"
)

// buildInsertQuery builds a parameterized INSERT statement using ? placeholders.
// Returns the query string and a flattened slice of values for Exec.
// mysql dialect
func buildInsertQuery(data *populate.PopulateTable) (query string, args []any) {
	if data == nil || len(data.Columns) == 0 || len(data.Data) == 0 {
		return "", nil
	}

	table := data.TableName
	if table == "" {
		table = "configs"
	}

	query = "INSERT INTO " + table + " ("
	for i, col := range data.Columns {
		if i > 0 {
			query += ", "
		}
		query += col
	}
	query += ") VALUES "

	placeholders := "("
	for i := range data.Columns {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
	}
	placeholders += ")"

	for i := range data.Data {
		if i > 0 {
			query += ", "
		}
		query += placeholders
		args = append(args, data.Data[i]...)
	}

	return query, args
}