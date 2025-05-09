package mcp

import (
	"fmt"
	"strings"

	"github.com/dnc-data-mcp/db"
)

// Service represents our MCP service
type Service struct {
	db *db.DB
}

// NewService creates a new MCP service
func NewService(db *db.DB) *Service {
	return &Service{db: db}
}

// QueryResult represents a single row from a query
type QueryResult map[string]interface{}

// QueryResponse represents the response from a query
type QueryResponse struct {
	Columns []string      `json:"columns"`
	Rows    []QueryResult `json:"rows"`
	Error   string        `json:"error,omitempty"`
}

// HandleQuery handles a natural language query and returns the results
func (s *Service) HandleQuery(query string) (*QueryResponse, error) {
	// Convert query to lowercase for easier matching
	queryLower := strings.ToLower(query)

	// Handle special commands
	if strings.HasPrefix(queryLower, "show tables") {
		return s.handleShowTables()
	}

	if strings.HasPrefix(queryLower, "describe table") {
		return s.handleDescribeTable(query)
	}

	// Handle natural language queries
	if strings.Contains(queryLower, "who are our partners") {
		sql := `
			SELECT DISTINCT p.id as partner_id, p.name, p.status, ps.status as status_name
			FROM partners.partners p
			JOIN partners.partner_status ps ON p.status = ps.id
			ORDER BY p.name;
		`
		return s.executeQuery(sql)
	}

	if strings.Contains(queryLower, "which partner made the most money last month") {
		sql := `
			WITH last_month AS (
				SELECT *
				FROM yer_analysis.yer_reports
				WHERE date_trunc('month', end_date) = date_trunc('month', current_date - interval '1 month')
			)
			SELECT 
				p.id as partner_id,
				p.name as partner_name,
				yi.start_date,
				yi.end_date,
				sum(yi.amount) as total_revenue
			FROM yer_analysis.v_yer_items yi
			JOIN partners.partners p ON yi.partner_id = p.id
			JOIN last_month r ON yi.yer_report_id = r.id
			GROUP BY p.id, p.name, yi.start_date, yi.end_date
			ORDER BY total_revenue DESC
			LIMIT 5;
		`
		return s.executeQuery(sql)
	}

	if strings.Contains(queryLower, "which source tags does") && strings.Contains(queryLower, "use") {
		// Extract partner name from query
		parts := strings.Split(queryLower, "which source tags does")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid query format")
		}
		partnerPart := strings.Split(parts[1], "use")[0]
		partnerName := strings.TrimSpace(partnerPart)

		sql := `
			SELECT DISTINCT 
				st.id as sourcetag_id,
				st.name as sourcetag_name
			FROM yer_analysis.v_yer_items yi
			JOIN partners.partners p ON yi.partner_id = p.id
			JOIN yer_analysis.source_tags st ON yi.sourcetag_id = st.id
			WHERE lower(p.name) = lower($1)
			ORDER BY st.name;
		`
		return s.executeQuery(sql, partnerName)
	}

	if strings.Contains(queryLower, "which source tags went up in tq last month") {
		sql := `
			WITH last_month AS (
				SELECT *
				FROM yer_analysis.yer_reports
				WHERE date_trunc('month', end_date) = date_trunc('month', current_date - interval '1 month')
			),
			prev_month AS (
				SELECT *
				FROM yer_analysis.yer_reports
				WHERE date_trunc('month', end_date) = date_trunc('month', current_date - interval '2 month')
			),
			last_month_data AS (
				SELECT 
					st.id as sourcetag_id,
					st.name as sourcetag_name,
					avg(yi.traffic_quality) as traffic_quality
				FROM yer_analysis.v_yer_items yi
				JOIN yer_analysis.source_tags st ON yi.sourcetag_id = st.id
				JOIN last_month lm ON yi.yer_report_id = lm.id
				GROUP BY st.id, st.name
			),
			prev_month_data AS (
				SELECT 
					st.id as sourcetag_id,
					avg(yi.traffic_quality) as traffic_quality
				FROM yer_analysis.v_yer_items yi
				JOIN yer_analysis.source_tags st ON yi.sourcetag_id = st.id
				JOIN prev_month pm ON yi.yer_report_id = pm.id
				GROUP BY st.id
			)
			SELECT 
				curr.sourcetag_name,
				curr.traffic_quality as current_tq,
				prev.traffic_quality as previous_tq,
				(curr.traffic_quality - prev.traffic_quality) as tq_change
			FROM last_month_data curr
			JOIN prev_month_data prev ON curr.sourcetag_id = prev.sourcetag_id
			WHERE curr.traffic_quality > prev.traffic_quality
			ORDER BY tq_change DESC
			LIMIT 10;
		`
		return s.executeQuery(sql)
	}

	if strings.Contains(queryLower, "which partners have been on the yer the most in the last 6 months") {
		sql := `
			SELECT 
				p.id as partner_id,
				p.name as partner_name,
				count(distinct yi.yer_report_id) as report_count
			FROM yer_analysis.v_yer_items yi
			JOIN partners.partners p ON yi.partner_id = p.id
			JOIN yer_analysis.yer_reports r ON yi.yer_report_id = r.id
			WHERE r.end_date >= current_date - interval '6 months'
			GROUP BY p.id, p.name
			ORDER BY report_count DESC
			LIMIT 10;
		`
		return s.executeQuery(sql)
	}

	if strings.Contains(queryLower, "what traffic sources does") && strings.Contains(queryLower, "use") {
		// Extract partner name from query
		parts := strings.Split(queryLower, "what traffic sources does")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid query format")
		}
		partnerPart := strings.Split(parts[1], "use")[0]
		partnerName := strings.TrimSpace(partnerPart)

		sql := `
			WITH last_month AS (
				SELECT *
				FROM yer_analysis.yer_reports
				WHERE date_trunc('month', end_date) = date_trunc('month', current_date - interval '1 month')
			)
			SELECT DISTINCT 
				p.id as partner_id,
				p.name as partner_name,
				st.id as sourcetag_id,
				st.name as sourcetag_name,
				sum(yi.total_searches) as total_searches,
				sum(yi.total_clicks) as total_clicks,
				avg(yi.traffic_quality) as traffic_quality
			FROM yer_analysis.v_yer_items yi
			JOIN partners.partners p ON yi.partner_id = p.id
			JOIN yer_analysis.source_tags st ON yi.sourcetag_id = st.id
			JOIN last_month r ON yi.yer_report_id = r.id
			WHERE lower(p.name) = lower($1)
			GROUP BY p.id, p.name, st.id, st.name
			ORDER BY total_searches DESC;
		`
		return s.executeQuery(sql, partnerName)
	}

	// If no natural language query matches, treat as direct SQL
	return s.executeQuery(query)
}

