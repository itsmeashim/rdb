package db

import (
	"context"
	"fmt"

	"github.com/itsmeashim/rdb/config"
	"github.com/itsmeashim/rdb/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

const createTableSQL = `
CREATE TABLE IF NOT EXISTS httpx_data (
    id SERIAL PRIMARY KEY,
    port TEXT,
    url TEXT,
    input TEXT,
    location TEXT,
    title TEXT,
    scheme TEXT,
    webserver TEXT,
    content_type TEXT,
    method TEXT,
    host TEXT,
    path TEXT,
    time TEXT,
    a JSONB,
    tech JSONB,
    words INT,
    lines INT,
    status_code INT,
    content_length INT,
    program TEXT DEFAULT 'default',
    platform TEXT DEFAULT 'default',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_url ON httpx_data(url);
CREATE INDEX IF NOT EXISTS idx_input ON httpx_data(input);
CREATE INDEX IF NOT EXISTS idx_webserver ON httpx_data(webserver);
CREATE INDEX IF NOT EXISTS idx_program ON httpx_data(program);
CREATE INDEX IF NOT EXISTS idx_platform ON httpx_data(platform);
`

func Init(cfg *config.Config) error {
	if cfg.ConnectionString == "" {
		return fmt.Errorf("connection string not configured. Run: rdb config --connection-string <postgres_url>")
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.ConnectionString)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxConnections)

	pool, err = pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	_, err = pool.Exec(context.Background(), createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

func Insert(ctx context.Context, data *models.HTTPXData) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO httpx_data (
			port, url, input, location, title, scheme, webserver,
			content_type, method, host, path, time, a, tech,
			words, lines, status_code, content_length, program, platform
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`, data.Port, data.URL, data.Input, data.Location, data.Title, data.Scheme, data.Webserver,
		data.ContentType, data.Method, data.Host, data.Path, data.Time, data.A, data.Tech,
		data.Words, data.Lines, data.StatusCode, data.ContentLength, data.Program, data.Platform)
	return err
}

type ListOptions struct {
	Query       string
	URL         string
	Input       string
	Title       string
	A           string
	Webserver   string
	Tech        string
	Host        string
	Scheme      string
	Port        string
	Method      string
	Path        string
	Location    string
	ContentType string
	StatusCode  int
	Program     string
	Platform    string
	SortBy      string
	SortOrder   string
	Limit       int
}

func List(ctx context.Context, opts ListOptions) ([]models.HTTPXData, error) {
	query := `SELECT id, port, url, input, location, title, scheme, webserver,
		content_type, method, host, path, time, a, tech, words, lines,
		status_code, content_length, program, platform
		FROM httpx_data WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if opts.Query != "" {
		// Single "search" term across common fields.
		ph := fmt.Sprintf("$%d", argNum)
		query += fmt.Sprintf(
			" AND (url ILIKE %[1]s OR input ILIKE %[1]s OR title ILIKE %[1]s OR host ILIKE %[1]s OR webserver ILIKE %[1]s OR content_type ILIKE %[1]s OR tech::text ILIKE %[1]s OR a::text ILIKE %[1]s OR program ILIKE %[1]s OR platform ILIKE %[1]s)",
			ph,
		)
		args = append(args, "%"+opts.Query+"%")
		argNum++
	}

	if opts.URL != "" {
		query += fmt.Sprintf(" AND url ILIKE $%d", argNum)
		args = append(args, "%"+opts.URL+"%")
		argNum++
	}
	if opts.Input != "" {
		query += fmt.Sprintf(" AND input ILIKE $%d", argNum)
		args = append(args, "%"+opts.Input+"%")
		argNum++
	}
	if opts.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argNum)
		args = append(args, "%"+opts.Title+"%")
		argNum++
	}
	if opts.A != "" {
		// Stored as JSONB array; text-cast works well for partial matches.
		query += fmt.Sprintf(" AND a::text ILIKE $%d", argNum)
		args = append(args, "%"+opts.A+"%")
		argNum++
	}
	if opts.Webserver != "" {
		query += fmt.Sprintf(" AND webserver ILIKE $%d", argNum)
		args = append(args, "%"+opts.Webserver+"%")
		argNum++
	}
	if opts.Tech != "" {
		query += fmt.Sprintf(" AND tech::text ILIKE $%d", argNum)
		args = append(args, "%"+opts.Tech+"%")
		argNum++
	}
	if opts.Host != "" {
		query += fmt.Sprintf(" AND host ILIKE $%d", argNum)
		args = append(args, "%"+opts.Host+"%")
		argNum++
	}
	if opts.Scheme != "" {
		query += fmt.Sprintf(" AND scheme = $%d", argNum)
		args = append(args, opts.Scheme)
		argNum++
	}
	if opts.Port != "" {
		query += fmt.Sprintf(" AND port = $%d", argNum)
		args = append(args, opts.Port)
		argNum++
	}
	if opts.Method != "" {
		query += fmt.Sprintf(" AND method = $%d", argNum)
		args = append(args, opts.Method)
		argNum++
	}
	if opts.Path != "" {
		query += fmt.Sprintf(" AND path ILIKE $%d", argNum)
		args = append(args, "%"+opts.Path+"%")
		argNum++
	}
	if opts.Location != "" {
		query += fmt.Sprintf(" AND location ILIKE $%d", argNum)
		args = append(args, "%"+opts.Location+"%")
		argNum++
	}
	if opts.ContentType != "" {
		query += fmt.Sprintf(" AND content_type ILIKE $%d", argNum)
		args = append(args, "%"+opts.ContentType+"%")
		argNum++
	}
	if opts.StatusCode != 0 {
		query += fmt.Sprintf(" AND status_code = $%d", argNum)
		args = append(args, opts.StatusCode)
		argNum++
	}
	if opts.Program != "" {
		query += fmt.Sprintf(" AND program = $%d", argNum)
		args = append(args, opts.Program)
		argNum++
	}
	if opts.Platform != "" {
		query += fmt.Sprintf(" AND platform = $%d", argNum)
		args = append(args, opts.Platform)
		argNum++
	}

	validSortColumns := map[string]bool{
		"port":           true,
		"url":            true,
		"input":          true,
		"title":          true,
		"scheme":         true,
		"webserver":      true,
		"content_type":   true,
		"method":         true,
		"host":           true,
		"path":           true,
		"location":       true,
		"a":              true,
		"tech":           true,
		"words":          true,
		"lines":          true,
		"status_code":    true,
		"content_length": true,
		"program":        true,
		"platform":       true,
		"created_at":     true,
	}
	sortBy := "created_at"
	if opts.SortBy != "" && validSortColumns[opts.SortBy] {
		sortBy = opts.SortBy
	}

	sortOrder := "DESC"
	if opts.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.HTTPXData, error) {
		var d models.HTTPXData
		err := row.Scan(&d.ID, &d.Port, &d.URL, &d.Input, &d.Location, &d.Title, &d.Scheme,
			&d.Webserver, &d.ContentType, &d.Method, &d.Host, &d.Path, &d.Time,
			&d.A, &d.Tech, &d.Words, &d.Lines, &d.StatusCode, &d.ContentLength,
			&d.Program, &d.Platform)
		return d, err
	})

	return results, err
}
