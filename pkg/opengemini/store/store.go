package store

import (
	"context"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/naming"
	"github.com/openGemini/openGemini-operator/pkg/specs"
	corev1 "k8s.io/api/core/v1"
)

func DataVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      naming.DataVolume,
		MountPath: naming.DataMountPath,
	}
}

func ConfigVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      naming.ConfigVolume,
		MountPath: naming.ConfigMountPath,
		ReadOnly:  true,
	}
}

func InstancePod(
	ctx context.Context,
	inCluster *opengeminiv1.GeminiCluster,
	inDataVolumeName string,
	outInstancePod *corev1.PodSpec,
) {
	configVolumeMount := ConfigVolumeMount()
	configVolume := corev1.Volume{
		Name: configVolumeMount.Name,
		VolumeSource: corev1.VolumeSource{
			Projected: &corev1.ProjectedVolumeSource{
				DefaultMode: &[]int32{0o600}[0],
				Sources: []corev1.VolumeProjection{
					{ConfigMap: specs.ConfigFileConfigmapProjection(naming.ClusterConfigMap(inCluster).Name)},
				},
			},
		},
	}

	dataVolumeMount := DataVolumeMount()
	dataVolume := corev1.Volume{
		Name: dataVolumeMount.Name,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: inDataVolumeName,
				ReadOnly:  false,
			},
		},
	}

	container := corev1.Container{
		Name: naming.ContainerStore,

		Image:     inCluster.Spec.Store.Image,
		Resources: inCluster.Spec.Store.Resources,

		Ports: []corev1.ContainerPort{{
			Name:          naming.PortStore,
			ContainerPort: 666,
			Protocol:      corev1.ProtocolTCP,
		}},

		SecurityContext: specs.RestrictedSecurityContext(),
		VolumeMounts: []corev1.VolumeMount{
			configVolumeMount,
			dataVolumeMount,
		},
	}

	outInstancePod.Volumes = []corev1.Volume{
		configVolume,
		dataVolume,
	}

	outInstancePod.Containers = []corev1.Container{container}
}
