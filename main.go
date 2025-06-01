package main

import (
	"Mehul-Kumar-27/dbporter/internal/database"
	"Mehul-Kumar-27/dbporter/internal/database/postgres"
	"Mehul-Kumar-27/dbporter/logger"
	"context"

	_ "github.com/lib/pq"
)

// postgresql://postgres:postgres@127.0.0.1:5435/cms?statusColor=686B6F&env=&name=CMS%20LOCAL%20POSTGRES&tLSMode=0&usePrivateKey=false&safeModeLevel=0&advancedSafeModeLevel=0&driverVersion=0&lazyload=false
func main() {
	log := logger.New(nil)
	log.Info("Starting DB Porter")
	sourceConfig := database.Config{
		Host:     "127.0.0.1",
		Port:     5435,
		Database: "cms",
		Username: "postgres",
		Password: "postgres",
		SSLMode:  "disable",
	}

	source_postgres := postgres.NewPostgresAdapter(&sourceConfig)
	ctx := context.Background()
	err := source_postgres.Connect(ctx)
	if err != nil {
		log.Error("error connecting to source database: %v\n", err)
	}
	log.Info("Connected to the source postgres adapter")
	destination_postgres := postgres.NewPostgresAdapter(&sourceConfig)
	err = destination_postgres.Connect(ctx)
	if err != nil {
		log.Error("error connecting to the destination postgres adapter")
	}
	log.Info("Connected to the destination postgres adapter")
}
