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

	"github.com/openGemini/openGemini-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	SuperuserSecretName string `json:"superuserSecret,omitempty"`
	// +kubebuilder:default:=false
	EnableSuperuserAccess *bool `json:"enableSuperuserAccess,omitempty"`
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

// GeminiClusterStatus defines the observed state of GeminiCluster
type GeminiClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// observedGeneration is the most recent generation observed for this StatefulSet. It corresponds to the
	// StatefulSet's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// replicas is the number of Pods created by the StatefulSet controller.
	Replicas int32 `json:"replicas" protobuf:"varint,2,opt,name=replicas"`

	// readyReplicas is the number of Pods created by the StatefulSet controller that have a Ready Condition.
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,3,opt,name=readyReplicas"`

	// currentReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version
	// indicated by currentRevision.
	CurrentReplicas int32 `json:"currentReplicas,omitempty" protobuf:"varint,4,opt,name=currentReplicas"`

	// updatedReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version
	// indicated by updateRevision.
	UpdatedReplicas int32 `json:"updatedReplicas,omitempty" protobuf:"varint,5,opt,name=updatedReplicas"`

	// currentRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the
	// sequence [0,currentReplicas).
	CurrentRevision string `json:"currentRevision,omitempty" protobuf:"bytes,6,opt,name=currentRevision"`

	// updateRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence
	// [replicas-updatedReplicas,replicas)
	UpdateRevision string `json:"updateRevision,omitempty" protobuf:"bytes,7,opt,name=updateRevision"`

	// collisionCount is the count of hash collisions for the StatefulSet. The StatefulSet controller
	// uses this field as a collision avoidance mechanism when it needs to create the name for the
	// newest ControllerRevision.
	// +optional
	CollisionCount *int32 `json:"collisionCount,omitempty" protobuf:"varint,9,opt,name=collisionCount"`

	CustomStatus  string `json:"customStatus,omitempty"`
	StatusDetails string `json:"statusDetails,omitempty"`

	// Conditions for cluster object
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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

func (cluster *GeminiCluster) GetEnableSuperuserAccess() bool {
	if cluster.Spec.EnableSuperuserAccess != nil {
		return *cluster.Spec.EnableSuperuserAccess
	}

	return true
}

func (cluster *GeminiCluster) GetSuperuserSecretName() string {
	if cluster.Spec.SuperuserSecretName != "" {
		return cluster.Spec.SuperuserSecretName
	}

	return fmt.Sprintf("%v%v", cluster.Name, SuperUserSecretSuffix)
}

func (cluster *GeminiCluster) SetInheritedMetadata(obj *metav1.ObjectMeta) {
	obj.Annotations = utils.MergeLabels(cluster.Spec.Metadata.GetAnnotationsOrNil())
	obj.Labels = utils.MergeLabels(cluster.Spec.Metadata.GetLabelsOrNil(),
		map[string]string{
			LabelCluster: cluster.Name,
		})
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
