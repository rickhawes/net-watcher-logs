package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func dbInit() {
	dsn, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err = createTable(); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Database initialized")
}

func createTable() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS reports (
			id               BIGSERIAL PRIMARY KEY,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
			start_time       TEXT NOT NULL,
			duration         DOUBLE PRECISION NOT NULL,			
			state_durations  JSONB NOT NULL DEFAULT '{}',
			action_durations JSONB NOT NULL DEFAULT '{}',
			start_memory     BIGINT NOT NULL,
			delta_memory     BIGINT NOT NULL,
			watcher_duration DOUBLE PRECISION NOT NULL,
			important_events JSONB NOT NULL DEFAULT '[]'
		)
	`)
	return err
}

func insertLogEntry(entry LogEntry) error {
	stateDurations, err := json.Marshal(entry.StateDurations)
	if err != nil {
		return fmt.Errorf("marshal state_durations: %w", err)
	}

	actionDurations, err := json.Marshal(entry.ActionDurations)
	if err != nil {
		return fmt.Errorf("marshal action_durations: %w", err)
	}

	importantEvents, err := json.Marshal(entry.ImportantEvents)
	if err != nil {
		return fmt.Errorf("marshal important_events: %w", err)
	}

	_, err = db.Exec(`
		INSERT INTO reports
			(state_durations, action_durations, watcher_duration, duration,
			 start_memory, delta_memory, important_events, start_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		stateDurations, actionDurations, entry.WatcherDuration, entry.Duration,
		entry.StartMemory, entry.DeltaMemory, importantEvents, entry.StartTime,
	)
	return err
}
