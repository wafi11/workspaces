package services

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sClient) CreateResourceQuota(ctx context.Context, userId string, quota QuotaConfig) error {

	rq := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workspace-quota",
			Namespace: generateNamespace(userId),
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceLimitsCPU:       resource.MustParse(quota.CPULimit),
				corev1.ResourceLimitsMemory:    resource.MustParse(quota.MemoryLimit),
				corev1.ResourceRequestsCPU:     resource.MustParse(quota.CPURequest),
				corev1.ResourceRequestsMemory:  resource.MustParse(quota.MemoryRequest),
				corev1.ResourceRequestsStorage: resource.MustParse(quota.StorageLimit),
			},
		},
	}

	_, err := k.Client.CoreV1().ResourceQuotas(generateNamespace(userId)).Create(ctx, rq, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create resource quota: %w", err)
	}

	return nil
}

func (k *K8sClient) UpdateResourceQuota(ctx context.Context, namespace string, quota QuotaConfig) error {

	rq, err := k.Client.CoreV1().ResourceQuotas(namespace).Get(ctx, "workspace-quota", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get resource quota: %w", err)
	}

	rq.Spec.Hard = corev1.ResourceList{
		corev1.ResourceLimitsCPU:       resource.MustParse(quota.CPULimit),
		corev1.ResourceLimitsMemory:    resource.MustParse(quota.MemoryLimit),
		corev1.ResourceRequestsCPU:     resource.MustParse(quota.CPURequest),
		corev1.ResourceRequestsMemory:  resource.MustParse(quota.MemoryRequest),
		corev1.ResourceRequestsStorage: resource.MustParse(quota.StorageLimit),
	}

	_, err = k.Client.CoreV1().ResourceQuotas(namespace).Update(ctx, rq, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update resource quota: %w", err)
	}

	return nil
}
