package postgres

import (
	"Mehul-Kumar-27/dbporter/internal/database"
	"Mehul-Kumar-27/dbporter/internal/models"
	"Mehul-Kumar-27/dbporter/logger"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type PostgresAdapter struct {
	config *database.Config
	db     *sql.DB
}

func NewPostgresAdapter(config *database.Config) *PostgresAdapter {
	return &PostgresAdapter{
		config: config,
	}
}

func (p *PostgresAdapter) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", p.config.Host, p.config.Port, p.config.Database, p.config.Username, p.config.Password, p.config.SSLMode)
	log := logger.New(nil)
	log.Info("Connection string for the postgres connector is %s", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	p.db = db
	return nil
}

func (p *PostgresAdapter) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresAdapter) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *PostgresAdapter) ListTables(ctx context.Context) ([]string, error) {
	query := `
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
        ORDER BY table_name`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, rows.Err()
}

func (p *PostgresAdapter) GetTableSchema(ctx context.Context, tableName string) (*models.TableSchema, error) {
	// Get columns
	columnQuery := `
        SELECT column_name, data_type, is_nullable, column_default
        FROM information_schema.columns
        WHERE table_name = $1 AND table_schema = 'public'
        ORDER BY ordinal_position`

	rows, err := p.db.QueryContext(ctx, columnQuery, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []models.ColumnSchema
	for rows.Next() {
		var col models.ColumnSchema
		var nullable, defaultVal sql.NullString

		if err := rows.Scan(&col.Name, &col.Type, &nullable, &defaultVal); err != nil {
			return nil, err
		}

		col.Nullable = nullable.String == "YES"
		if defaultVal.Valid {
			col.DefaultValue = defaultVal.String
		}

		columns = append(columns, col)
	}

	// Get primary keys
	pkQuery := `
        SELECT column_name
        FROM information_schema.key_column_usage k
        JOIN information_schema.table_constraints t ON k.constraint_name = t.constraint_name
        WHERE t.table_name = $1 AND t.constraint_type = 'PRIMARY KEY'`

	pkRows, err := p.db.QueryContext(ctx, pkQuery, tableName)
	if err != nil {
		return nil, err
	}
	defer pkRows.Close()

	pkColumns := make(map[string]bool)
	for pkRows.Next() {
		var colName string
		if err := pkRows.Scan(&colName); err != nil {
			return nil, err
		}
		pkColumns[colName] = true
	}

	// Mark primary key columns
	for i := range columns {
		if pkColumns[columns[i].Name] {
			columns[i].PrimaryKey = true
		}
	}

	// Get foreign keys
	foreignKeys, err := p.GetForeignKeys(ctx, tableName)
	if err != nil {
		return nil, err
	}

	// Map foreign keys to columns
	fkMap := make(map[string]*models.ForeignKeySchema)
	for i := range foreignKeys {
		fkMap[foreignKeys[i].Column] = &foreignKeys[i]
	}

	// Assign foreign keys to columns
	for i := range columns {
		if fk, exists := fkMap[columns[i].Name]; exists {
			columns[i].ForeignKey = fk
		}
	}

	return &models.TableSchema{
		Name:        tableName,
		Columns:     columns,
		ForeignKeys: foreignKeys,
	}, nil
}

func (p *PostgresAdapter) GetForeignKeys(ctx context.Context, tableName string) ([]models.ForeignKeySchema, error) {
	query := `
        SELECT
            tc.constraint_name,
            kcu.column_name,
            ccu.table_name AS foreign_table_name,
            ccu.column_name AS foreign_column_name,
            rc.delete_rule,
            rc.update_rule
        FROM information_schema.table_constraints AS tc
        JOIN information_schema.key_column_usage AS kcu
            ON tc.constraint_name = kcu.constraint_name
            AND tc.table_schema = kcu.table_schema
        JOIN information_schema.constraint_column_usage AS ccu
            ON ccu.constraint_name = tc.constraint_name
            AND ccu.table_schema = tc.table_schema
        JOIN information_schema.referential_constraints AS rc
            ON tc.constraint_name = rc.constraint_name
            AND tc.table_schema = rc.constraint_schema
        WHERE tc.constraint_type = 'FOREIGN KEY'
            AND tc.table_name = $1
            AND tc.table_schema = 'public'`

	rows, err := p.db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foreignKeys []models.ForeignKeySchema
	for rows.Next() {
		var fk models.ForeignKeySchema
		var deleteRule, updateRule sql.NullString

		if err := rows.Scan(&fk.Name, &fk.Column, &fk.ReferencedTable,
			&fk.ReferencedColumn, &deleteRule, &updateRule); err != nil {
			return nil, err
		}

		if deleteRule.Valid {
			fk.OnDelete = deleteRule.String
		}
		if updateRule.Valid {
			fk.OnUpdate = updateRule.String
		}

		foreignKeys = append(foreignKeys, fk)
	}

	return foreignKeys, rows.Err()
}

func (p *PostgresAdapter) CreateTable(ctx context.Context, schema *models.TableSchema) error {
	var columnDefs []string
	var primaryKeys []string

	for _, col := range schema.Columns {
		def := fmt.Sprintf("%s %s", p.EscapeIdentifier(col.Name), col.Type)

		if !col.Nullable {
			def += " NOT NULL"
		}

		if col.DefaultValue != "" {
			def += " DEFAULT " + col.DefaultValue
		}

		if col.PrimaryKey {
			primaryKeys = append(primaryKeys, col.Name)
		}

		columnDefs = append(columnDefs, def)
	}

	if len(primaryKeys) > 0 {
		pkDef := fmt.Sprintf("PRIMARY KEY (%s)",
			strings.Join(primaryKeys, ", "))
		columnDefs = append(columnDefs, pkDef)
	}

	query := fmt.Sprintf("CREATE TABLE %s (%s)",
		p.EscapeIdentifier(schema.Name),
		strings.Join(columnDefs, ", "))

	_, err := p.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	// Create foreign key constraints separately
	for _, fk := range schema.ForeignKeys {
		if err := p.CreateForeignKey(ctx, schema.Name, &fk); err != nil {
			return fmt.Errorf("failed to create foreign key %s: %w", fk.Name, err)
		}
	}

	return nil
}

func (p *PostgresAdapter) CreateForeignKey(ctx context.Context, tableName string, fk *models.ForeignKeySchema) error {
	query := fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s
        FOREIGN KEY (%s) REFERENCES %s(%s)`,
		p.EscapeIdentifier(tableName),
		p.EscapeIdentifier(fk.Name),
		p.EscapeIdentifier(fk.Column),
		p.EscapeIdentifier(fk.ReferencedTable),
		p.EscapeIdentifier(fk.ReferencedColumn))

	if fk.OnDelete != "" {
		query += " ON DELETE " + fk.OnDelete
	}

	if fk.OnUpdate != "" {
		query += " ON UPDATE " + fk.OnUpdate
	}

	_, err := p.db.ExecContext(ctx, query)
	return err
}

func (p *PostgresAdapter) EscapeIdentifier(identifier string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(identifier, `"`, `""`))
}
func (p *PostgresAdapter) DropTable(ctx context.Context, tableName string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", p.EscapeIdentifier(tableName))
	_, err := p.db.ExecContext(ctx, query)
	return err
}

func (p *PostgresAdapter) FetchData(ctx context.Context, query string, limit int) (*models.QueryResult, error) {
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var data []models.DataRow
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(models.DataRow)
		for i, col := range columns {
			row[col] = values[i]
		}

		data = append(data, row)
	}

	return &models.QueryResult{
		Columns: columns,
		Rows:    data,
		Count:   int64(len(data)),
	}, rows.Err()
}

