package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func APIInit(db *sql.DB) {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://localhost:4173/",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // fastest, usually best
	}))

	// POST /logs
	app.Post("/api/logs", func(c *fiber.Ctx) error {
		var log LogEntry
		if err := c.BodyParser(&log); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		// Validate and set fallback values
		if log.Level == "" {
			log.Level = "info"
		}
		if log.Message == "" {
			return c.Status(400).JSON(fiber.Map{"error": "message is required"})
		}
		if log.Source == "" {
			log.Source = "unknown"
		}
		if log.Hostname == "" {
			log.Hostname = "unknown"
		}
		if log.Environment == "" {
			log.Environment = "development"
		}

		err := InsertLog(db, log.Level, log.Message, log.Source, log.Hostname, log.Environment, log.Metadata)
		if err != nil {
			fmt.Println(err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save log: " + err.Error()})
		}

		return c.Status(201).JSON(fiber.Map{"status": "log saved"})
	})

	app.Get("/api/logs", func(c *fiber.Ctx) error {
		logs, err := GetLogs(db)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(logs)
	})

	app.Listen(":3000")
	log.Println("Backend API endpoints listening for :3000")
}
