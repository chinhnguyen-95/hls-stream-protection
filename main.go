package main

import (
	"database/sql"
	"log"
	"net/http"

	"hls-stream-protection/api"
	"hls-stream-protection/config"
	"hls-stream-protection/db"
	"hls-stream-protection/redis"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func(dbConn *sql.DB) {
		_ = dbConn.Close()
	}(dbConn)

	// Initialize Redis connection
	redisClient, err := redis.NewClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func(redisClient *redis.Client) {
		_ = redisClient.Close()
	}(redisClient)

	// Define API routes
	http.HandleFunc(
		"/hls/protect", func(w http.ResponseWriter, r *http.Request) {
			err := api.ProtectHLSStream(dbConn, redisClient, w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		},
	)

	http.HandleFunc(
		"/ffmpeg/parse", func(w http.ResponseWriter, r *http.Request) {
			err := api.ParseFFMPEGOutput(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		},
	)

	// Start HTTP server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
