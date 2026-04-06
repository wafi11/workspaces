package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/wafi11/workspaces/pkg/models"
)

type Client struct {
	es *es8.Client
}

func NewClient(address string) (*Client, error) {
	cfg := es8.Config{
		Addresses: []string{address},
	}
	client, err := es8.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{es: client}, nil
}

func (c *Client) QueryLogs(params models.QueryParams) ([]models.LogEntry, []interface{}, error) {
	must := []map[string]interface{}{}

	// PERBAIKAN: Gunakan .keyword untuk filter yang lebih akurat
	if params.Namespace != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"kubernetes.namespace_name.keyword": params.Namespace,
			},
		})
	}

	if params.Pod != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"kubernetes.pod_name.keyword": params.Pod,
			},
		})
	}

	if params.Level != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"level": params.Level,
			},
		})
	}

	if params.Keyword != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  params.Keyword,
				"fields": []string{"log", "message"},
			},
		})
	}

	if len(must) == 0 {
		must = append(must, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	query := map[string]interface{}{

		"size": params.Size,
		"sort": []map[string]interface{}{
			{"@timestamp": "asc"},
			{"_id": "asc"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte": params.From.Format(time.RFC3339),
							},
						},
					},
				},
			},
		},
	}

	if len(params.SearchAfter) > 0 {
		query["search_after"] = params.SearchAfter
	}

	body, _ := json.Marshal(query)

	res, err := c.es.Search(
		c.es.Search.WithIndex("logstash-*"),
		c.es.Search.WithBody(strings.NewReader(string(body))),
	)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return []models.LogEntry{}, nil, nil
	}

	var logs []models.LogEntry
	var lastSort []interface{}

	for _, hit := range hits {
		h := hit.(map[string]interface{})
		source := h["_source"].(map[string]interface{})

		entry := models.LogEntry{}
		if v, ok := source["@timestamp"].(string); ok {
			entry.Timestamp = v
		}
		if v, ok := source["cluster"].(string); ok {
			entry.Cluster = v
		}

		// PERBAIKAN: Parsing objek kubernetes yang nested
		if kube, ok := source["kubernetes"].(map[string]interface{}); ok {
			if v, ok := kube["namespace_name"].(string); ok {
				entry.Namespace = v
			}
			if v, ok := kube["pod_name"].(string); ok {
				entry.PodName = v
			}
			if v, ok := kube["container_name"].(string); ok {
				entry.ContainerName = v
			}
			if v, ok := kube["host"].(string); ok {
				entry.Host = v
			}
		}

		if v, ok := source["log"].(string); ok {
			entry.Log = v
		}
		if v, ok := source["level"].(string); ok {
			entry.Level = v
		}
		if v, ok := source["message"].(string); ok {
			entry.Message = v
		}

		logs = append(logs, entry)
		lastSort = h["sort"].([]interface{})
	}

	return logs, lastSort, nil
}

// Tambahkan perbaikan serupa untuk GetNamespaces dan GetPods
func (c *Client) GetNamespaces() ([]string, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"namespaces": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "kubernetes.namespace_name.keyword",
					"size":  100,
				},
			},
		},
	}

	body, _ := json.Marshal(query)
	res, err := c.es.Search(
		c.es.Search.WithIndex("logstash-*"), // Sesuaikan index
		c.es.Search.WithBody(strings.NewReader(string(body))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	// Pastikan aggregations ada untuk menghindari panic
	if aggs, ok := result["aggregations"].(map[string]interface{}); ok {
		if ns, ok := aggs["namespaces"].(map[string]interface{}); ok {
			buckets := ns["buckets"].([]interface{})
			var namespaces []string
			for _, b := range buckets {
				namespaces = append(namespaces, fmt.Sprintf("%v", b.(map[string]interface{})["key"]))
			}
			return namespaces, nil
		}
	}
	return []string{}, nil
}

func (c *Client) GetPods(namespace string) ([]string, error) {
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"kubernetes.namespace_name.keyword": namespace,
			},
		},
		"aggs": map[string]interface{}{
			"pods": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "kubernetes.pod_name.keyword",
					"size":  100,
				},
			},
		},
	}

	body, _ := json.Marshal(query)
	res, err := c.es.Search(
		c.es.Search.WithIndex("logstash-*"), // Sesuaikan index
		c.es.Search.WithBody(strings.NewReader(string(body))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	if aggs, ok := result["aggregations"].(map[string]interface{}); ok {
		if p, ok := aggs["pods"].(map[string]interface{}); ok {
			buckets := p["buckets"].([]interface{})
			var pods []string
			for _, b := range buckets {
				pods = append(pods, fmt.Sprintf("%v", b.(map[string]interface{})["key"]))
			}
			return pods, nil
		}
	}
	return []string{}, nil
}
