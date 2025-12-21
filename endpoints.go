package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func APIInit(db *sql.DB) {
	app := fiber.New()

	// POST /logs
	app.Post("/api/logs", func(c *fiber.Ctx) error {
		var log LogEntry
		if err := c.BodyParser(&log); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		err := InsertLog(db, log.Level, log.Message, log.Source, log.Hostname, log.Environment, log.Metadata)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save log"})
		}

		return c.Status(201).JSON(fiber.Map{"status": "log saved"})
	})

	app.Get("/api/logs", func(c *fiber.Ctx) error {
		log.Printf(c.Query("test"))
		logs, err := GetLogs(db)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(logs)
	})

	app.Listen(":3000")
	log.Println("Backend API endpoints listening for :3000")
}
