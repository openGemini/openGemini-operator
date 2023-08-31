/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"testing"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MockK8sClient struct {
	client.Client
}

func (MockK8sClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	obj.(*corev1.ConfigMap).Data = map[string]string{
		"opengemini.conf": `
[common] 
meta-join = ["opengemini-meta-01.test.svc.cluster.local:8092", "opengemini-meta-02.test.svc.cluster.local:8092", "opengemini-meta-03.test.svc.cluster.local:8092"]
executor-memory-size-limit= "16G"  
executor-memory-wait-time = "0s"
cpu-num = 8  
memory-size = "32G"`,
	}

	return nil
}

func (MockK8sClient) Scheme() *runtime.Scheme {
	return nil
}

func (MockK8sClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}

type MockK8sController struct {
	K8sController
}

func (MockK8sController) SetControllerReference(owner, controlled metav1.Object, scheme *runtime.Scheme) error {
	return nil
}

func Test_ReconcileClusterConfigMap(t *testing.T) {
	var r = &GeminiClusterReconciler{
		Client:     &MockK8sClient{},
		Controller: &MockK8sController{},
	}

	var metaReplicas int32 = 3

	cluster := &opengeminiv1.GeminiCluster{
		Spec: opengeminiv1.GeminiClusterSpec{
			Meta: opengeminiv1.MetaSpec{
				Replicas: &metaReplicas,
			},
			CustomConfigMapName: "config.conf",
		},
	}
	cluster.Name = "test"
	cluster.Namespace = "opengemini"

	err := r.reconcileClusterConfigMap(context.Background(), cluster)
	assert.NoError(t, err)
	assert.Equal(t, cluster.Status.AppliedConfigHash, "1e94d9f395144eaf69571fb83161d278")
}
