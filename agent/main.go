package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

type MCPResponse struct {
	Columns []string                 `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
	Error   string                   `json:"error,omitempty"`
}

func main() {
	// Start Ollama if it's not running
	startOllama()

	// Get the query from command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Please provide a query as a command line argument")
		os.Exit(1)
	}
	question := os.Args[1]

	// Create prompt for Ollama to convert question to SQL
	prompt := fmt.Sprintf(`You are a helpful assistant that can convert natural language questions to SQL queries.
The database has tables in various schemas including yer_analysis, trafficdata, and others.
The main table we are querying is yer_analysis.yer_reports which contains YER (Yield Enhancement Report) data.

Please convert the following question into a SQL query:

Question: %s

Please provide ONLY the SQL query, nothing else.`, question)

	// Call Ollama to convert question to SQL
	sqlQuery, err := callOllama(prompt)
	if err != nil {
		fmt.Printf("Error calling Ollama: %v\n", err)
		return
	}

	fmt.Printf("\nGenerated SQL Query:\n%s\n\n", sqlQuery)

	// Make request to our MCP service with the SQL query
	encodedQuery := url.QueryEscape(sqlQuery)
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/mcp/query?q=%s", encodedQuery))
	if err != nil {
		fmt.Printf("Error calling MCP service: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Debug: Print the raw response
	fmt.Printf("\nRaw MCP Response:\n%s\n\n", string(body))

	// Trim any trailing characters and parse the MCP response
	bodyStr := strings.TrimSpace(string(body))
	var mcpResp MCPResponse
	if err := json.Unmarshal([]byte(bodyStr), &mcpResp); err != nil {
		fmt.Printf("Error parsing MCP response: %v\n", err)
		return
	}

	// Check for errors
	if mcpResp.Error != "" {
		fmt.Printf("\nMCP Service Error:\n%s\n", mcpResp.Error)
		return
	}

	// Check if we have any data
	if len(mcpResp.Rows) == 0 {
		fmt.Printf("\nNo data found.\n")
		return
	}

	// Create prompt for Ollama to interpret the results
	resultPrompt := fmt.Sprintf(`Here is the response from the MCP service for the question "%s":

The query returned the following data:
{
  "columns": [%s],
  "rows": %s
}

Please provide a simple, direct answer to the original question based on this data.`,
		question,
		strings.Join(mcpResp.Columns, ", "),
		prettyPrintRows(mcpResp.Rows))

	// Call Ollama to interpret the results
	interpretation, err := callOllama(resultPrompt)
	if err != nil {
		fmt.Printf("Error calling Ollama: %v\n", err)
		return
	}

	fmt.Printf("\nMCP Service Results:\n%s\n", interpretation)
}

func startOllama() {
	// Check if Ollama is running
	_, err := http.Get("http://localhost:11434")
	if err == nil {
		return // Ollama is already running
	}

	// Start Ollama
	cmd := exec.Command("ollama", "serve")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting Ollama: %v\n", err)
		os.Exit(1)
	}
}

func callOllama(prompt string) (string, error) {
	// Prepare the request to Ollama
	reqBody := map[string]interface{}{
		"model":  "llama3.2:latest",
		"prompt": prompt,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Make the request to Ollama
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response line by line
	decoder := json.NewDecoder(resp.Body)
	var fullResponse string

	for decoder.More() {
		var response OllamaResponse
		if err := decoder.Decode(&response); err != nil {
			return "", fmt.Errorf("error decoding response: %v", err)
		}
		fullResponse += response.Response
	}

	return fullResponse, nil
}

func prettyPrintRows(rows []map[string]interface{}) string {
	rowsJSON, err := json.MarshalIndent(rows, "    ", "  ")
	if err != nil {
		return fmt.Sprintf("%v", rows)
	}
	return string(rowsJSON)
}
