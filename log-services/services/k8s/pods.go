package k8s

import (
	"encoding/json"
	"fmt"
	"strings"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

func GetPods(c *es8.Client, namespace string) ([]string, error) {
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
	res, err := c.Search(
		c.Search.WithIndex("logstash-*"), 
		c.Search.WithBody(strings.NewReader(string(body))),
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
