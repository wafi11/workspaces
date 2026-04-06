package logservice

import (
	"encoding/json"
	"net/http"

	"github.com/wafi11/workspaces/config"
)

func GetNamespaces(esClient *config.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")

        namespaces, err := esClient.GetNamespaces()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]interface{}{
            "data": namespaces,
        })
    }
}

func GetPods(esClient *config.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")

        namespace := r.URL.Query().Get("namespace")
        if namespace == "" {
            http.Error(w, "namespace required", http.StatusBadRequest)
            return
        }

        pods, err := esClient.GetPods(namespace)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]interface{}{
            "data": pods,
        })
    }
}