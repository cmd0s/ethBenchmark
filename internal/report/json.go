package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FormatJSON generates a JSON string of the report
func FormatJSON(r *Report) (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal report: %w", err)
	}
	return string(data), nil
}

// SaveJSON saves the report as a JSON file with timestamp in filename
func SaveJSON(r *Report, outputDir string) (string, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("ethbench-%s.json", timestamp)
	filepath := filepath.Join(outputDir, filename)

	// Marshal report to JSON
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal report: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write report file: %w", err)
	}

	return filepath, nil
}
