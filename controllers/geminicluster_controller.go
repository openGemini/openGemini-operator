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
	"strings"

	_ "github.com/influxdata/influxdb1-client"
	influxClient "github.com/influxdata/influxdb1-client/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/configfile"
	"github.com/openGemini/openGemini-operator/pkg/naming"
	"github.com/openGemini/openGemini-operator/pkg/resources"
	"github.com/openGemini/openGemini-operator/pkg/specs"
	"github.com/openGemini/openGemini-operator/pkg/utils"
)

const (
	ControllerManagerName = "openminicluster-controller"
)

type K8sController interface {
	SetControllerReference(owner, controlled metav1.Object, scheme *runtime.Scheme) error
}

// GeminiClusterReconciler reconciles a GeminiCluster object
type GeminiClusterReconciler struct {
	client.Client
	Owner  client.FieldOwner
	Scheme *runtime.Scheme

	Controller K8sController
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
func (r *GeminiClusterReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (result ctrl.Result, err error) {
	log := log.FromContext(ctx)

	cluster := &opengeminiv1.GeminiCluster{}
	if err = r.Get(ctx, req.NamespacedName, cluster); err != nil {
		log.Error(err, "unable to fetch GeminiCluster")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info(
		"Starting Reconcile Cluster object",
		"name",
		cluster.Name,
		"namespace",
		cluster.Namespace,
	)

	//TODO:Set Defaults

	// Keep a copy of cluster prior to any manipulations.
	before := cluster.DeepCopy()

	// Handle delete
	if !cluster.DeletionTimestamp.IsZero() {
		log.Info("GeminiCluster deleted", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, nil
	}

	defer func() {
		if !equality.Semantic.DeepEqual(before.Status, cluster.Status) {
			patchErr := r.Client.Status().Patch(ctx, cluster, client.MergeFrom(before), r.Owner)
			if patchErr != nil {
				log.Error(patchErr, "patching cluster status")
				if err == nil {
					err = patchErr
				}
				return
			}
			log.V(1).Info("patched cluster status")
		}
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
		return result, nil
	} else {
		meta.RemoveStatusCondition(&cluster.Status.Conditions, opengeminiv1.ClusterProgressing)
	}

	if err := r.reconcileClusterStatus(ctx, cluster); err != nil {
		return result, fmt.Errorf("reconcile Cluster status failed: %w", err)
	}

	// reconciler resource
	if err := r.reconcileClusterConfigMap(ctx, cluster); err != nil {
		return result, fmt.Errorf("reconcile Cluster configmap objects failed: %w", err)
	}

	if err := r.reconcileClusterServices(ctx, cluster); err != nil {
		return result, fmt.Errorf("reconcile Cluster service objects failed: %w", err)
	}

	if err := r.reconcileClusterInstances(ctx, cluster); err != nil {
		return result, fmt.Errorf("reconcile Cluster instances failed: %w", err)
	}

	if err := r.reconcileAdminUserSecret(ctx, cluster); err != nil {
		return result, fmt.Errorf("reconcile admin secret objects failed: %w", err)
	}
	if err := r.reconcileAdminAccount(ctx, cluster); err != nil {
		return result, fmt.Errorf("reconcile admin account failed: %w", err)
	}

	cluster.Status.ObservedGeneration = cluster.GetGeneration()
	log.V(1).Info("reconciled cluster")

	return result, nil
}

// +kubebuilder:rbac:groups="apps",resources="statefulsets",verbs={list}
// +kubebuilder:rbac:groups="apps",resources="deployments",verbs={list}

func (r *GeminiClusterReconciler) reconcileClusterStatus(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
) error {
	dbs := &appsv1.StatefulSetList{}
	runners := &appsv1.DeploymentList{}

	selector, err := metav1.LabelSelectorAsSelector(naming.ClusterInstances(cluster.Name))
	if err != nil {
		return fmt.Errorf("build label selector failed: %w", err)
	}
	if err = r.List(ctx, dbs,
		client.InNamespace(cluster.Namespace),
		client.MatchingLabelsSelector{Selector: selector},
	); err != nil {
		return fmt.Errorf("list instance statefulsets failed: %w", err)
	}
	if err = r.List(ctx, runners,
		client.InNamespace(cluster.Namespace),
		client.MatchingLabelsSelector{Selector: selector},
	); err != nil {
		return fmt.Errorf("list instance deployments failed: %w", err)
	}

	// Fill out status sorted by set name.
	cluster.Status.InstanceSets = cluster.Status.InstanceSets[:0]
	for _, name := range []string{naming.InstanceMeta, naming.InstanceStore, naming.InstanceSql} {
		status := opengeminiv1.InstanceSetStatus{Name: name}

		for _, instance := range dbs.Items {
			if instance.Labels[opengeminiv1.LabelInstanceSet] != name {
				continue
			}

			status.Replicas += instance.Status.Replicas
			status.ReadyReplicas += instance.Status.ReadyReplicas
			status.UpdatedReplicas += instance.Status.UpdatedReplicas
		}
		for _, instance := range runners.Items {
			if instance.Labels[opengeminiv1.LabelInstanceSet] != name {
				continue
			}

			status.Replicas += instance.Status.Replicas
			status.ReadyReplicas += instance.Status.ReadyReplicas
			status.UpdatedReplicas += instance.Status.UpdatedReplicas
		}

		cluster.Status.InstanceSets = append(cluster.Status.InstanceSets, status)
	}

	return err
}

// +kubebuilder:rbac:groups="",resources="configmaps",verbs={get,create,patch}

func (r *GeminiClusterReconciler) reconcileClusterConfigMap(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
) error {
	confdata, err := configfile.NewBaseConfiguration(cluster)
	if err != nil {
		return fmt.Errorf("cannot generate cluster configruation: %w", err)
	}

	updatedConfig, err := r.getUpdatedConfig(ctx, cluster, confdata)
	if err != nil {
		return err
	}
	if updatedConfig != "" {
		confdata = updatedConfig
	}

	clusterConfigMap := &corev1.ConfigMap{ObjectMeta: naming.ClusterConfigMap(cluster)}
	clusterConfigMap.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	clusterConfigMap.Data = map[string]string{
		naming.ConfigurationFile: confdata,
	}
	confHash := utils.CalcMd5Hash(confdata)

	cluster.SetInheritedMetadata(&clusterConfigMap.ObjectMeta)
	if err := r.setControllerReference(cluster, clusterConfigMap); err != nil {
		return fmt.Errorf("set controller reference failed: %w", err)
	}

	if err := r.apply(ctx, clusterConfigMap); err != nil {
		return err
	}
	cluster.Status.AppliedConfigHash = confHash

	return nil
}

// getUpdatedConfig returns the updated config items if set config map
func (r *GeminiClusterReconciler) getUpdatedConfig(ctx context.Context, cluster *opengeminiv1.GeminiCluster, tmplConf string) (string, error) {
	var customCM corev1.ConfigMap
	err := r.Get(
		ctx,
		client.ObjectKey{Namespace: cluster.Namespace, Name: cluster.Spec.CustomConfigMapName},
		&customCM)
	if errors.IsNotFound(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("fetch custom conf ConfigMap failed, err: %w", err)
	}

	var newConfigData string
	var ok bool
	if newConfigData, ok = customCM.Data[naming.ConfigurationFile]; !ok {
		return "", nil
	}

	return configfile.UpdateConfig(tmplConf, newConfigData)
}

// +kubebuilder:rbac:groups="",resources="services",verbs={get,create}

func (r *GeminiClusterReconciler) reconcileClusterServices(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
) error {
	readWriteService := specs.CreateClusterReadWriteService(cluster)
	cluster.SetInheritedMetadata(&readWriteService.ObjectMeta)
	if err := r.setControllerReference(cluster, readWriteService); err != nil {
		return fmt.Errorf("set controller reference failed: %w", err)
	}

	if err := resources.CreateIfNotFound(ctx, r.Client, readWriteService); client.IgnoreAlreadyExists(
		err,
	) != nil {
		return err
	}

	MaintainService := specs.CreateClusterMaintainService(cluster)
	cluster.SetInheritedMetadata(&MaintainService.ObjectMeta)
	if err := r.setControllerReference(cluster, MaintainService); err != nil {
		return fmt.Errorf("set controller reference failed: %w", err)
	}

	if err := resources.CreateIfNotFound(ctx, r.Client, MaintainService); client.IgnoreAlreadyExists(
		err,
	) != nil {
		return err
	}

	metaHeadlessSvcs := specs.CreateMetaHeadlessServices(cluster)
	for _, svc := range metaHeadlessSvcs {
		cluster.SetInheritedMetadata(&svc.ObjectMeta)
		if err := r.setControllerReference(cluster, svc); err != nil {
			return fmt.Errorf("set controller reference failed: %w", err)
		}

		if err := resources.CreateIfNotFound(ctx, r.Client, svc); client.IgnoreAlreadyExists(
			err,
		) != nil {
			return err
		}
	}

	return nil
}

// +kubebuilder:rbac:groups="",resources="secrets",verbs={get,create,delete}

func (r *GeminiClusterReconciler) reconcileAdminUserSecret(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
) (err error) {
	log := log.FromContext(ctx)
	if cluster.Status.AdminUserInitialized {
		log.Info("admin has been initialized, will skip reconcile secret")
		return nil
	}

	if cluster.GetEnableHttpAuth() {
		var (
			adminUsername string
			adminPassword string
		)

		if cluster.Spec.CustomAdminSecretName != "" {
			var customAdminSecret corev1.Secret
			err = r.Get(
				ctx,
				client.ObjectKey{Namespace: cluster.Namespace, Name: cluster.Spec.CustomAdminSecretName},
				&customAdminSecret)
			if err != nil {
				return fmt.Errorf("get custom admin secret failed: %w", err)
			}
			adminUsername = string(customAdminSecret.Data["username"])
			adminPassword = string(customAdminSecret.Data["password"])
			if adminUsername == "" || adminPassword == "" {
				return fmt.Errorf("custom admin secret is not valid")
			}
		} else {
			adminUsername = opengeminiv1.DefaultAdminUsername
			adminPassword = utils.GenPassword(16, 4, 4, 4)
		}

		adminUserSecret := specs.CreateSecret(
			cluster.GetAdminUserSecretName(),
			cluster.Namespace,
			adminUsername,
			adminPassword)
		cluster.SetInheritedMetadata(&adminUserSecret.ObjectMeta)
		if err := r.setControllerReference(cluster, adminUserSecret); err != nil {
			return fmt.Errorf("set controller reference failed: %w", err)
		}

		if err := resources.CreateIfNotFound(ctx, r.Client, adminUserSecret); client.IgnoreAlreadyExists(err) != nil {
			return err
		}
	}

	return nil
}

func (r *GeminiClusterReconciler) reconcileAdminAccount(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
) (err error) {
	log := log.FromContext(ctx)
	if cluster.Status.AdminUserInitialized {
		log.Info("admin has been initialized, will skip reconcile account")
		return nil
	}

	if cluster.GetEnableHttpAuth() && cluster.IsSqlReady() {
		var (
			adminUsername string
			adminPassword string
		)

		var adminUserSecret corev1.Secret
		err = r.Get(
			ctx,
			client.ObjectKey{Namespace: cluster.Namespace, Name: cluster.GetAdminUserSecretName()},
			&adminUserSecret)
		if err != nil {
			return fmt.Errorf("get admin user secret failed: %w", err)
		}
		adminUsername = string(adminUserSecret.Data["username"])
		adminPassword = string(adminUserSecret.Data["password"])
		if adminUsername == "" || adminPassword == "" {
			return fmt.Errorf("admin user secret data is not valid")
		}

		sqlHost := fmt.Sprintf("%s.%s", cluster.GetServiceReadWriteName(), cluster.Namespace)
		cli, err := influxClient.NewHTTPClient(influxClient.HTTPConfig{
			Addr:     fmt.Sprintf("http://%s:8086", sqlHost),
			Username: adminUsername,
			Password: adminPassword,
		})
		if err != nil {
			return fmt.Errorf("create opengemini client failed, err: %s", err.Error())
		}
		defer cli.Close()

		adminQuery := fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s' WITH ALL PRIVILEGES", adminUsername, adminPassword)
		response, err := cli.Query(influxClient.NewQuery(adminQuery, "", ""))
		if err != nil {
			return err
		}
		if response.Error() != nil {
			if !strings.Contains(response.Error().Error(), "already exists") && !strings.Contains(response.Error().Error(), "is existed") {
				return response.Error()
			}
		}
		cluster.Status.AdminUserInitialized = true
	}

	return nil
}

func (r *GeminiClusterReconciler) reconcileClusterInstances(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
) error {

	for i := 0; i < int(*cluster.Spec.Meta.Replicas); i++ {
		if err := r.reconcileMetaInstance(ctx, cluster, i); err != nil {
			return err
		}
	}
	for i := 0; i < int(*cluster.Spec.Store.Replicas); i++ {
		if err := r.reconcileStoreInstance(ctx, cluster, i); err != nil {
			return err
		}
	}
	if err := r.reconcileSqlInstance(ctx, cluster); err != nil {
		return err
	}
	return nil
}

func (r *GeminiClusterReconciler) setControllerReference(
	owner *opengeminiv1.GeminiCluster, controlled client.Object,
) error {
	return r.Controller.SetControllerReference(owner, controlled, r.Client.Scheme())
	//return controllerutil.SetControllerReference(owner, controlled, r.Client.Scheme())
}

// func (r *GeminiClusterReconciler) setOwnerReference(
// 	owner *opengeminiv1.GeminiCluster, controlled client.Object,
// ) error {
// 	return controllerutil.SetOwnerReference(owner, controlled, r.Client.Scheme())
// }

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
