package services

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sClient) SetupRBAC(ctx context.Context, namespace, userID string) error {
	// 1. Buat ServiceAccount untuk user
	if err := k.createServiceAccount(ctx, namespace, userID); err != nil {
		return err
	}

	// 2. Buat Role — hanya bisa akses namespace sendiri
	if err := k.createRole(ctx, namespace); err != nil {
		return err
	}

	// 3. Bind Role ke ServiceAccount
	if err := k.createRoleBinding(ctx, namespace, userID); err != nil {
		return err
	}

	return nil
}

func (k *K8sClient) createServiceAccount(ctx context.Context, namespace, userID string) error {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("user-%s", userID),
			Namespace: namespace,
			Labels: map[string]string{
				"managed-by": "workspace-controller",
				"user-id":    userID,
			},
		},
	}

	_, err := k.Client.CoreV1().ServiceAccounts(namespace).Create(ctx, sa, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service account: %w", err)
	}

	return nil
}

func (k *K8sClient) createRole(ctx context.Context, namespace string) error {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workspace-user-role",
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				// Boleh lihat pods, logs
				APIGroups: []string{""},
				Resources: []string{"pods", "pods/log", "services"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				// Boleh akses exec ke pod (untuk terminal di vscode)
				APIGroups: []string{""},
				Resources: []string{"pods/exec"},
				Verbs:     []string{"create"},
			},
			{
				// Tidak boleh delete/create deployment sendiri
				APIGroups: []string{"apps"},
				Resources: []string{"deployments"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}

	_, err := k.Client.RbacV1().Roles(namespace).Create(ctx, role, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

func (k *K8sClient) createRoleBinding(ctx context.Context, namespace, userID string) error {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workspace-user-rolebinding",
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      fmt.Sprintf("user-%s", userID),
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "workspace-user-role",
		},
	}

	_, err := k.Client.RbacV1().RoleBindings(namespace).Create(ctx, rb, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create role binding: %w", err)
	}

	return nil
}

func (k *K8sClient) DeleteRBAC(ctx context.Context, userID string) error {
	namespace := fmt.Sprintf("workspace-%s", userID)

	// Hapus RoleBinding
	k.Client.RbacV1().RoleBindings(namespace).Delete(ctx, "workspace-user-rolebinding", metav1.DeleteOptions{})

	// Hapus Role
	k.Client.RbacV1().Roles(namespace).Delete(ctx, "workspace-user-role", metav1.DeleteOptions{})

	// Hapus ServiceAccount
	k.Client.CoreV1().ServiceAccounts(namespace).Delete(ctx, fmt.Sprintf("user-%s", userID), metav1.DeleteOptions{})

	return nil
}
