package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rdb",
	Short: "Recon Database - Store and query httpx data",
	Long: `rdb (Recon Database) is a CLI tool to store and query httpx JSON output.

It uses PostgreSQL as the backend and supports filtering and sorting of stored data.

Quick start:
  1. Configure database connection:
     rdb config --connection-string "postgres://user:pass@localhost:5432/rdb"

  2. Store httpx output:
     httpx -l targets.txt -json | rdb store -p myprogram -P hackerone

  3. List stored data:
     rdb list --webserver nginx --sort url`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
