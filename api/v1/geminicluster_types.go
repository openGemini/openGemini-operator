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

package v1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openGemini/openGemini-operator/pkg/utils"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GeminiClusterSpec defines the desired state of GeminiCluster
type GeminiClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Version string `json:"version"`
	// +optianal
	// +kubebuilder:default:=false
	Paused *bool `json:"paused,omitempty"`
	// +optional
	Metadata   *Metadata      `json:"metadata,omitempty"`
	SQL        SQLSpec        `json:"sql"`
	Meta       MetaSpec       `json:"meta"`
	Store      StoreSpec      `json:"store"`
	Monitoring MonitoringSpec `json:"monitoring"`
	Affinity   AffinitySpec   `json:"affinity"`

	// +optional
	CustomAdminSecretName string `json:"customAdminSecretName,omitempty"`
	// +kubebuilder:default:=false
	EnableHttpAuth *bool `json:"enableHttpAuth,omitempty"`
	// +optional
	CustomConfigMapName string `json:"customConfigMapName,omitempty"`
}

type SQLSpec struct {
	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	Replicas  *int32                      `json:"replicas,omitempty"`
	Image     string                      `json:"image"`
	Port      string                      `json:"port"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// +optional
	Parameters SQLParamsSpec `json:"parameters"`
}

type SQLParamsSpec struct {
	WriteTimeout       string `json:"write-timeout"`
	MaxConnectionLimit int32  `json:"max-connection-limit"`
}

type MetaSpec struct {
	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	Replicas *int32 `json:"replicas,omitempty"`
	Image    string `json:"image"`
	// +kubebuilder:validation:Required
	DataVolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"dataVolumeClaimSpec"`
	Resources           corev1.ResourceRequirements      `json:"resources,omitempty"`
	// +optional
	Parameters MetaParamsSpec `json:"parameters"`
}

type MetaParamsSpec struct {
	RetentionAutocreate bool `json:"retention-autocreate"`
}

type StoreSpec struct {
	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	Replicas *int32 `json:"replicas,omitempty"`
	Image    string `json:"image"`
	// +kubebuilder:validation:Required
	DataVolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"dataVolumeClaimSpec"`
	Resources           corev1.ResourceRequirements      `json:"resources,omitempty"`
	// +optional
	Parameters StoreParamsSpec `json:"parameters"`
}

type StoreParamsSpec struct {
	WalEnabled           bool  `json:"wal-enabled"`
	WriteConcurrentLimit int32 `json:"write-concurrent-limit"`
}

type MonitoringSpec struct {
	Type string `json:"type"`
}

type AffinitySpec struct {
	// +kubebuilder:default:=false
	EnablePodAntiAffinity bool `json:"enablePodAntiAffinity"`
}

// Metadata contains metadata for custom resources
type Metadata struct {
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

func (meta *Metadata) GetLabelsOrNil() map[string]string {
	if meta == nil {
		return nil
	}
	return meta.Labels
}

func (meta *Metadata) GetAnnotationsOrNil() map[string]string {
	if meta == nil {
		return nil
	}
	return meta.Annotations
}

type InstanceSetStatus struct {
	Name string `json:"name"`

	// Total number of ready pods.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// Total number of pods.
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// Total number of pods that have the desired specification.
	// +optional
	UpdatedReplicas int32 `json:"updatedReplicas,omitempty"`
}

// GeminiClusterStatus defines the observed state of GeminiCluster
type GeminiClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// observedGeneration represents the .metadata.generation on which the status was based.
	// +optional
	// +kubebuilder:validation:Minimum=0
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Current state of instances.
	// +listType=map
	// +listMapKey=name
	// +optional
	InstanceSets []InstanceSetStatus `json:"instances,omitempty"`

	CustomStatus  string `json:"customStatus,omitempty"`
	StatusDetails string `json:"statusDetails,omitempty"`

	// conditions represent the observations of cluster's current state.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors={"urn:alm:descriptor:io.kubernetes.conditions"}
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// if admin user has initialized
	// +kubebuilder:default:=false
	AdminUserInitialized bool `json:"adminUserInitialized"`
	// md5 hash of applied config file content
	AppliedConfigHash string `json:"appliedConfigHash"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=ogi
//+kubebuilder:subresource:status

// GeminiCluster is the Schema for the geminiclusters API
type GeminiCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GeminiClusterSpec   `json:"spec,omitempty"`
	Status GeminiClusterStatus `json:"status,omitempty"`
}

func (cluster *GeminiCluster) GetServiceMaintainName() string {
	return fmt.Sprintf("%v%v", cluster.Name, ServiceMaintainSuffix)
}

func (cluster *GeminiCluster) GetServiceReadWriteName() string {
	return fmt.Sprintf("%v%v", cluster.Name, ServiceReadWriteSuffix)
}

func (cluster *GeminiCluster) GetEnableHttpAuth() bool {
	if cluster.Spec.EnableHttpAuth != nil {
		return *cluster.Spec.EnableHttpAuth
	}

	return false
}

func (cluster *GeminiCluster) GetAdminUserSecretName() string {
	return fmt.Sprintf("%v%v", cluster.Name, AdminUserSecretSuffix)
}

func (cluster *GeminiCluster) SetInheritedMetadata(obj *metav1.ObjectMeta) {
	obj.Annotations = utils.MergeLabels(cluster.Spec.Metadata.GetAnnotationsOrNil())
	obj.Labels = utils.MergeLabels(cluster.Spec.Metadata.GetLabelsOrNil(),
		map[string]string{
			LabelCluster: cluster.Name,
		})
}

func (cluster *GeminiCluster) IsSqlReady() bool {
	for _, status := range cluster.Status.InstanceSets {
		if status.Name == "sql" {
			if status.Replicas != 0 && status.ReadyReplicas == status.Replicas {
				return true
			}
		}
	}
	return false
}

//+kubebuilder:object:root=true

// GeminiClusterList contains a list of GeminiCluster
type GeminiClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GeminiCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GeminiCluster{}, &GeminiClusterList{})
}