func (p *PostgresAdapter) InsertData(ctx context.Context, tableName string, data []models.DataRow) error {
	return p.BulkInsert(ctx, tableName, data, 1000)
}

func (p *PostgresAdapter) BulkInsert(ctx context.Context, tableName string, data []models.DataRow, batchSize int) error {
	if len(data) == 0 {
		return nil
	}

	// Get column names from first row
	var columns []string
	for col := range data[0] {
		columns = append(columns, col)
	}

	query := p.BuildInsertQuery(tableName, columns)

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		batch := data[i:end]
		if err := p.insertBatch(ctx, query, columns, batch); err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresAdapter) insertBatch(ctx context.Context, query string, columns []string, batch []models.DataRow) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range batch {
		values := make([]interface{}, len(columns))
		for i, col := range columns {
			values[i] = row[col]
		}

		if _, err := stmt.ExecContext(ctx, values...); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (p *PostgresAdapter) GetDatabaseType() string {
	return "postgres"
}

func (p *PostgresAdapter) BuildInsertQuery(tableName string, columns []string) string {
	escapedCols := make([]string, len(columns))
	placeholders := make([]string, len(columns))

	for i, col := range columns {
		escapedCols[i] = p.EscapeIdentifier(col)
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		p.EscapeIdentifier(tableName),
		strings.Join(escapedCols, ", "),
		strings.Join(placeholders, ", "))
}

func (p *PostgresAdapter) DropForeignKey(ctx context.Context, tableName, constraintName string) error {
	query := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s",
		p.EscapeIdentifier(tableName),
		p.EscapeIdentifier(constraintName))

	_, err := p.db.ExecContext(ctx, query)
	return err
}
