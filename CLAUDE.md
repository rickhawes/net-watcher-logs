# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

NetWatcher Logs is a lightweight Go HTTP server that receives JSON log data from the NetWatcher monitoring system and persists it to a Postgres database. It uses the Fiber/v2 web framework (built on fasthttp) and `lib/pq` for Postgres.

## Build & Run Commands

- `make init` - Download Go module dependencies
- `make run` - Run locally with `go run .` (serves on port 5555 by default)
- `make build` - Auto-increment git tag via `autotag.sh`, then build Docker images (AMD64 and multi-platform)
- `make push` - Push Docker images to Docker Hub (`rickhawes/net-watcher-logs`)
- `go build -o app .` - Build the binary directly

There are no tests. Manual testing can be done with `examples.http` (REST client format) by POSTing JSON to `localhost:5555`.

## Architecture

Three-file Go application:

- `main.go` - Fiber app setup, HTTP handlers, entry point
- `model.go` - `LogEntry` and `Event` structs matching the JSON payload
- `db.go` - Postgres connection (`dbInit`), table auto-creation, `insertLogEntry`

Two endpoints:

- `GET /` - Health check, returns "NetWatcher Logs"
- `POST /` - Accepts JSON body (single object or array of objects), inserts each entry into `log_entries` table

The `log_entries` table uses scalar columns for queryable fields (`duration`, `watcher_duration`, `start_memory`, `delta_memory`) and JSONB for variable-key maps (`state_durations`, `action_durations`) and nested structures (`important_events`).

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | *(required)* | Postgres connection string |
| `PORT` | `5555` | Server listen port |

## Docker

Multi-stage build: `golang:1.25-alpine` for building, `alpine:latest` for runtime. Uses `America/Los_Angeles` timezone. `DATABASE_URL` must be provided at runtime. Docker requires a running daemon (e.g., OrbStack).
