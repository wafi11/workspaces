package main

import (
	"fmt"
	"os"

	"gihub.com/wafi11/workspaces/log-service/config"
	"gihub.com/wafi11/workspaces/log-service/services/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "https://elasticsearch.wfdnstore.online" 
	}

	es8Client, err := config.NewClient(esURL)
	if err != nil {
		fmt.Printf("Error creating Elasticsearch client: %v\n", err)
		return
	}

	e := echo.New()
	e.GET("/stream", handler.StreamLogs(es8Client))
	e.Logger.Fatal(e.Start(":8080"))
}