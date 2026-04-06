package logservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/models"
)

// services/log-service/handler.go

func StreamLogs(esClient *config.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		w := c.Response().Writer
		r := c.Request()

		// 1. Set SSE Headers lewat Echo Context
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		// 2. Pastikan server mendukung Flushing
		flusher, ok := w.(http.Flusher)
		if !ok {
			return c.String(http.StatusInternalServerError, "Streaming not supported")
		}

		// 3. Ambil Query Params menggunakan Echo
		params := models.QueryParams{
			Namespace: c.QueryParam("namespace"),
			Pod:       c.QueryParam("pod"),
			Level:     c.QueryParam("level"),
			From:      time.Now().Add(-1 * time.Minute),
			Size:      50,
		}

		var searchAfter []interface{}
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		heartbeat := time.NewTicker(30 * time.Second)
		defer heartbeat.Stop()

		// Loop SSE
		for {
			select {
			case <-r.Context().Done():
				// Koneksi terputus (user close browser)
				return nil

			case <-heartbeat.C:
				fmt.Fprintf(w, ": heartbeat\n\n")
				flusher.Flush()

			case <-ticker.C:
				params.SearchAfter = searchAfter
				logs, lastSort, err := esClient.QueryLogs(params)

				if err != nil {
					fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
					flusher.Flush()
					continue
				}

				for _, log := range logs {
					data, _ := json.Marshal(log)
					// Format SSE: "data: <json>\n\n"
					if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
						return nil
					}
					flusher.Flush()
				}

				if len(lastSort) > 0 {
					searchAfter = lastSort
				}
			}
		}
	}
}
