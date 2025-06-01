package database

import (
	"Mehul-Kumar-27/dbporter/internal/models"
	"context"
)

type DatabaseAdapter interface {
	// Connection management
	Connect(ctx context.Context) error
	Close() error
	Ping(ctx context.Context) error

	// Schema operations
	ListTables(ctx context.Context) ([]string, error)
	GetTableSchema(ctx context.Context, tableName string) (*models.TableSchema, error)
	GetForeignKeys(ctx context.Context, tableName string) ([]models.ForeignKeySchema, error)
	CreateTable(ctx context.Context, schema *models.TableSchema) error
	CreateForeignKey(ctx context.Context, tableName string, fk *models.ForeignKeySchema) error
	DropForeignKey(ctx context.Context, tableName, constraintName string) error
	DropTable(ctx context.Context, tableName string) error

	// Data operations
	FetchData(ctx context.Context, query string, limit int) (*models.QueryResult, error)
	InsertData(ctx context.Context, tableName string, data []models.DataRow) error
	BulkInsert(ctx context.Context, tableName string, data []models.DataRow, batchSize int) error

	// Utility methods
	GetDatabaseType() string
	EscapeIdentifier(identifier string) string
	BuildInsertQuery(tableName string, columns []string) string
}

type Config struct {
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Database string            `json:"database"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	SSLMode  string            `json:"ssl_mode,omitempty"`
	Schema   string            `json:"schema,omitempty"`
	Extra    map[string]string `json:"extra,omitempty"`
}
