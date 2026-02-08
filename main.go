package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "5555"
	}
	return port
}

func handlePost(c *fiber.Ctx) error {
	c.Accepts("application/json")

	body := c.Body()

	// Try to parse as array first (most common), then as single object
	var entries []LogEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		var single LogEntry
		if err := json.Unmarshal(body, &single); err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
		}
		entries = []LogEntry{single}
	}

	for _, entry := range entries {
		if err := insertLogEntry(entry); err != nil {
			log.Printf("Error inserting log entry: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save log")
		}
	}

	log.Printf("Inserted %d log entries\n", len(entries))
	return c.Status(fiber.StatusCreated).SendString("Log received")
}

func main() {
	app := fiber.New()
	dbInit()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("NetWatcher Logs")
	})
	app.Post("/", handlePost)

	port := getPort()
	log.Fatal(app.Listen(":" + port))
}
