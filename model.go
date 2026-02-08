package main

import (
	"log"
	"os"
	"strings"
	"time"
)

const customLayout = "2006-01-02 15:04:05"

// Custom formated time type for JSON serialization
type WatcherTime time.Time

func (ct *WatcherTime) UnmarshalJSON(b []byte) (err error) {
	// Trim surrounding quotes from the JSON string
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		*ct = WatcherTime(time.Time{})
		return nil
	}
	// Parse the time using the custom layout
	ts, err := time.ParseInLocation(customLayout, s, location)
	if err != nil {
		return err
	}
	*ct = WatcherTime(ts)
	return nil
}

func (ct WatcherTime) MarshalJSON() ([]byte, error) {
	// Format the time using the custom layout and add quotes for JSON string
	s := time.Time(ct).Format(customLayout)
	return []byte(`"` + s + `"`), nil
}

var location *time.Location

func modelInit() {
	tz, exists := os.LookupEnv("TZ")
	if !exists {
		tz = "America/Los_Angeles"
	}
	var err error
	if location, err = time.LoadLocation(tz); err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
}

type Report struct {
	// Database
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	// Period
	StartTime WatcherTime `json:"start_time"`
	Duration  float64     `json:"duration"`

	// Watcher Information
	StartMemory     int64   `json:"start_memory"`
	DeltaMemory     int64   `json:"delta_memory"`
	WatcherDuration float64 `json:"watcher_duration"`

	// Summary
	StateDurations  map[string]float64 `json:"state_durations"`
	ActionDurations map[string]float64 `json:"action_durations"`

	// Events
	ImportantEvents []Event `json:"important_events"`
}

type Event struct {
	At        WatcherTime `json:"at"`
	EventType string      `json:"event_type"`
	State     string      `json:"state"`
	Message   string      `json:"message"`
	Duration  float64     `json:"duration"`
	Context   []Event     `json:"context"`
}
