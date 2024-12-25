package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"hls-stream-protection/redis"
)

// ProtectHLSStream validates HLS stream access using GET parameters
func ProtectHLSStream(db *sql.DB, redisClient *redis.Client, w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	// Extract required parameters from the request
	streamID := r.URL.Query().Get("stream_id")
	accessToken := r.URL.Query().Get("access_token")
	if streamID == "" || accessToken == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return fmt.Errorf("missing required parameters")
	}

	// Check Redis for access token validation
	cachedToken, err := redisClient.Get(ctx, streamID)
	if err != nil {
		return fmt.Errorf("failed to check Redis: %w", err)
	}

	if cachedToken == accessToken {
		// Token is valid
		response := map[string]string{"status": "success", "message": "Stream access granted"}
		_ = json.NewEncoder(w).Encode(response)
		return nil
	}

	// Check database for stream access
	query := "SELECT access_token FROM streams WHERE stream_id = ?"
	row := db.QueryRow(query, streamID)
	var dbToken string
	if err := row.Scan(&dbToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Stream not found", http.StatusNotFound)
			return fmt.Errorf("stream not found: %w", err)
		}
		return fmt.Errorf("failed to query database: %w", err)
	}

	if dbToken != accessToken {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return fmt.Errorf("invalid access token")
	}

	// Cache the token in Redis for future requests
	if err := redisClient.Set(ctx, streamID, accessToken, 0); err != nil {
		return fmt.Errorf("failed to cache token in Redis: %w", err)
	}

	// Respond with success
	response := map[string]string{"status": "success", "message": "Stream access granted"}
	_ = json.NewEncoder(w).Encode(response)
	return nil
}
