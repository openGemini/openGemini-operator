package sql

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/naming"
	"github.com/openGemini/openGemini-operator/pkg/specs"
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
	inCluster *opengeminiv1.OpenGeminiCluster,
	inInstanceName string,
	outInstancePod *corev1.PodSpec,
) {
	configVolumeMount := ConfigVolumeMount()
	configVolume := corev1.Volume{
		Name: configVolumeMount.Name,
		VolumeSource: corev1.VolumeSource{
			Projected: &corev1.ProjectedVolumeSource{
				DefaultMode: &[]int32{0o600}[0],
				Sources: []corev1.VolumeProjection{
					{
						ConfigMap: specs.ConfigFileConfigmapProjection(
							naming.ClusterConfigMap(inCluster).Name,
						),
					},
				},
			},
		},
	}

	dataVolumeMount := DataVolumeMount()
	dataVolume := corev1.Volume{
		Name: dataVolumeMount.Name,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}

	container := corev1.Container{
		Name:      naming.ContainerSql,
		Command:   []string{naming.EntrypointFilePath()},
		Image:     inCluster.Spec.SQL.Image,
		Resources: inCluster.Spec.SQL.Resources,
		Env: []corev1.EnvVar{
			{
				Name: "HOST_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
			{
				Name:  "DATA_DIR",
				Value: naming.DataMountPath,
			},
			{
				Name:  "CONFIG_PATH",
				Value: naming.ConfigFilePath(),
			},
			{
				Name:  "APP",
				Value: "sql",
			},
			{
				Name:  "SQL_DOMAIN",
				Value: naming.GenerateInstanceFQDN(inInstanceName, inCluster.Namespace),
			},
		},

		Ports: []corev1.ContainerPort{{
			Name:          naming.PortMeta,
			ContainerPort: 666,
			Protocol:      corev1.ProtocolTCP,
		}},

		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/ping",
					Port: intstr.FromInt(8086),
				},
			},
		},

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
