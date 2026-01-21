# rdb - Recon Database

A fast CLI tool for storing and querying [httpx](https://github.com/projectdiscovery/httpx) reconnaissance data in PostgreSQL.

## Installation

### Using Go

```bash
go install github.com/itsmeashim/rdb@latest
```

### From Releases

Download the latest binary from [Releases](https://github.com/itsmeashim/rdb/releases).

### Build from source

```bash
git clone https://github.com/itsmeashim/rdb.git
cd rdb
go build -o rdb .
```

## Prerequisites

- PostgreSQL database
- [httpx](https://github.com/projectdiscovery/httpx) for generating reconnaissance data

## Quick Start

```bash
# 1. Configure database connection
rdb config --connection-string "postgres://user:password@localhost:5432/rdb"

# 2. Store httpx data
httpx -l targets.txt -json | rdb store -p myprogram

# 3. Query stored data
rdb list --webserver nginx --limit 50
```

## Commands

### `rdb config`

Configure database connection and defaults.

```bash
# Set connection string
rdb config --connection-string "postgres://user:pass@localhost:5432/rdb"

# Set defaults
rdb config --default-program bugbounty --default-platform hackerone

# View current config
rdb config
```

| Flag | Description |
|------|-------------|
| `--connection-string` | PostgreSQL connection URL |
| `--max-connections` | Connection pool size (default: 10) |
| `--default-program` | Default program name |
| `--default-platform` | Default platform name |

### `rdb store`

Store httpx JSON output from stdin.

```bash
# Basic usage
httpx -l targets.txt -json | rdb store

# With program and platform tags
httpx -l targets.txt -json | rdb store -p myprogram --platform hackerone

# From file
cat httpx_output.json | rdb store -p myprogram
```

| Flag | Short | Description |
|------|-------|-------------|
| `--program` | `-p` | Program identifier |
| `--platform` | | Platform identifier |

### `rdb list`

Query stored data with filters.

```bash
# Filter by webserver
rdb list --webserver nginx

# Filter by technology
rdb list --tech Cloudflare

# Filter by program
rdb list --program myprogram --platform hackerone

# Combine filters with sorting
rdb list --webserver nginx --tech PHP --sort url --order asc --limit 100

# JSON output
rdb list --program myprogram --json

# Custom separator for piping
rdb list --sep "," | cut -d',' -f1
```

#### Filter Options

| Flag | Match Type | Description |
|------|------------|-------------|
| `--url` | partial | Filter by URL |
| `--input` | partial | Filter by input domain |
| `--webserver` | partial | Filter by web server |
| `--tech` | partial | Filter by technology |
| `--program` | exact | Filter by program name |
| `--platform` | exact | Filter by platform name |

#### Sort & Output Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--sort` | | `created_at` | Sort field |
| `--order` | | `desc` | Sort direction (asc/desc) |
| `--limit` | `-n` | all | Limit results |
| `--json` | `-j` | false | JSON output |
| `--sep` | `-s` | | Custom separator |

Valid sort fields: `url`, `input`, `webserver`, `tech`, `program`, `platform`, `created_at`

## Data Model

Each record stores the following httpx fields:

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Full URL |
| `input` | string | Original input domain |
| `status_code` | int | HTTP status code |
| `title` | string | Page title |
| `webserver` | string | Web server type |
| `tech` | []string | Detected technologies |
| `content_length` | int | Response size |
| `content_type` | string | Content-Type header |
| `host` | string | Hostname/IP |
| `port` | string | Port number |
| `scheme` | string | http/https |
| `path` | string | URL path |
| `method` | string | HTTP method |
| `location` | string | Redirect location |
| `a` | []string | DNS A records |
| `words` | int | Word count |
| `lines` | int | Line count |
| `time` | string | Response time |
| `program` | string | Custom program tag |
| `platform` | string | Custom platform tag |

## Examples

### Bug Bounty Workflow

```bash
# Configure once
rdb config --connection-string "postgres://user:pass@localhost:5432/recon"
rdb config --default-platform hackerone

# Daily recon
subfinder -d target.com | httpx -json | rdb store -p target-program

# Find interesting targets
rdb list --tech WordPress --program target-program
rdb list --webserver nginx --status 200 --limit 50
rdb list --input admin --json > admin_panels.json
```

### Multi-program Management

```bash
# Store data for different programs
cat program1_httpx.json | rdb store -p program1 --platform bugcrowd
cat program2_httpx.json | rdb store -p program2 --platform hackerone

# Query specific program
rdb list --program program1 --sort url

# Query all from platform
rdb list --platform hackerone --limit 1000 --json
```

## Database Schema

The tool automatically creates the required table and indexes:

```sql
CREATE TABLE httpx_data (
    id SERIAL PRIMARY KEY,
    url TEXT,
    input TEXT,
    -- ... other fields
    program TEXT DEFAULT 'default',
    platform TEXT DEFAULT 'default',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for fast queries
CREATE INDEX idx_url ON httpx_data(url);
CREATE INDEX idx_input ON httpx_data(input);
CREATE INDEX idx_webserver ON httpx_data(webserver);
CREATE INDEX idx_program ON httpx_data(program);
CREATE INDEX idx_platform ON httpx_data(platform);
```

## Configuration

Config file location: `~/.config/rdb/config.json`

```json
{
  "connection_string": "postgres://user:pass@localhost:5432/rdb",
  "max_connections": 10,
  "default_program": "default",
  "default_platform": "default"
}
```

## License

MIT
