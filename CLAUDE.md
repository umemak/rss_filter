# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an RSS filter service built in Go that aggregates and filters RSS feeds from various sources. It's designed to run as a Cloud Run service that:

1. Fetches RSS feeds from multiple configured sources
2. Filters content based on configurable criteria (title, description, URL patterns, etc.)
3. Merges multiple feeds into categorized collections
4. Returns filtered RSS output

## Architecture

- **Package**: `main` - Single package containing all functionality
- **Main Handler**: `Handler()` function serves HTTP requests with query parameters:
  - `?name=<config_name>` - Returns merged RSS for a specific configuration
  - `?url=<rss_url>` - Returns filtered RSS for a single URL
- **Configuration**: RSS sources and filter rules are defined in `config.go`
- **Core Logic**: 
  - `getRSS()` - Fetches and parses RSS feeds using `gofeed`
  - `isIgnore()` - Applies filtering logic based on multiple filter arrays
  - `GetByName()` / `GetByURL()` - Main processing functions

## Key Dependencies

- `github.com/gorilla/feeds` - RSS feed generation
- `github.com/mmcdole/gofeed` - RSS feed parsing

## Configuration Structure

The `configs` array contains RSS source configurations with:
- `name` - Configuration identifier
- `rsss` - Array of RSS sources (title, URL)
- `ignores` - Array of ignore rules (condition, text)

Filter arrays (`filter1` through `filter6`) define different filtering strategies:
- Content filtering (title, description, content)
- URL pattern filtering
- Author filtering
- Title replacement rules

## Development Commands

### Build and Run
```bash
go build                    # Build the package
go run .                    # Run locally (requires HTTP server setup)
go test                     # Run tests
go mod tidy                 # Clean up dependencies
```

### Deployment
Deploy to Cloud Run:
```bash
gcloud run deploy rss-filter --source . --platform managed --region us-central1 --allow-unauthenticated
```

#### Prerequisites
- Ensure Cloud Run API is enabled:
```bash
gcloud services enable run.googleapis.com
```

#### Docker Support
The service includes a Dockerfile for containerized deployment using Go 1.25.0.

### Configuration Management
- Edit `config.go` to modify RSS sources and filtering rules
- Copy `config.go.sample` to `config.go` for initial setup
- The `config.go` file is gitignored to allow local customization

## Development Notes

- RSS configurations support Japanese sources and include extensive filtering for unwanted content
- Filter conditions are defined as constants (TITLE_INCLUDES, URL_INCLUDES, etc.)
- All RSS items are sorted by creation date before output
- The service adds source attribution to RSS descriptions