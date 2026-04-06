package logservice

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/models"
)

func SearchLogs(esClient *config.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")

        params := models.QueryParams{
            Namespace: r.URL.Query().Get("namespace"),
            Pod:       r.URL.Query().Get("pod"),
            Level:     r.URL.Query().Get("level"),
            Keyword:   r.URL.Query().Get("keyword"),
            From:      time.Now().Add(-1 * time.Hour),
            Size:      100,
        }

        logs, _, err := esClient.QueryLogs(params)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]interface{}{
            "data":  logs,
            "total": len(logs),
        })
    }
}