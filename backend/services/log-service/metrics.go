package logservices

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type PodMetrics struct {
	PodName     string             `json:"pod_name"`
	AppName     string             `json:"app_name"`
	Containers  []ContainerMetrics `json:"containers"`
	TotalCPU    float64            `json:"total_cpu_cores"`
	TotalMemory int64              `json:"total_memory_mb"`
}

type ContainerMetrics struct {
	Name   string  `json:"name"`
	CPU    float64 `json:"cpu_cores"`
	Memory int64   `json:"memory_mb"`
}

func GetUserWorkspaceMetrics(ctx context.Context, metricsClient *metricsv.Clientset, userID string) ([]PodMetrics, error) {
	namespace := fmt.Sprintf("ws-%s", userID)

	pods, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []PodMetrics
	for _, pod := range pods.Items {
		pm := PodMetrics{
			PodName: pod.Name,
			AppName: extractAppName(pod.Name),
		}
		for _, c := range pod.Containers {
			cpu := float64(c.Usage.Cpu().MilliValue()) / 1000
			mem := c.Usage.Memory().Value() / 1024 / 1024

			pm.Containers = append(pm.Containers, ContainerMetrics{
				Name:   c.Name,
				CPU:    cpu,
				Memory: mem,
			})
			pm.TotalCPU += cpu
			pm.TotalMemory += mem
		}
		result = append(result, pm)
	}

	return result, nil
}

func extractAppName(podName string) string {
	parts := strings.Split(podName, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return podName
}
