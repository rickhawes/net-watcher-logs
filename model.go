package main

import "time"

type LogEntry struct {
	// Database
	ID              int64              `json:"id"`
	CreatedAt       time.Time          `json:"created_at"`
	
	// Period
	StartTime       string             `json:"start_time"`
	Duration        float64            `json:"duration"`
	
	// Watcher Information
	StartMemory     int64              `json:"start_memory"`
	DeltaMemory     int64              `json:"delta_memory"`
	WatcherDuration float64            `json:"watcher_duration"`
	
	// Summary
	StateDurations  map[string]float64 `json:"state_durations"`
	ActionDurations map[string]float64 `json:"action_durations"`
	
	// Events
	ImportantEvents []Event            `json:"important_events"`
}

type Event struct {
	At        string  `json:"at"`
	EventType string  `json:"event_type"`
	State     string  `json:"state"`
	Message   string  `json:"message"`
	Duration  float64 `json:"duration"`
	Context   []Event `json:"context"`
}
