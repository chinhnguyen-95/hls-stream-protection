package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ParseFFMPEGOutput parses raw ffmpeg log output into a structured JSON response
func ParseFFMPEGOutput(w http.ResponseWriter, r *http.Request) error {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return fmt.Errorf("invalid request method")
	}

	// Read raw ffmpeg log data from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return fmt.Errorf("failed to read request body: %w", err)
	}

	// Parse ffmpeg log lines
	logLines := strings.Split(string(body), "\n")
	parsedLogs := make([]map[string]string, 0)

	for _, line := range logLines {
		if strings.TrimSpace(line) == "" {
			continue // Skip empty lines
		}

		// Example: Simple key-value extraction
		// Modify this logic to handle specific ffmpeg log formats
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		parsedLogs = append(parsedLogs, map[string]string{key: value})
	}

	// Encode parsed logs as JSON and respond
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(parsedLogs); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
