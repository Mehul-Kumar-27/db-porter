package cmd

import (
	"Mehul-Kumar-27/dbporter/internal/core"
	"Mehul-Kumar-27/dbporter/logger"

	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer data between specified sources",
	Long: `Transfer data between specified data sources using configuration files or flags.

Examples:
  dbporter transfer --config config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.New("DBPORTER", nil)
		logger.Info("Reading configuration file at %s", cfgFile)
		pipelineConfig, err := core.LoadPipelineConfig(cfgFile)
		if err != nil {
			logger.Error("Error loading pipeline config: %s", err)
			return
		}
		logger.Info("Pipeline config loaded successfully")
		logger.Info("Source config: %v", pipelineConfig.SourceConfig.Type)
		logger.Info("Destination config: %v", pipelineConfig.DestinationConfig.Type)
	},
}
