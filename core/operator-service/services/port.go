package services

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (k *K8sClient) CreatePort(ctx context.Context, userId,workspaceName string, port int) error {
	namespace :=  fmt.Sprintf("ws-%s",userId)
	serviceName := fmt.Sprintf("%s-%d-svc", workspaceName, port)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
			Labels: map[string]string{
				"app": workspaceName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": workspaceName, 
			},
			Ports: []corev1.ServicePort{
				{
					Name:       fmt.Sprintf("port-%d", port),
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// Buat Service di K8s
	_, err := k.Client.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Failed to make services : %s",err.Error())
	}

	return nil
}
func (k *K8sClient) ExposeToIngress(ctx context.Context, useerId, workspace_name, serviceName, domain string, port int32) error {
    ns := generateNamespace(useerId)

    ingress, err := k.Client.NetworkingV1().Ingresses(ns).Get(ctx, workspace_name, metav1.GetOptions{})
    if err != nil {
        return fmt.Errorf("Ingress %s tidak ditemukan: %v", workspace_name, err)
    }

    newRule := networkingv1.IngressRule{
        Host: domain,
        IngressRuleValue: networkingv1.IngressRuleValue{
            HTTP: &networkingv1.HTTPIngressRuleValue{
                Paths: []networkingv1.HTTPIngressPath{
                    {
                        Path:     "/",
                        PathType: func() *networkingv1.PathType { pt := networkingv1.PathTypePrefix; return &pt }(),
                        Backend: networkingv1.IngressBackend{
                            Service: &networkingv1.IngressServiceBackend{
                                Name: serviceName,
                                Port: networkingv1.ServiceBackendPort{ Number: port },
                            },
                        },
                    },
                },
            },
        },
    }

    exists := false
    for _, rule := range ingress.Spec.Rules {
        if rule.Host == domain {
            exists = true
            break
        }
    }

    if !exists {
        ingress.Spec.Rules = append(ingress.Spec.Rules, newRule)
        _, err = k.Client.NetworkingV1().Ingresses(ns).Update(ctx, ingress, metav1.UpdateOptions{})
        if err != nil {
            return fmt.Errorf("failed to update ingress : %s",err.Error())
        }
    }

    return nil
}

func (k *K8sClient) DeletePort(ctx context.Context, userId, workspaceName string, port int) error {
	namespace := fmt.Sprintf("ws-%s", userId)
	serviceName := fmt.Sprintf("%s-%d-svc", workspaceName, port)

	err := k.Client.CoreV1().Services(namespace).Delete(ctx, serviceName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete service: %s", err.Error())
	}

	return nil
}

func (k *K8sClient) RemoveFromIngress(ctx context.Context, userId, workspaceName, domain string) error {
	ns := generateNamespace(userId)

	ingress, err := k.Client.NetworkingV1().Ingresses(ns).Get(ctx, workspaceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("ingress %s tidak ditemukan: %v", workspaceName, err)
	}

	filtered := make([]networkingv1.IngressRule, 0, len(ingress.Spec.Rules))
	for _, rule := range ingress.Spec.Rules {
		if rule.Host != domain {
			filtered = append(filtered, rule)
		}
	}

	// Ga ada yang dihapus, skip update
	if len(filtered) == len(ingress.Spec.Rules) {
		return nil
	}

	ingress.Spec.Rules = filtered
	_, err = k.Client.NetworkingV1().Ingresses(ns).Update(ctx, ingress, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ingress: %s", err.Error())
	}

	return nil
}