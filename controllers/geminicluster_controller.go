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
	"fmt"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/configfile"
	"github.com/openGemini/openGemini-operator/pkg/naming"
	"github.com/openGemini/openGemini-operator/pkg/resources"
	"github.com/openGemini/openGemini-operator/pkg/specs"
	"github.com/sethvargo/go-password/password"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ControllerManagerName = "openminicluster-controller"
)

// GeminiClusterReconciler reconciles a GeminiCluster object
type GeminiClusterReconciler struct {
	client.Client
	Owner  client.FieldOwner
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

	cluster := &opengeminiv1.GeminiCluster{}
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
			Type:               opengeminiv1.ClusterProgressing,
			Status:             metav1.ConditionFalse,
			Reason:             "Paused",
			Message:            "No spec changes will be applied and no other statuses will be updated.",
			ObservedGeneration: cluster.GetGeneration(),
		})
		return ctrl.Result{}, nil
	} else {
		meta.RemoveStatusCondition(&cluster.Status.Conditions, opengeminiv1.ClusterProgressing)
	}

	// reconciler resource
	ogConf, err := configfile.NewConfiguration()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot generate cluster configruation: %w", err)
	}

	if err := r.reconcileClusterConfigMap(ctx, cluster, ogConf); err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot create Cluster configmap objects: %w", err)
	}

	if err := r.reconcileClusterServices(ctx, cluster); err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot create Cluster service objects: %w", err)
	}

	if err := r.reconcileSuperuserSecret(ctx, cluster); err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot create Superuser secret objects: %w", err)
	}

	if err := r.reconcileClusterInstances(ctx, cluster); err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot create Cluster instances: %w", err)
	}

	return ctrl.Result{}, nil
}

// +kubebuilder:rbac:groups="",resources="configmaps",verbs={get,create}

func (r *GeminiClusterReconciler) reconcileClusterConfigMap(
	ctx context.Context, cluster *opengeminiv1.GeminiCluster, opengeminiConf string,
) error {
	clusterConfigMap := &corev1.ConfigMap{ObjectMeta: naming.ClusterConfigMap(cluster)}
	clusterConfigMap.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	cluster.SetInheritedMetadata(&clusterConfigMap.ObjectMeta)
	if err := r.setControllerReference(cluster, clusterConfigMap); err != nil {
		return fmt.Errorf("set controller reference failed: %w", err)
	}

	if err := resources.CreateIfNotFound(ctx, r.Client, clusterConfigMap); client.IgnoreAlreadyExists(err) != nil {
		return err
	}

	return nil
}

// +kubebuilder:rbac:groups="",resources="services",verbs={get,create}

func (r *GeminiClusterReconciler) reconcileClusterServices(ctx context.Context, cluster *opengeminiv1.GeminiCluster) error {
	readWriteService := specs.CreateClusterReadWriteService(*cluster)
	cluster.SetInheritedMetadata(&readWriteService.ObjectMeta)
	if err := r.setControllerReference(cluster, readWriteService); err != nil {
		return fmt.Errorf("set controller reference failed: %w", err)
	}

	if err := resources.CreateIfNotFound(ctx, r.Client, readWriteService); client.IgnoreAlreadyExists(err) != nil {
		return err
	}

	MaintainService := specs.CreateClusterMaintainService(*cluster)
	cluster.SetInheritedMetadata(&MaintainService.ObjectMeta)
	if err := r.setControllerReference(cluster, MaintainService); err != nil {
		return fmt.Errorf("set controller reference failed: %w", err)
	}

	if err := resources.CreateIfNotFound(ctx, r.Client, MaintainService); client.IgnoreAlreadyExists(err) != nil {
		return err
	}

	return nil
}

// +kubebuilder:rbac:groups="",resources="services",verbs={get,create,delete}

func (r *GeminiClusterReconciler) reconcileSuperuserSecret(ctx context.Context, cluster *opengeminiv1.GeminiCluster) error {
	if cluster.GetEnableSuperuserAccess() && cluster.Spec.SuperuserSecretName == "" {
		superuserPassword, err := password.Generate(64, 10, 0, false, true)
		if err != nil {
			return err
		}

		superuserSecret := specs.CreateSecret(
			cluster.GetSuperuserSecretName(),
			cluster.Namespace,
			"root",
			superuserPassword)
		cluster.SetInheritedMetadata(&superuserSecret.ObjectMeta)
		if err := r.setControllerReference(cluster, superuserSecret); err != nil {
			return fmt.Errorf("set controller reference failed: %w", err)
		}

		if err := resources.CreateIfNotFound(ctx, r.Client, superuserSecret); client.IgnoreAlreadyExists(err) != nil {
			return err
		}
	}

	if !cluster.GetEnableSuperuserAccess() {
		var secret corev1.Secret
		err := r.Get(
			ctx,
			client.ObjectKey{Namespace: cluster.Namespace, Name: cluster.GetSuperuserSecretName()},
			&secret)
		if err != nil {
			if apierrs.IsNotFound(err) || apierrs.IsForbidden(err) {
				return nil
			}
			return err
		}

		if _, owned := IsOwnedByCluster(&secret); owned {
			return r.Delete(ctx, &secret)
		}
	}

	return nil
}

func (r *GeminiClusterReconciler) reconcileClusterInstances(ctx context.Context, cluster *opengeminiv1.GeminiCluster) error {

	for i := 0; i < int(*cluster.Spec.Meta.Replicas); i++ {
		r.reconcileMetaInstance(ctx, cluster, i)
	}
	for i := 0; i < int(*cluster.Spec.Store.Replicas); i++ {
		r.reconcileStoreInstance(ctx, cluster, i)
	}
	for i := 0; i < int(*cluster.Spec.SQL.Replicas); i++ {
		r.reconcileSqlInstance(ctx, cluster, i)
	}
	return nil
}

func (r *GeminiClusterReconciler) setControllerReference(
	owner *opengeminiv1.GeminiCluster, controlled client.Object,
) error {
	return controllerutil.SetControllerReference(owner, controlled, r.Client.Scheme())
}

func (r *GeminiClusterReconciler) setOwnerReference(
	owner *opengeminiv1.GeminiCluster, controlled client.Object,
) error {
	return controllerutil.SetOwnerReference(owner, controlled, r.Client.Scheme())
}

func IsOwnedByCluster(obj client.Object) (string, bool) {
	owner := metav1.GetControllerOf(obj)
	if owner == nil {
		return "", false
	}

	if owner.APIVersion != opengeminiv1.GroupVersion.String() {
		return "", false
	}

	return owner.Name, true
}

// +kubebuilder:rbac:groups="",resources="configmaps",verbs={get,list,watch}
// +kubebuilder:rbac:groups="",resources="services",verbs={get,list,watch}
// +kubebuilder:rbac:groups="",resources="secrets",verbs={get,list,watch}
// +kubebuilder:rbac:groups="",resources="persistentvolumeclaims",verbs={get,list,watch}
// +kubebuilder:rbac:groups="apps",resources="deployments",verbs={get,list,watch}
// +kubebuilder:rbac:groups="apps",resources="statefulsets",verbs={get,list,watch}

// SetupWithManager sets up the controller with the Manager.
func (r *GeminiClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&opengeminiv1.GeminiCluster{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
