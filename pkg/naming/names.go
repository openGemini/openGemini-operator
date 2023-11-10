package naming

import (
	"fmt"
	"path"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/utils"
)

const (
	DataVolume      = "data-volume"
	DataMountPath   = "/ogdata"
	ConfigVolume    = "config-volume"
	ConfigMountPath = "/etc/opengemini"

	InstanceMeta  = "meta"
	InstanceStore = "store"
	InstanceSql   = "sql"

	ContainerMeta  = "meta"
	ContainerStore = "store"
	ContainerSql   = "sql"

	PortMeta  = ""
	PortStore = ""

	ConfigurationFile = "opengemini.conf"
)

func ClusterConfigMap(cluster *opengeminiv1.GeminiCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-config",
	}
}

func ConfigFilePath() string {
	return path.Join(ConfigMountPath, ConfigurationFile)
}

func GenerateMetaInstance(cluster *opengeminiv1.GeminiCluster, index int) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-" + "meta" + "-" + fmt.Sprintf("%d", index),
	}
}

func GenerateMetaHeadlessSvc(cluster *opengeminiv1.GeminiCluster, index int) string {
	return fmt.Sprintf(
		"%s.%s.svc.%s",
		GenerateMetaInstance(cluster, index).Name,
		cluster.Namespace,
		utils.GetClusterDomain(),
	)
}

func GenerateStoreInstance(cluster *opengeminiv1.GeminiCluster, index int) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-" + "store" + "-" + fmt.Sprintf("%d", index),
	}
}

func GenerateSqlInstance(cluster *opengeminiv1.GeminiCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-" + "sql",
	}
}

func GeneratePVC(
	cluster *opengeminiv1.GeminiCluster,
	instance *appsv1.StatefulSet,
) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      DataVolume,
		Labels:    instance.Labels,
	}
}
