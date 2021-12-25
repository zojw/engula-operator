/*
Copyright 2021.

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
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	clusterv1alpha1 "github.com/engula/engula-operator/api/v1alpha1"
)

// EngulaClusterReconciler reconciles a EngulaCluster object
type EngulaClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cluster.engula.io,resources=engulaclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.engula.io,resources=engulaclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.engula.io,resources=engulaclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EngulaCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *EngulaClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling engula cluster")

	original := &clusterv1alpha1.EngulaCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
	}
	err := r.fetch(ctx, req.Namespace, original)
	if err != nil {
		log.Error(err, "failed to retrieve EngulaCluster")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cluster := original.DeepCopy()
	cluster.Default()
	if len(cluster.Status.Conditions) == 0 {
		now := metav1.Now()
		cluster.Status.SetState(clusterv1alpha1.InitializedCondition, metav1.ConditionFalse, now)
	}

	return ctrl.Result{}, nil
}

func (r *EngulaClusterReconciler) fetch(ctx context.Context, ns string, o client.Object) (err error) {
	var (
		accessor metav1.Object
	)
	accessor, err = meta.Accessor(o)
	if err != nil {
		return
	}
	err = r.Client.Get(ctx, makeKey(ns, accessor.GetName()), o)
	return
}

func makeKey(ns, name string) types.NamespacedName {
	return types.NamespacedName{
		Name:      name,
		Namespace: ns,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *EngulaClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1alpha1.EngulaCluster{}).
		Complete(r)
}
