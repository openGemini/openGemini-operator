package specs

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/naming"
	"github.com/openGemini/openGemini-operator/pkg/opengemini"
	"github.com/openGemini/openGemini-operator/pkg/utils"
)

func CreateClusterMaintainService(cluster *opengeminiv1.OpenGeminiCluster) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetServiceMaintainName(),
			Namespace: cluster.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       opengemini.MaintainPortName,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(opengemini.MaintainPort),
					Port:       opengemini.MaintainPort,
				},
			},
			Selector: map[string]string{
				opengeminiv1.LabelCluster:     cluster.Name,
				opengeminiv1.LabelInstanceSet: naming.InstanceMeta,
			},
		},
	}
}

func CreateClusterReadWriteService(cluster *opengeminiv1.OpenGeminiCluster) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetServiceReadWriteName(),
			Namespace: cluster.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       opengemini.HttpPortName,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(opengemini.HttpPort),
					Port:       opengemini.HttpPort,
				},
			},
			Selector: map[string]string{
				opengeminiv1.LabelCluster:     cluster.Name,
				opengeminiv1.LabelInstanceSet: naming.InstanceSql,
			},
		},
	}
}

func instanceHeadlessService(cluster *opengeminiv1.OpenGeminiCluster, instanceName, instanceSet string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: cluster.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
			Selector: utils.MergeLabels(
				cluster.Spec.Metadata.GetLabelsOrNil(),
				map[string]string{
					opengeminiv1.LabelCluster:     cluster.Name,
					opengeminiv1.LabelInstanceSet: instanceSet,
					opengeminiv1.LabelInstance:    instanceName,
				}),
		},
	}
}

func CreateInstanceHeadlessServices(cluster *opengeminiv1.OpenGeminiCluster) []*corev1.Service {
	svcs := []*corev1.Service{}
	for i := 0; i < int(*cluster.Spec.Meta.Replicas); i++ {
		metaInstanceName := naming.GenerateMetaInstance(cluster, i).Name
		svc := instanceHeadlessService(cluster, metaInstanceName, naming.InstanceMeta)
		svcs = append(svcs, svc)
	}
	for i := 0; i < int(*cluster.Spec.Store.Replicas); i++ {
		storeInstanceName := naming.GenerateStoreInstance(cluster, i).Name
		svc := instanceHeadlessService(cluster, storeInstanceName, naming.InstanceStore)
		svcs = append(svcs, svc)
	}
	for i := 0; i < int(*cluster.Spec.SQL.Replicas); i++ {
		sqlInstanceName := naming.GenerateSqlInstance(cluster, i).Name
		svc := instanceHeadlessService(cluster, sqlInstanceName, naming.InstanceSql)
		svcs = append(svcs, svc)
	}
	return svcs
}
