package mcp

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dnc-data-mcp/config"
	"github.com/dnc-data-mcp/db"
	"github.com/dnc-data-mcp/tunnel"
)

func TestMCPService(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

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

	// Initialize database connection
	database, err := db.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create MCP service
	service := NewService(database)

	// Test cases
	testCases := []struct {
		name     string
		query    string
		validate func(*testing.T, *QueryResponse)
	}{
		{
			name:  "Show Tables",
			query: "show tables",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Columns) != 2 {
					t.Errorf("Expected 2 columns, got %d", len(resp.Columns))
				}
				if len(resp.Rows) == 0 {
					t.Error("Expected at least one table")
				}
				// Print the tables for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("Tables:\n%s", string(jsonData))
			},
		},
		{
			name:  "Describe Table",
			query: "describe table yer_analysis.yer_reports",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Columns) != 3 {
					t.Errorf("Expected 3 columns, got %d", len(resp.Columns))
				}
				if len(resp.Rows) == 0 {
					t.Error("Expected at least one column")
				}
				// Print the table structure for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("Table structure:\n%s", string(jsonData))
			},
		},
		{
			name:  "Query Data",
			query: "SELECT * FROM yer_analysis.yer_reports LIMIT 2",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Rows) != 2 {
					t.Errorf("Expected 2 rows, got %d", len(resp.Rows))
				}
				// Print the data for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("Query results:\n%s", string(jsonData))
			},
		},
		{
			name:  "List Partners",
			query: "who are our partners",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Columns) != 4 {
					t.Errorf("Expected 4 columns, got %d", len(resp.Columns))
				}
				if len(resp.Rows) == 0 {
					t.Error("Expected at least one partner")
				}
				// Print the partners for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("Partners:\n%s", string(jsonData))
			},
		},
		{
			name:  "Top Revenue Partners",
			query: "which partner made the most money last month",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Columns) != 5 {
					t.Errorf("Expected 5 columns, got %d", len(resp.Columns))
				}
				if len(resp.Rows) == 0 {
					t.Error("Expected at least one partner")
				}
				// Print the revenue data for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("Revenue data:\n%s", string(jsonData))
			},
		},
		{
			name:  "Partner Source Tags",
			query: "which source tags does DNC use",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Columns) != 2 {
					t.Errorf("Expected 2 columns, got %d", len(resp.Columns))
				}
				// Print the source tags for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("Source tags:\n%s", string(jsonData))
			},
		},
		{
			name:  "TQ Improvement",
			query: "which source tags went up in tq last month",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Columns) != 4 {
					t.Errorf("Expected 4 columns, got %d", len(resp.Columns))
				}
				// Print the TQ changes for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("TQ changes:\n%s", string(jsonData))
			},
		},
		{
			name:  "YER Frequency",
			query: "which partners have been on the yer the most in the last 6 months",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Columns) != 3 {
					t.Errorf("Expected 3 columns, got %d", len(resp.Columns))
				}
				// Print the YER frequency for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("YER frequency:\n%s", string(jsonData))
			},
		},
		{
			name:  "Partner Traffic Sources",
			query: "what traffic sources does DNC use",
			validate: func(t *testing.T, resp *QueryResponse) {
				if resp.Error != "" {
					t.Errorf("Unexpected error: %s", resp.Error)
				}
				if len(resp.Columns) != 7 {
					t.Errorf("Expected 7 columns, got %d", len(resp.Columns))
				}
				// Print the traffic sources for debugging
				jsonData, _ := json.MarshalIndent(resp, "", "  ")
				t.Logf("Traffic sources:\n%s", string(jsonData))
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.HandleQuery(tc.query)
			if err != nil {
				t.Fatalf("Failed to handle query: %v", err)
			}
			tc.validate(t, resp)
		})
	}
}

func TestDescribePartners(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

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

	// Initialize database connection
	database, err := db.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create MCP service
	service := NewService(database)

	// Test describing partners table
	resp, err := service.HandleQuery("describe table partners.partners")
	if err != nil {
		t.Fatalf("Failed to describe partners table: %v", err)
	}

	// Print the table structure
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Partners table structure:\n%s", string(jsonData))
}

func TestDescribePartnerStatus(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

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

	// Initialize database connection
	database, err := db.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create MCP service
	service := NewService(database)

	// Test describing partner_status table
	resp, err := service.HandleQuery("describe table partners.partner_status")
	if err != nil {
		t.Fatalf("Failed to describe partner_status table: %v", err)
	}

	// Print the table structure
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Partner status table structure:\n%s", string(jsonData))
}

func TestDescribeYerData(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

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

	// Initialize database connection
	database, err := db.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create MCP service
	service := NewService(database)

	// Test describing yer_data table
	resp, err := service.HandleQuery("describe table yer_analysis.yer_data")
	if err != nil {
		t.Fatalf("Failed to describe yer_data table: %v", err)
	}

	// Print the table structure
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("YER data table structure:\n%s", string(jsonData))
}

func TestDescribeYerDataDirect(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

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

	// Initialize database connection
	database, err := db.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create MCP service
	service := NewService(database)

	// Test describing yer_data table with a direct query
	resp, err := service.HandleQuery(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'yer_analysis'
		AND table_name = 'yer_data'
		ORDER BY ordinal_position;
	`)
	if err != nil {
		t.Fatalf("Failed to describe yer_data table: %v", err)
	}

	// Print the table structure
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("YER data table structure (direct query):\n%s", string(jsonData))
}

func TestDescribeYerItems(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

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

	// Initialize database connection
	database, err := db.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create MCP service
	service := NewService(database)

	// Test describing v_yer_items view
	resp, err := service.HandleQuery(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'yer_analysis'
		AND table_name = 'v_yer_items'
		ORDER BY ordinal_position;
	`)
	if err != nil {
		t.Fatalf("Failed to describe v_yer_items view: %v", err)
	}

	// Print the view structure
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("YER items view structure:\n%s", string(jsonData))
}

func TestDescribeYerItemsTags(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

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

	// Initialize database connection
	database, err := db.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create MCP service
	service := NewService(database)

	// Test describing v_yer_items_tags view
	resp, err := service.HandleQuery(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'yer_analysis'
		AND table_name = 'v_yer_items_tags'
		ORDER BY ordinal_position;
	`)
	if err != nil {
		t.Fatalf("Failed to describe v_yer_items_tags view: %v", err)
	}

	// Print the view structure
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("YER items tags view structure:\n%s", string(jsonData))
}
