package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/itsmeashim/rdb/config"
	"github.com/itsmeashim/rdb/models"
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
	URL       string
	Input     string
	Webserver string
	Tech      string
	Program   string
	Platform  string
	SortBy    string
	SortOrder string
	Limit     int
}

func List(ctx context.Context, opts ListOptions) ([]models.HTTPXData, error) {
	query := `SELECT id, port, url, input, location, title, scheme, webserver,
		content_type, method, host, path, time, a, tech, words, lines,
		status_code, content_length, program, platform
		FROM httpx_data WHERE 1=1`
	args := []interface{}{}
	argNum := 1

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
		"url": true, "input": true, "webserver": true, "tech": true,
		"program": true, "platform": true, "created_at": true,
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
