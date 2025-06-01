package models

type TableSchema struct {
    Name         string               `json:"name"`
    Columns      []ColumnSchema       `json:"columns"`
    Indexes      []IndexSchema        `json:"indexes,omitempty"`
    ForeignKeys  []ForeignKeySchema   `json:"foreign_keys,omitempty"`
}

type ColumnSchema struct {
    Name         string            `json:"name"`
    Type         string            `json:"type"`
    Nullable     bool              `json:"nullable"`
    PrimaryKey   bool              `json:"primary_key"`
    DefaultValue string            `json:"default_value,omitempty"`
    ForeignKey   *ForeignKeySchema `json:"foreign_key,omitempty"`
}

type ForeignKeySchema struct {
    Name              string `json:"name"`
    Column            string `json:"column"`
    ReferencedTable   string `json:"referenced_table"`
    ReferencedColumn  string `json:"referenced_column"`
    OnDelete          string `json:"on_delete,omitempty"`
    OnUpdate          string `json:"on_update,omitempty"` 
}

type IndexSchema struct {
    Name    string   `json:"name"`
    Columns []string `json:"columns"`
    Unique  bool     `json:"unique"`
}

type DataRow map[string]interface{}

type QueryResult struct {
    Columns []string  `json:"columns"`
    Rows    []DataRow `json:"rows"`
    Count   int64     `json:"count"`
}
