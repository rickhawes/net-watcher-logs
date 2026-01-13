package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func logInit() string {
	dir, exists := os.LookupEnv("LOG_DIR")
	if !exists {
		dir = "logs"
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v\n", err)
	}
	return dir
}

func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "5555"
	}
	return port
}

func handlePost(c *fiber.Ctx) error {
	c.Accepts("application/json")

	// Read the request body
	bytes := c.Body()

	// Format the timestamp
	formattedTime := time.Now().Local().Format("20060102-150405")

	// Create a timestamped log fil
	filePath := fmt.Sprintf("%s/%s-%03d.json", logDir, formattedTime, logCount%1000)
	err := os.WriteFile(filePath, bytes, 0666)
	if err != nil {
		log.Printf("Error writing log to file: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to write log")
	}
	logCount++
	log.Printf("Log written to: %s\n", filePath)

	return c.Status(fiber.StatusCreated).SendString("Log received")
}

var logDir string = ""
var logCount int = 0

func main() {
	app := fiber.New()
	logDir = logInit()
	log.Printf("Logging directory initialized at: %s\n", logDir)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("NetWatcher Logs")
	})
	app.Post("/", handlePost)

	port := getPort()
	log.Fatal(app.Listen(":" + port))
}
