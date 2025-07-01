package gsheet

type GsheetConfig struct {
	SheetID   string                      `yaml:"sheet_id"`
	Worksheets map[string]WorksheetConfig `yaml:"worksheets"`
}

type WorksheetConfig struct {
	Columns []ColumnConfig `yaml:"columns"`
}

type ColumnConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}
