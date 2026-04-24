package logservices

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type StorageMetrics struct {
	PodName string        `json:"pod_name"`
	Volumes []VolumeStats `json:"volumes"`
}

type VolumeStats struct {
	UsedMB      int64 `json:"used_mb"`
	CapacityMB  int64 `json:"capacity_mb"`
	AvailableMB int64 `json:"available_mb"`
}

func GetPodStorageFromKubelet(ctx context.Context, clientset kubernetes.Interface, namespace, podName string) (*StorageMetrics, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	nodeName := pod.Spec.NodeName
	path := fmt.Sprintf("/api/v1/nodes/%s/proxy/stats/summary", nodeName)

	data, err := clientset.CoreV1().RESTClient().Get().
		AbsPath(path).
		DoRaw(ctx)
	if err != nil {
		return nil, fmt.Errorf("kubelet proxy: %w", err)
	}

	var summary KubeletSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		return nil, err
	}

	for _, p := range summary.Pods {
		if p.PodRef.Name != podName || p.PodRef.Namespace != namespace {
			continue
		}

		var volumes []VolumeStats
		for _, v := range p.Volumes {
			if v.PVCRef == nil {
				continue
			}
			vol := VolumeStats{}
			if v.UsedBytes != nil {
				vol.UsedMB = int64(*v.UsedBytes) / 1024 / 1024
			}
			if v.CapacityBytes != nil {
				vol.CapacityMB = int64(*v.CapacityBytes) / 1024 / 1024
			}
			if v.AvailableBytes != nil {
				vol.AvailableMB = int64(*v.AvailableBytes) / 1024 / 1024
			}
			volumes = append(volumes, vol)
		}

		return &StorageMetrics{
			PodName: podName,
			Volumes: volumes,
		}, nil
	}
	return nil, fmt.Errorf("pod %s not found in kubelet summary", podName)
}

type KubeletSummary struct {
	Pods []struct {
		PodRef struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"podRef"`
		EphemeralStorage *struct {
			UsedBytes     *uint64 `json:"usedBytes"`
			CapacityBytes *uint64 `json:"capacityBytes"`
		} `json:"ephemeral-storage"`
		Volumes []struct {
			Name           string  `json:"name"`
			UsedBytes      *uint64 `json:"usedBytes"`
			CapacityBytes  *uint64 `json:"capacityBytes"`
			AvailableBytes *uint64 `json:"availableBytes"`
			PVCRef         *struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"pvcRef"`
		} `json:"volume"`
	} `json:"pods"`
}
