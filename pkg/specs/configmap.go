package specs

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/openGemini/openGemini-operator/pkg/naming"
)

func ConfigFileConfigmapProjection(configmapName string) *corev1.ConfigMapProjection {
	return &corev1.ConfigMapProjection{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: configmapName,
		},
		Items: []corev1.KeyToPath{
			{
				Key:  naming.ConfigurationFile,
				Path: naming.ConfigurationFile,
			},
			{
				Key:  naming.EntrypointFile,
				Path: naming.EntrypointFile,
				Mode: &[]int32{0o700}[0],
			},
		},
	}
}
