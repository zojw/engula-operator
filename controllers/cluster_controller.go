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
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
	"time"

	clusterv1alpha1 "github.com/engula/engula-operator/api/v1alpha1"
	"github.com/engula/engula-operator/resource"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Clientset *kubernetes.Clientset
}

//+kubebuilder:rbac:groups=cluster.engula.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.engula.io,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.engula.io,resources=clusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=services/finalizers,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets/scale,verbs=get;watch;update
//+kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=statefulsets/finalizers,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;create;watch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;create;watch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;create;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling cluster")

	res := resource.NewManagedResources(ctx, req.Namespace, r.Client)

	cluster := resource.ClusterQueryObject(req.Name)
	err := res.Fetch(cluster)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "failed to retrieve cluster")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	cluster = cluster.DeepCopy()
	cluster.Default()

	if err = r.reconcileRBAC(res, log, cluster); err != nil {
		return ctrl.Result{}, err
	}

	labels := resource.MakeLabels(cluster.Name, resource.Cluster)
	var componentBuilders = []resource.ResourceBuilder{
		&resource.KernelBuilder{Own: cluster, Selector: labels},
		&resource.EngineBuilder{Own: cluster, Selector: labels},
		&resource.StorageBuilder{Own: cluster, Selector: labels},
		&resource.JournalBuilder{Own: cluster, Selector: labels},
		&resource.BackgroundBuilder{Own: cluster, Selector: labels},
	}
	for _, builder := range componentBuilders {
		_, err = res.ReconcileOne(builder)
		if err != nil {
			log.Error(err, "failed to reconcile cluster")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}
	}

	log.Info("reconcile cluster completed")
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) reconcileRBAC(res *resource.ManagedResources, log logr.Logger, cluster *clusterv1alpha1.Cluster) (err error) {
	var rbacReady bool
	if rbacReady, err = resource.CheckClusterServiceAccountExist(r.Clientset, cluster); err != nil {
		return
	}
	if !rbacReady {
		var rbacBuilders = []resource.ResourceBuilder{
			&resource.ServiceAccountBuilder{Cluster: cluster},
			&resource.RoleBuilder{Cluster: cluster},
			&resource.RoleBindingBuilder{Cluster: cluster},
		}
		for _, builder := range rbacBuilders {
			_, err = res.ReconcileOne(builder)
			if err != nil {
				log.Error(err, "failed to reconcile rbac")
				return
			}
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1alpha1.Cluster{}).
		Complete(r)
}
