package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/itsmeashim/rdb/config"
)

var (
	connString     string
	maxConnections int
	defaultProgram string
	defaultPlatform string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure rdb settings",
	Long:  `Set configuration options like database connection string and defaults.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		changed := false

		if connString != "" {
			cfg.ConnectionString = connString
			changed = true
		}
		if maxConnections > 0 {
			cfg.MaxConnections = maxConnections
			changed = true
		}
		if defaultProgram != "" {
			cfg.DefaultProgram = defaultProgram
			changed = true
		}
		if defaultPlatform != "" {
			cfg.DefaultPlatform = defaultPlatform
			changed = true
		}

		if changed {
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Println("configuration saved")
		}

		path, _ := config.ConfigPath()
		fmt.Printf("\nconfig file: %s\n", path)
		fmt.Printf("connection_string: %s\n", maskConnString(cfg.ConnectionString))
		fmt.Printf("max_connections: %d\n", cfg.MaxConnections)
		fmt.Printf("default_program: %s\n", cfg.DefaultProgram)
		fmt.Printf("default_platform: %s\n", cfg.DefaultPlatform)

		return nil
	},
}

func maskConnString(s string) string {
	if s == "" {
		return "(not set)"
	}
	if len(s) > 20 {
		return s[:10] + "..." + s[len(s)-10:]
	}
	return "***"
}

func init() {
	configCmd.Flags().StringVar(&connString, "connection-string", "", "PostgreSQL connection string")
	configCmd.Flags().IntVar(&maxConnections, "max-connections", 0, "Maximum database connections")
	configCmd.Flags().StringVar(&defaultProgram, "default-program", "", "Default program name")
	configCmd.Flags().StringVar(&defaultPlatform, "default-platform", "", "Default platform name")
	rootCmd.AddCommand(configCmd)
}
