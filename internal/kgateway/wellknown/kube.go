package wellknown

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

var (
	SecretGVK         = corev1.SchemeGroupVersion.WithKind("Secret")
	ConfigMapGVK      = corev1.SchemeGroupVersion.WithKind("ConfigMap")
	ServiceGVK        = corev1.SchemeGroupVersion.WithKind("Service")
	ServiceAccountGVK = corev1.SchemeGroupVersion.WithKind("ServiceAccount")
	PodGVK            = corev1.SchemeGroupVersion.WithKind("Pod")

	// RBAC GVKs
	ClusterRoleBindingGVK = rbacv1.SchemeGroupVersion.WithKind("ClusterRoleBinding")

	DeploymentGVK = appsv1.SchemeGroupVersion.WithKind("Deployment")
)
