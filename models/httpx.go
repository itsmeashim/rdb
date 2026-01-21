package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringArray is a custom type for PostgreSQL text arrays
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}
	return json.Unmarshal(bytes, a)
}

// HTTPXData represents the httpx JSON output (excluding unwanted fields)
type HTTPXData struct {
	ID            int64       `json:"id,omitempty" db:"id"`
	Port          string      `json:"port" db:"port"`
	URL           string      `json:"url" db:"url"`
	Input         string      `json:"input" db:"input"`
	Location      string      `json:"location" db:"location"`
	Title         string      `json:"title" db:"title"`
	Scheme        string      `json:"scheme" db:"scheme"`
	Webserver     string      `json:"webserver" db:"webserver"`
	ContentType   string      `json:"content_type" db:"content_type"`
	Method        string      `json:"method" db:"method"`
	Host          string      `json:"host" db:"host"`
	Path          string      `json:"path" db:"path"`
	Time          string      `json:"time" db:"time"`
	A             StringArray `json:"a" db:"a"`
	Tech          StringArray `json:"tech" db:"tech"`
	Words         int         `json:"words" db:"words"`
	Lines         int         `json:"lines" db:"lines"`
	StatusCode    int         `json:"status_code" db:"status_code"`
	ContentLength int         `json:"content_length" db:"content_length"`
	Program       string      `json:"program" db:"program"`
	Platform      string      `json:"platform" db:"platform"`
}
