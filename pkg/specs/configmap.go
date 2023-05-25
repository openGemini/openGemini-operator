package specs

import (
	"github.com/openGemini/openGemini-operator/pkg/naming"
	corev1 "k8s.io/api/core/v1"
)

func ConfigFileConfigmapProjection(configmapName string) *corev1.ConfigMapProjection {
	return &corev1.ConfigMapProjection{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: configmapName,
		},
		Items: []corev1.KeyToPath{
			{
				Key:  naming.ConfigurationFile,
				Path: naming.ConfigurationFilePath,
			},
		},
	}
}
