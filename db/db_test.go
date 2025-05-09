package db

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/dnc-data-mcp/config"
	"github.com/dnc-data-mcp/tunnel"
)

func TestDatabaseConnection(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Debug: Print the entire config structure
	configJSON, _ := json.MarshalIndent(cfg, "", "  ")
	t.Logf("Loaded configuration:\n%s", string(configJSON))

	// Set up SSH tunnel
	sshTunnel, err := tunnel.NewSSHTunnel(cfg)
	if err != nil {
		t.Fatalf("Failed to create SSH tunnel: %v", err)
	}
	defer sshTunnel.Close()

	// Wait a moment for the tunnel to be ready
	time.Sleep(2 * time.Second)

	// Update the database connection to use the tunnel
	cfg.Database.ROTraffic.Server = "localhost"
	cfg.Database.ROTraffic.Port = 5433

	// Debug: Print the DSN string (without password)
	dsn := cfg.GetDSN()
	t.Logf("DSN (without password): %s", dsn)

	// Initialize database connection
	db, err := NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test query to yer_analysis.yer_reports
	rows, err := db.Query("SELECT * FROM yer_analysis.yer_reports LIMIT 5")
	if err != nil {
		t.Fatalf("Failed to query yer_analysis.yer_reports: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("Failed to get column names: %v", err)
	}

	// Log the column names for debugging
	t.Logf("Found columns: %v", columns)

	// Create a slice to hold the values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Print each row
	rowCount := 0
	for rows.Next() {
		rowCount++
		err := rows.Scan(valuePtrs...)
		if err != nil {
			t.Fatalf("Error scanning row: %v", err)
		}

		// Convert values to strings for printing
		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case []byte:
				rowData[col] = string(v)
			case time.Time:
				rowData[col] = v.Format(time.RFC3339)
			case nil:
				rowData[col] = nil
			default:
				rowData[col] = v
			}
		}

		// Pretty print the row
		rowJSON, _ := json.MarshalIndent(rowData, "", "  ")
		t.Logf("Row %d:\n%s", rowCount, string(rowJSON))
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("Error while iterating rows: %v", err)
	}

	t.Logf("Successfully retrieved %d rows from yer_analysis.yer_reports", rowCount)
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
