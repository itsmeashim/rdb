package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/itsmeashim/rdb/config"
	"github.com/itsmeashim/rdb/db"
	"github.com/itsmeashim/rdb/models"
)

var (
	program  string
	platform string
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store httpx JSON data from stdin",
	Long:  `Reads httpx JSON output from stdin (piped) and stores it in the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := db.Init(cfg); err != nil {
			return err
		}
		defer db.Close()

		if program == "" {
			program = cfg.DefaultProgram
		}
		if platform == "" {
			platform = cfg.DefaultPlatform
		}

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no input provided. Pipe httpx JSON output to this command")
		}

		scanner := bufio.NewScanner(os.Stdin)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)

		count := 0
		ctx := context.Background()

		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var data models.HTTPXData
			if err := json.Unmarshal(line, &data); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to parse JSON: %v\n", err)
				continue
			}

			data.Program = program
			data.Platform = platform

			if err := db.Insert(ctx, &data); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to insert: %v\n", err)
				continue
			}
			count++
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		fmt.Printf("stored %d records\n", count)
		return nil
	},
}

func init() {
	storeCmd.Flags().StringVarP(&program, "program", "p", "", "Program name (e.g., bugcrowd-program)")
	storeCmd.Flags().StringVar(&platform, "platform", "", "Platform name (e.g., hackerone, bugcrowd)")
	rootCmd.AddCommand(storeCmd)
}
