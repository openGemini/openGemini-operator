package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/naming"
	"github.com/openGemini/openGemini-operator/pkg/opengemini/meta"
	"github.com/openGemini/openGemini-operator/pkg/opengemini/sql"
	"github.com/openGemini/openGemini-operator/pkg/opengemini/store"
	"github.com/openGemini/openGemini-operator/pkg/utils"
)

// +kubebuilder:rbac:groups="apps",resources="statefulsets",verbs={get,create,patch}
// +kubebuilder:rbac:groups="",resources="persistentvolumeclaims",verbs={get,create,patch}

func (r *GeminiClusterReconciler) reconcileMetaInstance(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
	index int,
) error {
	instance := &appsv1.StatefulSet{}
	instance.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("StatefulSet"))
	instance.ObjectMeta = naming.GenerateMetaInstance(cluster, index)
	if err := r.setControllerReference(cluster, instance); err != nil {
		return err
	}

	generateInstanceStatefulSetIntent(ctx, cluster, naming.InstanceMeta, instance)
	instance.Spec.ServiceName = instance.Name

	pvc := &corev1.PersistentVolumeClaim{}
	pvc.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim"))
	pvc.ObjectMeta = naming.GeneratePVC(cluster, instance)
	if err := r.setControllerReference(cluster, pvc); err != nil {
		return err
	}
	pvc.Spec = cluster.Spec.Meta.DataVolumeClaimSpec
	instance.Spec.VolumeClaimTemplates = append(instance.Spec.VolumeClaimTemplates, *pvc)

	meta.InstancePod(ctx, cluster, pvc.Name, instance.Name, &instance.Spec.Template.Spec)

	if err := r.apply(ctx, instance); err != nil {
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups="apps",resources="statefulsets",verbs={get,create,patch}
// +kubebuilder:rbac:groups="",resources="persistentvolumeclaims",verbs={get,create,patch}

func (r *GeminiClusterReconciler) reconcileStoreInstance(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
	index int,
) error {
	instance := &appsv1.StatefulSet{}
	instance.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("StatefulSet"))
	instance.ObjectMeta = naming.GenerateStoreInstance(cluster, index)
	if err := r.setControllerReference(cluster, instance); err != nil {
		return err
	}

	generateInstanceStatefulSetIntent(ctx, cluster, naming.InstanceStore, instance)

	pvc := &corev1.PersistentVolumeClaim{}
	pvc.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim"))
	pvc.ObjectMeta = naming.GeneratePVC(cluster, instance)
	if err := r.setControllerReference(cluster, pvc); err != nil {
		return err
	}
	pvc.Spec = cluster.Spec.Meta.DataVolumeClaimSpec
	instance.Spec.VolumeClaimTemplates = append(instance.Spec.VolumeClaimTemplates, *pvc)

	store.InstancePod(ctx, cluster, pvc.Name, &instance.Spec.Template.Spec)
	if err := r.apply(ctx, instance); err != nil {
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups="apps",resources="deployments",verbs={get,create,patch}

func (r *GeminiClusterReconciler) reconcileSqlInstance(
	ctx context.Context,
	cluster *opengeminiv1.GeminiCluster,
) error {
	instance := &appsv1.Deployment{}
	instance.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Deployment"))
	instance.ObjectMeta = naming.GenerateSqlInstance(cluster)
	if err := r.setControllerReference(cluster, instance); err != nil {
		return err
	}

	generateInstanceDeploymentIntent(ctx, cluster, naming.InstanceSql, instance)

	sql.InstancePod(ctx, cluster, &instance.Spec.Template.Spec)
	if err := r.apply(ctx, instance); err != nil {
		return err
	}
	return nil
}

func generateInstanceStatefulSetIntent(
	_ context.Context,
	cluster *opengeminiv1.GeminiCluster,
	setName string,
	sts *appsv1.StatefulSet,
) {
	sts.Annotations = utils.MergeLabels(
		cluster.Spec.Metadata.GetAnnotationsOrNil())
	sts.Labels = utils.MergeLabels(
		cluster.Spec.Metadata.GetLabelsOrNil(),
		map[string]string{
			opengeminiv1.LabelCluster:     cluster.Name,
			opengeminiv1.LabelInstanceSet: setName,
			opengeminiv1.LabelInstance:    sts.Name,
		})
	sts.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			opengeminiv1.LabelCluster:     cluster.Name,
			opengeminiv1.LabelInstanceSet: setName,
			opengeminiv1.LabelInstance:    sts.Name,
		},
	}
	sts.Spec.Template.Annotations = utils.MergeLabels(
		cluster.Spec.Metadata.GetAnnotationsOrNil(),
	)
	sts.Spec.Template.Labels = utils.MergeLabels(
		cluster.Spec.Metadata.GetLabelsOrNil(),
		map[string]string{
			opengeminiv1.LabelCluster:     cluster.Name,
			opengeminiv1.LabelInstanceSet: setName,
			opengeminiv1.LabelInstance:    sts.Name,
			opengeminiv1.LabelConfigHash:  cluster.Status.AppliedConfigHash,
		})
	sts.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyAlways
	sts.Spec.Template.Spec.ShareProcessNamespace = &[]bool{true}[0]
	sts.Spec.Template.Spec.EnableServiceLinks = &[]bool{false}[0]

	sts.Spec.RevisionHistoryLimit = &[]int32{0}[0]
	sts.Spec.Replicas = &[]int32{1}[0]
}

func generateInstanceDeploymentIntent(
	_ context.Context,
	cluster *opengeminiv1.GeminiCluster,
	setName string,
	deploy *appsv1.Deployment,
) {
	deploy.Annotations = utils.MergeLabels(
		cluster.Spec.Metadata.GetAnnotationsOrNil())
	deploy.Labels = utils.MergeLabels(
		cluster.Spec.Metadata.GetLabelsOrNil(),
		map[string]string{
			opengeminiv1.LabelCluster:     cluster.Name,
			opengeminiv1.LabelInstanceSet: setName,
			opengeminiv1.LabelInstance:    deploy.Name,
		})
	deploy.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			opengeminiv1.LabelCluster:     cluster.Name,
			opengeminiv1.LabelInstanceSet: setName,
			opengeminiv1.LabelInstance:    deploy.Name,
		},
	}
	deploy.Spec.Template.Annotations = utils.MergeLabels(
		cluster.Spec.Metadata.GetAnnotationsOrNil(),
	)
	deploy.Spec.Template.Labels = utils.MergeLabels(
		cluster.Spec.Metadata.GetLabelsOrNil(),
		map[string]string{
			opengeminiv1.LabelCluster:     cluster.Name,
			opengeminiv1.LabelInstanceSet: setName,
			opengeminiv1.LabelInstance:    deploy.Name,
			opengeminiv1.LabelConfigHash:  cluster.Status.AppliedConfigHash,
		})
	deploy.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyAlways
	deploy.Spec.Template.Spec.ShareProcessNamespace = &[]bool{true}[0]
	deploy.Spec.Template.Spec.EnableServiceLinks = &[]bool{false}[0]

	deploy.Spec.RevisionHistoryLimit = &[]int32{0}[0]
	deploy.Spec.Replicas = cluster.Spec.SQL.Replicas
}
