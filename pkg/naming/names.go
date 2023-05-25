package naming

import (
	"fmt"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DataVolume      = "data-volume"
	DataMountPath   = "/ogdata"
	ConfigVolume    = "config-volume"
	ConfigMountPath = "/etc/opengemini"

	ContainerMeta  = "meta"
	ContainerStore = "store"

	PortMeta  = ""
	PortStore = ""

	ConfigurationFile     = "opengemini.conf"
	ConfigurationFilePath = "/etc/opengemini"
)

func ClusterConfigMap(cluster *opengeminiv1.GeminiCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-config",
	}
}

func GenerateMetaInstance(cluster *opengeminiv1.GeminiCluster, index int) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-" + "meta" + "-" + fmt.Sprintf("%2d", index),
	}
}

func GenerateStoreInstance(cluster *opengeminiv1.GeminiCluster, index int) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-" + "data" + "-" + fmt.Sprintf("%2d", index),
	}
}

func GenerateSqlInstance(cluster *opengeminiv1.GeminiCluster, index int) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-" + "sql" + "-" + fmt.Sprintf("%2d", index),
	}
}

func GeneratePVC(cluster *opengeminiv1.GeminiCluster, instance *appsv1.StatefulSet) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      DataVolume,
		Labels:    instance.Labels,
	}
}
