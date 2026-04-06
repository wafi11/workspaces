package k8s

import "time"

type LogEntry struct {
    Timestamp      string `json:"@timestamp"`
    Cluster        string `json:"cluster"`
    Namespace      string `json:"kubernetes.namespace_name"`
    PodName        string `json:"kubernetes.pod_name"`
    ContainerName  string `json:"kubernetes.container_name"`
    Host           string `json:"kubernetes.host"`
    Log            string `json:"log"`
    Level          string `json:"level"`
    Message        string `json:"message"`
}

type QueryParams struct {
    Namespace string
    Pod       string
    Level     string
    Keyword   string
    From      time.Time
    Size      int
    SearchAfter []interface{}
}