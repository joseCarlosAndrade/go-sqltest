package storage

type SeedConfig struct {
	TableName string
	Columns   []string
	Data      [][]any
}

// buildInsertQuery builds a parameterized INSERT statement using ? placeholders.
// Returns the query string and a flattened slice of values for Exec.
func buildInsertQuery(data *SeedConfig) (query string, args []any) {
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