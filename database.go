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

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "logs.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
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
    )
		CREATE INDEX idx_timestamp ON logs(timestamp);
		CREATE INDEX idx_level ON logs(level);
		CREATE INDEX idx_source ON logs(source);
		CREATE INDEX idx_hostname ON logs(hostname);
		CREATE INDEX idx_environment ON logs(environment);

		CREATE INDEX idx_level_timestamp ON logs(level, timestamp);
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

func GetLogs(db *sql.DB) ([]LogEntry, error) {
	rows, err := db.Query(`SELECT id, timestamp, level, message, source, hostname, environment, metadata FROM logs`)
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
