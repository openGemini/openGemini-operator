package specs

import (
	corev1 "k8s.io/api/core/v1"
)

func PodSecurityContext() *corev1.PodSecurityContext {
	onRootMismatch := corev1.FSGroupChangeOnRootMismatch
	return &corev1.PodSecurityContext{
		FSGroupChangePolicy: &onRootMismatch,
	}
}

func RestrictedSecurityContext() *corev1.SecurityContext {
	return &corev1.SecurityContext{
		AllowPrivilegeEscalation: &[]bool{false}[0],
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
		Privileged:             &[]bool{false}[0],
		ReadOnlyRootFilesystem: &[]bool{true}[0],
		RunAsNonRoot:           &[]bool{true}[0],
	}
}
