package naming

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
)

func ClusterInstances(cluster string) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{
			opengeminiv1.LabelCluster: cluster,
		},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{Key: opengeminiv1.LabelInstance, Operator: metav1.LabelSelectorOpExists},
		},
	}
}
