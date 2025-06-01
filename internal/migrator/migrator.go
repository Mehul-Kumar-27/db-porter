package migrator

import (
	"Mehul-Kumar-27/dbporter/internal/database"
	"context"
	"fmt"
	"log"
)

type Migrator struct {
	source      database.DatabaseAdapter
	destination database.DatabaseAdapter
	logger      *log.Logger
}

func NewMigrator(source, destination database.DatabaseAdapter, logger *log.Logger) *Migrator {
	return &Migrator{
		source:      source,
		destination: destination,
		logger:      logger,
	}
}

func (m *Migrator) MigrateTable(ctx context.Context, tableName string, recreateTable bool) error {
	m.logger.Printf("Starting migration for table: %s", tableName)

	// Get source table schema
	schema, err := m.source.GetTableSchema(ctx, tableName)
	if err != nil {
		return fmt.Errorf("failed to get source table schema: %w", err)
	}

	// Create table in destination if needed
	if recreateTable {
		if err := m.destination.DropTable(ctx, tableName); err != nil {
			m.logger.Printf("Warning: failed to drop table %s: %v", tableName, err)
		}

		if err := m.destination.CreateTable(ctx, schema); err != nil {
			return fmt.Errorf("failed to create destination table: %w", err)
		}
		m.logger.Printf("Created table %s in destination", tableName)
	}

	// Fetch data from source
	query := fmt.Sprintf("SELECT * FROM %s", m.source.EscapeIdentifier(tableName))
	result, err := m.source.FetchData(ctx, query, 0)
	if err != nil {
		return fmt.Errorf("failed to fetch data from source: %w", err)
	}

	if len(result.Rows) == 0 {
		m.logger.Printf("No data to migrate for table %s", tableName)
		return nil
	}

	// Insert data into destination
	if err := m.destination.BulkInsert(ctx, tableName, result.Rows, 1000); err != nil {
		return fmt.Errorf("failed to insert data into destination: %w", err)
	}

	m.logger.Printf("Migrated %d rows for table %s", len(result.Rows), tableName)
	return nil
}
