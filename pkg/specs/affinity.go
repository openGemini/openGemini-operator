package specs

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultAffinityWeight = 100

func CreateAffinity(enable bool, matchLabels map[string]string) *corev1.Affinity {
	affinity := corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{},
	}

	affinityTerm := corev1.PodAffinityTerm{
		TopologyKey: corev1.LabelHostname,
		LabelSelector: &metav1.LabelSelector{
			MatchLabels: matchLabels,
		},
	}

	if enable {
		affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []corev1.PodAffinityTerm{
			affinityTerm,
		}
	} else {
		affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.WeightedPodAffinityTerm{
			{
				Weight:          defaultAffinityWeight,
				PodAffinityTerm: affinityTerm,
			},
		}
	}
	return &affinity
}
