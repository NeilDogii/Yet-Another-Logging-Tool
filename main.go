package main

import (
	"fmt"
	"log"
)

func main() {
	db, err := InitDB()

	if err != nil {
		log.Fatal("Failed to Initialize database!", err)
	} else {
		log.Println("Database connected successfully.")
	}

	defer db.Close()
	fmt.Printf("InsertLog: ")
	var log_str string
	fmt.Scan(&log_str)
	if log_str != "" {
		err = InsertLog(db, "INFO", log_str, "localhost", "localhost", "development", map[string]interface{}{"example_key": "example_value"})
		if err != nil {
			log.Fatal("Failed to insert log:", err)
		}
	}

	logs, err := GetLogs(db)
	if err != nil {
		log.Fatal("Could not fetch the logs from the db")
	}

	for _, logEntry := range logs {
		fmt.Println("Log Entry:", logEntry)
	}

	fmt.Println("Logger Initiated")
}
