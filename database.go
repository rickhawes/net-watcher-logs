package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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
			start_time       TIMESTAMPTZ NOT NULL,
			duration         DOUBLE PRECISION NOT NULL,			
			state_durations  JSONB NOT NULL DEFAULT '{}',
			action_durations JSONB NOT NULL DEFAULT '{}',
			start_memory     BIGINT NOT NULL,
			delta_memory     BIGINT NOT NULL,
			watcher_duration DOUBLE PRECISION NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create reports table: %w", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_reports_start_time ON reports(start_time)`)
	if err != nil {
		return fmt.Errorf("create start_time index on reports: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id               BIGSERIAL PRIMARY KEY,
			report_id        BIGINT NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
			at               TIMESTAMPTZ NOT NULL,
			event_type       TEXT NOT NULL,
			state            TEXT NOT NULL,
			message          TEXT NOT NULL,
			duration         DOUBLE PRECISION NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create events table: %w", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_events_report_id ON events(report_id)`)
	if err != nil {
		return fmt.Errorf("create report_id index on events: %w", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_events_at ON events(at)`)
	if err != nil {
		return fmt.Errorf("create at index on events: %w", err)
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_events_unique ON events(at, event_type, state, message, duration)`)
	if err != nil {
		return fmt.Errorf("create unique index on events: %w", err)
	}

	return nil
}

func insertReport(report Report) error {
	stateDurations, err := json.Marshal(report.StateDurations)
	if err != nil {
		return fmt.Errorf("marshal state_durations: %w", err)
	}

	actionDurations, err := json.Marshal(report.ActionDurations)
	if err != nil {
		return fmt.Errorf("marshal action_durations: %w", err)
	}

	result := db.QueryRow(`
		INSERT INTO reports
			(state_durations, action_durations, watcher_duration, duration,
			 start_memory, delta_memory, start_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		stateDurations, actionDurations, report.WatcherDuration, report.Duration,
		report.StartMemory, report.DeltaMemory, time.Time(report.StartTime),
	)
	var reportID int64
	if err := result.Scan(&reportID); err != nil {
		return fmt.Errorf("scan report ID: %w", err)
	}

	for _, important_event := range report.ImportantEvents {
		for _, event := range important_event.Context {
			if err := insertEvent(event, reportID); err != nil {
				return err
			}
		}
		if err := insertEvent(important_event, reportID); err != nil {
			return err
		}
	}
	return nil
}

func insertEvent(event Event, reportID int64) error {
	_, err := db.Exec(`
		INSERT INTO events (report_id, at, event_type, state, message, duration)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (at, event_type, state, message, duration) DO NOTHING
	`, reportID, time.Time(event.At), event.EventType, event.State, event.Message, event.Duration)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}
