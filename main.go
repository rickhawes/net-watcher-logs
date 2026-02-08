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
	var reports []Report
	if err := json.Unmarshal(body, &reports); err != nil {
		var single Report
		if err := json.Unmarshal(body, &single); err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
		}
		reports = []Report{single}
	}

	for _, report := range reports {
		if err := insertReport(report); err != nil {
			log.Printf("Error inserting report: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save report")
		}
	}

	log.Printf("Inserted %d reports\n", len(reports))
	return c.Status(fiber.StatusCreated).SendString("Report received")
}

func main() {
	app := fiber.New()
	dbInit()
	modelInit()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("NetWatcher Logs")
	})
	app.Post("/", handlePost)

	port := getPort()
	log.Fatal(app.Listen(":" + port))
}
