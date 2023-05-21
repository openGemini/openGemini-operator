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

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	opengeminioperatorv1 "github.com/openGemini/openGemini-operator/api/v1"
)

// GeminiClusterReconciler reconciles a GeminiCluster object
type GeminiClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=opengemini-operator.opengemini.org,resources=geminiclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=opengemini-operator.opengemini.org,resources=geminiclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=opengemini-operator.opengemini.org,resources=geminiclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GeminiCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *GeminiClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	cluster := &opengeminioperatorv1.GeminiCluster{}
	if err := r.Get(ctx, req.NamespacedName, cluster); err != nil {
		log.Error(err, "unable to fetch GeminiCluster")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//TODO:Set Defaults

	// Keep a copy of cluster prior to any manipulations.
	before := cluster.DeepCopy()

	// Handle delete
	if !cluster.DeletionTimestamp.IsZero() {
		log.Info("GeminiCluster deleted", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, nil
	}

	defer func() (ctrl.Result, error) {
		if !equality.Semantic.DeepEqual(before.Status, cluster.Status) {
			// NOTE(cbandy): Kubernetes prior to v1.16.10 and v1.17.6 does not track
			// managed fields on the status subresource: https://issue.k8s.io/88901
			err := r.Client.Status().Patch(ctx, cluster, client.MergeFrom(before))
			if err != nil {
				log.Error(err, "patching cluster status")
				return ctrl.Result{}, err
			}
			log.V(1).Info("patched cluster status")
		}
		return ctrl.Result{}, nil
	}()

	// handle paused
	if cluster.Spec.Paused != nil && *cluster.Spec.Paused {
		meta.SetStatusCondition(&cluster.Status.Conditions, metav1.Condition{
			Type:               opengeminioperatorv1.ClusterProgressing,
			Status:             metav1.ConditionFalse,
			Reason:             "Paused",
			Message:            "No spec changes will be applied and no other statuses will be updated.",
			ObservedGeneration: cluster.GetGeneration(),
		})
		return ctrl.Result{}, nil
	} else {
		meta.RemoveStatusCondition(&cluster.Status.Conditions, opengeminioperatorv1.ClusterProgressing)
	}

	// reconciler resource

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GeminiClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&opengeminioperatorv1.GeminiCluster{}).
		Complete(r)
}
