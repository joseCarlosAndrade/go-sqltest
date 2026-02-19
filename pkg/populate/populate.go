package populate

type PopulateTable struct {
	TableName    string
	Columns []string
	Data    [][]any
}

func New(tableName string, columns ...string) *PopulateTable {
	// columns := make([]string, len(columns))
	return &PopulateTable{
		TableName:    tableName,
		Columns: columns,
		Data:    make([][]any, 0),
	}
}

func (p *PopulateTable) Insert(values ...any) *PopulateTable { 
	p.Data = append(p.Data, values)
	return p
}