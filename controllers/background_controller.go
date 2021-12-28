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
	"strconv"
	"time"

	clusterv1alpha1 "github.com/engula/engula-operator/api/v1alpha1"
	"github.com/engula/engula-operator/resource"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// BackgroundReconciler reconciles a Background object
type BackgroundReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cluster.engula.io,resources=backgrounds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.engula.io,resources=backgrounds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.engula.io,resources=backgrounds/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=services/finalizers,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets/scale,verbs=get;watch;update
//+kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=statefulsets/finalizers,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *BackgroundReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling background")

	res := resource.NewManagedResources(ctx, req.Namespace, r.Client)

	background := resource.BackgroundQueryObject(req.Name)
	err := res.Fetch(background)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "failed to retrieve background")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	background = background.DeepCopy()
	background.Default()

	var clusterName string
	if clusterName, err = resource.GetClusterName(background.ObjectMeta); err != nil {
		return ctrl.Result{}, err
	}
	labels := resource.MakeLabels(clusterName, resource.Background)
	background = background.DeepCopy()
	var builders = []resource.ResourceBuilder{
		resource.NewServiceBuilder(req.Name, labels, func() *resource.ServiceResource {
			return &resource.ServiceResource{
				Port: background.Spec.Port,
			}
		}),
		resource.NewStatefulsetBuilder(req.Name, clusterName,
			labels,
			func() *resource.StatefulResource {
				return &resource.StatefulResource{
					Replicas:  background.Spec.Replicas,
					Image:     background.Spec.Image,
					Port:      background.Spec.Port,
					Resources: background.Spec.Resources,
				}
			},
			resource.Background,
			"0.0.0.0:"+strconv.Itoa(int(*background.Spec.Port)),
		),
	}

	for _, builder := range builders {
		if _, err = res.ReconcileOne(builder); err != nil {
			log.Error(err, "failed to reconcile cluster")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}
	}

	log.Info("reconcile background completed")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackgroundReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1alpha1.Background{}).
		Complete(r)
}
