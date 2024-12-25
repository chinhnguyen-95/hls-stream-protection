# HLS Stream Protection

## Overview
This project demonstrates modifying an open-source project to include:
- Integration with PostgreSQL/MySQL and Redis.
- HLS stream protection using GET parameters.
- ffmpeg log parsing into JSON responses.
- Migration of JSON-based data to a proper database table.

## Features
- **Database Integration**: Supports PostgreSQL and MySQL for stream validation.
- **Redis Caching**: Optimizes token validation performance.
- **HLS Protection**: Ensures secure access to HLS streams via API endpoints.
- **ffmpeg Log Parsing**: Converts raw ffmpeg output into structured JSON.
- **Data Migration**: Moves JSON-based records into relational databases.

## Directory Structure
```
.
├── api/                  # API logic for HLS protection and ffmpeg parsing
│   ├── hls_protection.go # Protect HLS streams
│   └── ffmpeg_parser.go  # Parse ffmpeg logs
├── config/               # Configuration management
│   └── config.go         # Load environment-based configuration
├── db/                   # Database logic
│   ├── db.go             # Connect to PostgreSQL/MySQL
│   └── migration.go      # Migrate JSON data to database tables
├── redis/                # Redis client logic
│   └── redis_client.go   # Handle Redis operations
├── main.go               # Entry point of the application
├── go.mod                # Go module dependencies
├── go.sum                # Dependency checksums
```

## Prerequisites
- **Go**: Version 1.23+
- **PostgreSQL/MySQL**: A running instance of either database.
- **Redis**: A running Redis instance.

## Configuration
### Environment Variables
Set the following environment variables:
- `SERVER_ADDRESS`: Address to run the HTTP server (default `:8080`).
- `DATABASE_URL`: URL for the database connection (e.g., `postgres://user:password@localhost/dbname?sslmode=disable`).
- `REDIS_URL`: URL for the Redis connection (default `redis://localhost:6379`).

### Default Configuration
Default values are defined in `config/config.go`. Update as needed.

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/chinhnguyen-95/hls-stream-protection.git
   cd hls-stream-protection
   ```
2. Initialize Go modules:
   ```bash
   go mod tidy
   ```

## Usage
1. Start the application:
   ```bash
   go run main.go
   ```
2. Available endpoints:
    - **HLS Protection**: `POST /hls/protect` (validate HLS stream access).
    - **ffmpeg Parsing**: `POST /ffmpeg/parse` (convert ffmpeg logs to JSON).

## Example Requests
### HLS Protection
Request:
```bash
curl -X POST "http://localhost:8080/hls/protect?stream_id=abc123&access_token=token123"
```
Response:
```json
{
  "status": "success",
  "message": "Stream access granted"
}
```

### ffmpeg Parsing
Request:
```bash
curl -X POST -d @ffmpeg.log "http://localhost:8080/ffmpeg/parse"
```
Response:
```json
[
  {"key": "value"},
  {"key": "value"}
]
```