package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type LogEntry struct {
	ID          int64                  `db:"id"`
	Timestamp   time.Time              `db:"timestamp"`
	Level       string                 `db:"level"`
	Message     string                 `db:"message"`
	Source      string                 `db:"source"`
	Hostname    string                 `db:"hostname"`
	Environment string                 `db:"environment"`
	Metadata    map[string]interface{} `db:"metadata"`
}

const FETCH_LIMIT = 100

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "logs.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Enable WAL mode for better concurrency
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        level TEXT NOT NULL DEFAULT 'INFO' CHECK(level IN ('TRACE', 'DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL')),
        message TEXT NOT NULL,
        source TEXT DEFAULT 'unknown',
        hostname TEXT DEFAULT 'localhost',
        environment TEXT DEFAULT 'development',
        metadata TEXT DEFAULT '{}'
    );
		CREATE INDEX IF NOT EXISTS idx_timestamp ON logs(timestamp);
		CREATE INDEX IF NOT EXISTS idx_level ON logs(level);
		CREATE INDEX IF NOT EXISTS idx_source ON logs(source);
		CREATE INDEX IF NOT EXISTS idx_hostname ON logs(hostname);
		CREATE INDEX IF NOT EXISTS idx_environment ON logs(environment);
		CREATE INDEX IF NOT EXISTS idx_level_timestamp ON logs(level, timestamp);
`)

	if err != nil {
		return nil, err
	}

	log.Println("Database Initialized")
	return db, nil
}

func InsertLog(db *sql.DB, level, message, source, hostname, environment string, metadata map[string]interface{}) error {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
        INSERT INTO logs (level, message, source, hostname, environment, metadata)
        VALUES (?, ?, ?, ?, ?, ?)
    `, level, message, source, hostname, environment, string(metadataJSON))

	return err
}

func GetLogs(db *sql.DB, page ...int) ([]LogEntry, error) {
	pageNum := 0
	if len(page) > 0 {
		pageNum = page[0]
	}

	offset := pageNum * FETCH_LIMIT

	rows, err := db.Query(`SELECT id, timestamp, level, message, source, hostname, environment, metadata FROM logs LIMIT ? OFFSET ?`, FETCH_LIMIT, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var entry LogEntry
		var metadataStr string

		err := rows.Scan(
			&entry.ID,
			&entry.Timestamp,
			&entry.Level,
			&entry.Message,
			&entry.Source,
			&entry.Hostname,
			&entry.Environment,
			&metadataStr,
		)
		if err != nil {
			return nil, err
		}

		if metadataStr != "" {
			err = json.Unmarshal([]byte(metadataStr), &entry.Metadata)
			if err != nil {
				entry.Metadata = make(map[string]interface{})
			}
		} else {
			entry.Metadata = make(map[string]interface{})
		}

		logs = append(logs, entry)
	}

	return logs, rows.Err()
}

type PaginationInfo struct {
	TotalLogs   int64 `json:"total_logs"`
	TotalPages  int   `json:"total_pages"`
	LogsPerPage int   `json:"logs_per_page"`
}

func GetPaginationInfo(db *sql.DB) (*PaginationInfo, error) {
	logCount, err := db.Query(`SELECT count(*) FROM logs`)
	if err != nil {
		return nil, err
	}
	defer logCount.Close()

	var totalLogs int64
	if logCount.Next() {
		err = logCount.Scan(&totalLogs)
		if err != nil {
			return nil, err
		}
	}

	totalPages := int(totalLogs) / FETCH_LIMIT
	if int(totalLogs)%FETCH_LIMIT != 0 {
		totalPages++
	}

	return &PaginationInfo{
		TotalLogs:   totalLogs,
		TotalPages:  totalPages,
		LogsPerPage: FETCH_LIMIT,
	}, nil
}