// handleShowTables returns a list of all tables in the database
func (s *Service) handleShowTables() (*QueryResponse, error) {
	query := `
		SELECT table_schema, table_name 
		FROM information_schema.tables 
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_schema, table_name
	`
	return s.executeQuery(query)
}

// handleDescribeTable returns the structure of a specific table
func (s *Service) handleDescribeTable(query string) (*QueryResponse, error) {
	// Extract table name from query
	parts := strings.Fields(query)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid describe table query format")
	}
	fullTableName := parts[2]

	// Split schema and table name
	var schema, tableName string
	if strings.Contains(fullTableName, ".") {
		schemaParts := strings.Split(fullTableName, ".")
		schema = schemaParts[0]
		tableName = schemaParts[1]
	} else {
		schema = "public"
		tableName = fullTableName
	}

	// Query to describe table in PostgreSQL
	sql := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = $1
		AND table_name = $2
		ORDER BY ordinal_position;
	`

	return s.executeQuery(sql, schema, tableName)
}

// handleDirectQuery executes a direct SQL query
func (s *Service) handleDirectQuery(query string) (*QueryResponse, error) {
	return s.executeQuery(query)
}

// executeQuery executes a query and returns the results
func (s *Service) executeQuery(query string, args ...interface{}) (*QueryResponse, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return &QueryResponse{Error: err.Error()}, nil
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return &QueryResponse{Error: err.Error()}, nil
	}

	// Create a slice to hold the values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Process rows
	var results []QueryResult
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return &QueryResponse{Error: err.Error()}, nil
		}

		// Convert values to a map
		row := make(QueryResult)
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case []byte:
				row[col] = string(v)
			case nil:
				row[col] = nil
			default:
				row[col] = v
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return &QueryResponse{Error: err.Error()}, nil
	}

	return &QueryResponse{
		Columns: columns,
		Rows:    results,
	}, nil
}
