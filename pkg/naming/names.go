package naming

import (
	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ClusterConfigMap(cluster *opengeminiv1.GeminiCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-config",
	}
}
