package specs

import (
	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/naming"
	"github.com/openGemini/openGemini-operator/pkg/opengemini"
	"github.com/openGemini/openGemini-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateClusterMaintainService(cluster *opengeminiv1.GeminiCluster) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetServiceMaintainName(),
			Namespace: cluster.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "maintain",
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(opengemini.MaintainPort),
					Port:       opengemini.MaintainPort,
				},
			},
			Selector: map[string]string{
				opengeminiv1.LabelCluster: cluster.Name,
			},
		},
	}
}

func CreateClusterReadWriteService(cluster *opengeminiv1.GeminiCluster) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetServiceReadWriteName(),
			Namespace: cluster.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "readwrite",
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(opengemini.HttpPort),
					Port:       opengemini.HttpPort,
				},
			},
			Selector: map[string]string{
				opengeminiv1.LabelCluster: cluster.Name,
			},
		},
	}
}

func CreateMetaHeadlessServices(cluster *opengeminiv1.GeminiCluster) []*corev1.Service {
	svcs := []*corev1.Service{}
	for i := 0; i < int(*cluster.Spec.Meta.Replicas); i++ {
		metaInstanceName := naming.GenerateMetaInstance(cluster, i).Name
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      metaInstanceName,
				Namespace: cluster.Namespace,
			},
			Spec: corev1.ServiceSpec{
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: corev1.ClusterIPNone,
				Selector: utils.MergeLabels(
					cluster.Spec.Metadata.GetLabelsOrNil(),
					map[string]string{
						opengeminiv1.LabelCluster:     cluster.Name,
						opengeminiv1.LabelInstanceSet: "meta",
						opengeminiv1.LabelInstance:    metaInstanceName,
					}),
			},
		}
		svcs = append(svcs, svc)
	}
	return svcs
}
