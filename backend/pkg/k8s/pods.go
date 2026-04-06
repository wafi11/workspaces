package k8s

// func (k *K8sClient) CreateDeployment(ctx context.Context, job WorkspaceJob) error {
// 	var envVars []corev1.EnvVar
// 	for key, val := range job.EnvVars {
// 		envVars = append(envVars, corev1.EnvVar{
// 			Name:  key,
// 			Value: fmt.Sprintf("%v", val),
// 		})
// 	}

// 	// 1. buat PVC dulu
// 	if err := k.createPVC(ctx, job); err != nil {
// 		return fmt.Errorf("failed to create pvc: %w", err)
// 	}

// 	// 2. deployment
// 	deployment := &appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      job.WorkspaceId,
// 			Namespace: job.Namespace,
// 			Labels: map[string]string{
// 				"app":          job.WorkspaceId,
// 				"managed-by":   "workspace-operator",
// 				"workspace-id": job.WorkspaceId,
// 			},
// 		},
// 		Spec: appsv1.DeploymentSpec{
// 			Replicas: int32Ptr(1),
// 			Selector: &metav1.LabelSelector{
// 				MatchLabels: map[string]string{"app": job.WorkspaceId},
// 			},
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: map[string]string{
// 						"app":          job.WorkspaceId,
// 						"workspace-id": job.WorkspaceId,
// 					},
// 				},
// 				Spec: corev1.PodSpec{
// 					Containers: []corev1.Container{
// 						{
// 							Name:    "terminal",
// 							Image:   "tsl0922/ttyd:latest",
// 							Command: []string{"ttyd"},
// 							Args:    []string{"-W", "-p", "7681", "/bin/sh"},
// 							Env:     envVars,
// 							Ports: []corev1.ContainerPort{
// 								{ContainerPort: 7681},
// 							},
// 							Resources: corev1.ResourceRequirements{
// 								Requests: corev1.ResourceList{
// 									corev1.ResourceCPU:    resource.MustParse("100m"),
// 									corev1.ResourceMemory: resource.MustParse("128Mi"),
// 								},
// 								Limits: corev1.ResourceList{
// 									corev1.ResourceCPU:    resource.MustParse("500m"),
// 									corev1.ResourceMemory: resource.MustParse("256Mi"),
// 								},
// 							},
// 							VolumeMounts: []corev1.VolumeMount{
// 								{
// 									Name:      "workspace-storage",
// 									MountPath: "/workspace",
// 								},
// 							},
// 						},
// 					},
// 					Volumes: []corev1.Volume{
// 						{
// 							Name: "workspace-storage",
// 							VolumeSource: corev1.VolumeSource{
// 								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
// 									ClaimName: fmt.Sprintf("pvc-%s", job.WorkspaceId),
// 								},
// 							},
// 						},
// 					},
// 					AutomountServiceAccountToken: boolPtr(false),
// 				},
// 			},
// 		},
// 	}

// 	_, err := k.Client.AppsV1().Deployments(job.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
// 	if err != nil {
// 		return fmt.Errorf("failed to create deployment: %w", err)
// 	}

// 	// 3. service
// 	if err := k.createService(ctx, job); err != nil {
// 		return fmt.Errorf("failed to create service: %w", err)
// 	}

// 	// 4. ingress
// 	if err := k.createIngress(ctx, job); err != nil {
// 		return fmt.Errorf("failed to create ingress: %w", err)
// 	}

// 	return nil
// }

// // ─── PVC (Longhorn) ───────────────────────────────────────────────────────────

// func (k *K8sClient) createPVC(ctx context.Context, job WorkspaceJob) error {
// 	storageClass := "longhorn"
// 	pvc := &corev1.PersistentVolumeClaim{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      fmt.Sprintf("pvc-%s", job.WorkspaceId),
// 			Namespace: job.Namespace,
// 			Annotations: map[string]string{
// 				"longhorn.io/numberOfReplicas": "2",
// 			},
// 			Labels: map[string]string{
// 				"managed-by":   "workspace-operator",
// 				"workspace-id": job.WorkspaceId,
// 			},
// 		},
// 		Spec: corev1.PersistentVolumeClaimSpec{
// 			AccessModes: []corev1.PersistentVolumeAccessMode{
// 				corev1.ReadWriteOnce,
// 			},
// 			StorageClassName: &storageClass,
// 			Resources: corev1.VolumeResourceRequirements{
// 				Requests: corev1.ResourceList{
// 					corev1.ResourceStorage: resource.MustParse("2Gi"),
// 				},
// 			},
// 		},
// 	}

// 	_, err := k.Client.CoreV1().PersistentVolumeClaims(job.Namespace).Create(ctx, pvc, metav1.CreateOptions{})
// 	if err != nil {
// 		return fmt.Errorf("failed to create pvc: %w", err)
// 	}

// 	return nil
// }

// // ─── Service ──────────────────────────────────────────────────────────────────

// func (k *K8sClient) createService(ctx context.Context, job workspaceservice.WorkspaceJob) error {
// 	svc := &corev1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      fmt.Sprintf("svc-%s", job.WorkspaceId),
// 			Namespace: job.Namespace,
// 			Labels: map[string]string{
// 				"managed-by":   "workspace-operator",
// 				"workspace-id": job.WorkspaceId,
// 			},
// 		},
// 		Spec: corev1.ServiceSpec{
// 			Selector: map[string]string{
// 				"app": job.WorkspaceId,
// 			},
// 			Ports: []corev1.ServicePort{
// 				{
// 					Name:     "terminal",
// 					Port:     7681,
// 					Protocol: corev1.ProtocolTCP,
// 				},
// 			},
// 			Type: corev1.ServiceTypeClusterIP,
// 		},
// 	}

// 	_, err := k.Client.CoreV1().Services(job.Namespace).Create(ctx, svc, metav1.CreateOptions{})
// 	if err != nil {
// 		return fmt.Errorf("failed to create service: %w", err)
// 	}

// 	return nil
// }

// // ─── Ingress ──────────────────────────────────────────────────────────────────

// func (k *K8sClient) createIngress(ctx context.Context, job workspaceservice.WorkspaceJob) error {
// 	hostname := fmt.Sprintf("%s.workspace.local", job.WorkspaceId)
// 	pathType := networkingv1.PathTypePrefix
// 	ingressClass := "nginx"

// 	ingress := &networkingv1.Ingress{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      fmt.Sprintf("ing-%s", job.WorkspaceId),
// 			Namespace: job.Namespace,
// 			Labels: map[string]string{
// 				"managed-by":   "workspace-operator",
// 				"workspace-id": job.WorkspaceId,
// 			},
// 			Annotations: map[string]string{
// 				"nginx.ingress.kubernetes.io/proxy-read-timeout":      "3600",
// 				"nginx.ingress.kubernetes.io/proxy-send-timeout":      "3600",
// 				"nginx.ingress.kubernetes.io/proxy-http-version":      "1.1",
// 				"nginx.ingress.kubernetes.io/proxy-buffering":         "off",
// 				"nginx.ingress.kubernetes.io/proxy-request-buffering": "off",
// 				"nginx.ingress.kubernetes.io/configuration-snippet": `
// 				proxy_set_header Upgrade $http_upgrade;
// 				proxy_set_header Connection "Upgrade";
// 				proxy_set_header Host $host;
// 			`,
// 			},
// 		},
// 		Spec: networkingv1.IngressSpec{
// 			IngressClassName: &ingressClass,
// 			Rules: []networkingv1.IngressRule{
// 				{
// 					Host: hostname,
// 					IngressRuleValue: networkingv1.IngressRuleValue{
// 						HTTP: &networkingv1.HTTPIngressRuleValue{
// 							Paths: []networkingv1.HTTPIngressPath{
// 								{
// 									Path:     "/",
// 									PathType: &pathType,
// 									Backend: networkingv1.IngressBackend{
// 										Service: &networkingv1.IngressServiceBackend{
// 											Name: fmt.Sprintf("svc-%s", job.WorkspaceId),
// 											Port: networkingv1.ServiceBackendPort{
// 												Number: 7681,
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	_, err := k.Client.NetworkingV1().Ingresses(job.Namespace).Create(ctx, ingress, metav1.CreateOptions{})
// 	if err != nil {
// 		return fmt.Errorf("failed to create ingress: %w", err)
// 	}

// 	return nil
// }

// func int32Ptr(i int32) *int32 { return &i }
// func boolPtr(b bool) *bool    { return &b }
