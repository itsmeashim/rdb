package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/itsmeashim/rdb/config"
	"github.com/itsmeashim/rdb/db"
	"github.com/spf13/cobra"
)

var (
	filterQuery       string
	filterURL         string
	filterInput       string
	filterTitle       string
	filterA           string
	filterWebserver   string
	filterTech        string
	filterHost        string
	filterScheme      string
	filterPort        string
	filterMethod      string
	filterPath        string
	filterLocation    string
	filterContentType string
	filterStatusCode  int
	filterProgram     string
	filterPlatform    string
	sortBy            string
	sortOrder         string
	limit             int
	outputJSON        bool
	separator         string
	listURLs          bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored httpx data",
	Long:  `List httpx data from the database with optional filters and sorting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := db.Init(cfg); err != nil {
			return err
		}
		defer db.Close()

		opts := db.ListOptions{
			Query:       filterQuery,
			URL:         filterURL,
			Input:       filterInput,
			Title:       filterTitle,
			A:           filterA,
			Webserver:   filterWebserver,
			Tech:        filterTech,
			Host:        filterHost,
			Scheme:      filterScheme,
			Port:        filterPort,
			Method:      filterMethod,
			Path:        filterPath,
			Location:    filterLocation,
			ContentType: filterContentType,
			StatusCode:  filterStatusCode,
			Program:     filterProgram,
			Platform:    filterPlatform,
			SortBy:      sortBy,
			SortOrder:   sortOrder,
			Limit:       limit,
		}

		results, err := db.List(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("failed to query data: %w", err)
		}

		if listURLs {
			for _, r := range results {
				fmt.Fprintln(os.Stdout, r.URL)
			}
			return nil
		}

		if outputJSON {
			encoder := json.NewEncoder(os.Stdout)
			for _, r := range results {
				encoder.Encode(r)
			}
			return nil
		}

		if len(results) == 0 {
			fmt.Println("no records found")
			return nil
		}

		if separator != "" {
			for _, r := range results {
				tech := strings.Join(r.Tech, ",")
				fmt.Printf("%s%s%d%s%s%s%s%s%s%s%s%s%s\n",
					r.URL, separator, r.StatusCode, separator, r.Webserver, separator,
					tech, separator, r.Title, separator, r.Program, separator, r.Platform)
			}
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, r := range results {
			title := r.Title
			if len(title) > 30 {
				title = title[:27] + "..."
			}
			tech := strings.Join(r.Tech, ",")
			if len(tech) > 30 {
				tech = tech[:27] + "..."
			}
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n",
				r.URL, r.StatusCode, r.Webserver, tech, title, r.Program, r.Platform)
		}
		w.Flush()
		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&filterQuery, "query", "q", "", "Search across common fields (url, input, title, host, webserver, content-type, tech, a, program, platform)")
	listCmd.Flags().StringVar(&filterURL, "url", "", "Filter by URL (partial match)")
	listCmd.Flags().StringVar(&filterInput, "input", "", "Filter by input (partial match)")
	listCmd.Flags().StringVar(&filterTitle, "title", "", "Filter by title (partial match)")
	listCmd.Flags().StringVar(&filterA, "a", "", "Filter by DNS A record (partial match)")
	listCmd.Flags().StringVar(&filterWebserver, "webserver", "", "Filter by webserver (partial match)")
	listCmd.Flags().StringVar(&filterTech, "tech", "", "Filter by technology (partial match)")
	listCmd.Flags().StringVar(&filterHost, "host", "", "Filter by host (partial match)")
	listCmd.Flags().StringVar(&filterScheme, "scheme", "", "Filter by scheme (http/https)")
	listCmd.Flags().StringVar(&filterPort, "port", "", "Filter by port (exact)")
	listCmd.Flags().StringVar(&filterMethod, "method", "", "Filter by method (exact)")
	listCmd.Flags().StringVar(&filterPath, "path", "", "Filter by path (partial match)")
	listCmd.Flags().StringVar(&filterLocation, "location", "", "Filter by redirect location (partial match)")
	listCmd.Flags().StringVar(&filterContentType, "content-type", "", "Filter by content-type (partial match)")
	listCmd.Flags().IntVar(&filterStatusCode, "status", 0, "Filter by HTTP status code (exact)")
	listCmd.Flags().StringVar(&filterProgram, "program", "", "Filter by program name")
	listCmd.Flags().StringVar(&filterPlatform, "platform", "", "Filter by platform name")
	listCmd.Flags().StringVar(&sortBy, "sort", "created_at", "Sort by field (url, input, title, host, scheme, port, method, path, location, content_type, status_code, content_length, words, lines, webserver, tech, program, platform, created_at)")
	listCmd.Flags().StringVar(&sortOrder, "order", "desc", "Sort order (asc, desc)")
	listCmd.Flags().IntVarP(&limit, "limit", "n", 0, "Limit number of results (0 = all)")
	listCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "Output as JSON")
	listCmd.Flags().StringVarP(&separator, "sep", "s", "", "Field separator for piping (e.g., ',' or '|')")
	listCmd.Flags().BoolVar(&listURLs, "urls", false, "Only output URLs")
	rootCmd.AddCommand(listCmd)
}
