//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AffinitySpec) DeepCopyInto(out *AffinitySpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AffinitySpec.
func (in *AffinitySpec) DeepCopy() *AffinitySpec {
	if in == nil {
		return nil
	}
	out := new(AffinitySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeminiCluster) DeepCopyInto(out *GeminiCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeminiCluster.
func (in *GeminiCluster) DeepCopy() *GeminiCluster {
	if in == nil {
		return nil
	}
	out := new(GeminiCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GeminiCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeminiClusterList) DeepCopyInto(out *GeminiClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GeminiCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeminiClusterList.
func (in *GeminiClusterList) DeepCopy() *GeminiClusterList {
	if in == nil {
		return nil
	}
	out := new(GeminiClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GeminiClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeminiClusterSpec) DeepCopyInto(out *GeminiClusterSpec) {
	*out = *in
	if in.Paused != nil {
		in, out := &in.Paused, &out.Paused
		*out = new(bool)
		**out = **in
	}
	in.SQL.DeepCopyInto(&out.SQL)
	in.Meta.DeepCopyInto(&out.Meta)
	in.Store.DeepCopyInto(&out.Store)
	out.Monitoring = in.Monitoring
	out.UserSecret = in.UserSecret
	out.Affinity = in.Affinity
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeminiClusterSpec.
func (in *GeminiClusterSpec) DeepCopy() *GeminiClusterSpec {
	if in == nil {
		return nil
	}
	out := new(GeminiClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeminiClusterStatus) DeepCopyInto(out *GeminiClusterStatus) {
	*out = *in
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = new(int64)
		**out = **in
	}
	if in.CollisionCount != nil {
		in, out := &in.CollisionCount, &out.CollisionCount
		*out = new(int32)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeminiClusterStatus.
func (in *GeminiClusterStatus) DeepCopy() *GeminiClusterStatus {
	if in == nil {
		return nil
	}
	out := new(GeminiClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetaParamsSpec) DeepCopyInto(out *MetaParamsSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetaParamsSpec.
func (in *MetaParamsSpec) DeepCopy() *MetaParamsSpec {
	if in == nil {
		return nil
	}
	out := new(MetaParamsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetaSpec) DeepCopyInto(out *MetaSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.DataVolumeClaimSpec.DeepCopyInto(&out.DataVolumeClaimSpec)
	in.Resource.DeepCopyInto(&out.Resource)
	out.Parameters = in.Parameters
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetaSpec.
func (in *MetaSpec) DeepCopy() *MetaSpec {
	if in == nil {
		return nil
	}
	out := new(MetaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringSpec) DeepCopyInto(out *MonitoringSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringSpec.
func (in *MonitoringSpec) DeepCopy() *MonitoringSpec {
	if in == nil {
		return nil
	}
	out := new(MonitoringSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLParamsSpec) DeepCopyInto(out *SQLParamsSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLParamsSpec.
func (in *SQLParamsSpec) DeepCopy() *SQLParamsSpec {
	if in == nil {
		return nil
	}
	out := new(SQLParamsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLSpec) DeepCopyInto(out *SQLSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.Resource.DeepCopyInto(&out.Resource)
	out.Parameters = in.Parameters
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLSpec.
func (in *SQLSpec) DeepCopy() *SQLSpec {
	if in == nil {
		return nil
	}
	out := new(SQLSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StoreParamsSpec) DeepCopyInto(out *StoreParamsSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StoreParamsSpec.
func (in *StoreParamsSpec) DeepCopy() *StoreParamsSpec {
	if in == nil {
		return nil
	}
	out := new(StoreParamsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StoreSpec) DeepCopyInto(out *StoreSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.DataVolumeClaimSpec.DeepCopyInto(&out.DataVolumeClaimSpec)
	in.Resource.DeepCopyInto(&out.Resource)
	out.Parameters = in.Parameters
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StoreSpec.
func (in *StoreSpec) DeepCopy() *StoreSpec {
	if in == nil {
		return nil
	}
	out := new(StoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UserSecretSpec) DeepCopyInto(out *UserSecretSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UserSecretSpec.
func (in *UserSecretSpec) DeepCopy() *UserSecretSpec {
	if in == nil {
		return nil
	}
	out := new(UserSecretSpec)
	in.DeepCopyInto(out)
	return out
}
