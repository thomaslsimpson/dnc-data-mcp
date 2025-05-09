package main

import (
	"net/http"

	"github.com/dnc-data-mcp/config"
	"github.com/dnc-data-mcp/db"
	"github.com/dnc-data-mcp/mcp"
	"github.com/dnc-data-mcp/tunnel"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		panic(err)
	}

	// Set up SSH tunnel
	sshTunnel, err := tunnel.NewSSHTunnel(cfg)
	if err != nil {
		panic(err)
	}
	defer sshTunnel.Close()

	// Create a copy of the config for the database connection
	dbConfig := *cfg
	// Update the database connection to use the tunnel
	dbConfig.Database.ROTraffic.Server = "localhost"
	dbConfig.Database.ROTraffic.Port = 5433

	// Initialize database connection
	database, err := db.NewDB(&dbConfig)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	// Create MCP service
	service := mcp.NewService(database)

	// Set up Gin router
	r := gin.Default()

	// MCP endpoints
	r.GET("/mcp/query", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing query parameter 'q'"})
			return
		}

		resp, err := service.HandleQuery(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	// Start server
	r.Run(":8080")
}
