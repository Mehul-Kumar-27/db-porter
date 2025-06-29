package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbporter",
	Short: "Db-Porter helps you to move data around various data source",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
