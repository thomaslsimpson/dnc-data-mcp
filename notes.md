# Project Notes

## MCP Service
- Runs on port 8080
- Uses SSH tunnel to connect to database (port 5433 -> ads-prod-reporting-dbi.dnc.io:5432)
- Endpoint: GET http://localhost:8080/mcp/query?q=<url-encoded-sql>
- Returns JSON in format: {"columns": [...], "rows": [...]}

## Agent
- Uses Ollama (llama3.2:latest) for:
  1. Converting natural language to SQL
  2. Interpreting query results
- Requires proper URL encoding of SQL queries
- Needs better schema understanding for complex queries

## Database Schema
- Main table: yer_analysis.yer_reports
- Contains YER (Yield Enhancement Report) data
- Need to document full schema for better query generation

## Known Issues
1. Port conflicts:
   - Need to ensure port 5433 is free before starting MCP service
   - Need to ensure port 8080 is free before starting MCP service
2. Query generation:
   - Ollama needs better understanding of table relationships
   - Need to improve error handling for invalid queries
3. Response interpretation:
   - Need to improve prompts for better result interpretation
   - Need to handle empty results better

## Next Steps
1. Document full database schema
2. Improve Ollama prompts with schema information
3. Add better error handling and recovery
4. Add support for more complex queries
5. Improve response formatting and interpretation

## Commands
```bash
# Start MCP service
cd /Users/tsimpson/github/dnc-data-mcp && GIN_MODE=debug go run main.go

# Run agent
cd /Users/tsimpson/github/dnc-data-mcp/agent && go run main.go "your question here"

# Kill existing processes
pkill -f "go run main.go"
```

## Response Format
The MCP service returns JSON in this format:
```json
{
  "columns": ["column1", "column2", ...],
  "rows": [
    {"column1": "value1", "column2": "value2", ...},
    ...
  ]
}
```

## Error Handling
- Check for port conflicts before starting
- Handle database connection errors
- Handle invalid SQL queries
- Handle empty results
- Handle malformed responses 