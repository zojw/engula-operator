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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	DefaultKernelPort     int32 = 245677
	DefaultEnginePort     int32 = 245678
	DefaultStoragePort    int32 = 245679
	DefaultJournalPort    int32 = 245680
	DefaultBackgroundPort int32 = 245681
)

// log is for logging in this package.
var engulaclusterlog = logf.Log.WithName("engulacluster-resource")

func (r *EngulaCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-cluster-engula-io-v1alpha1-engulacluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=cluster.engula.io,resources=engulaclusters,verbs=create;update,versions=v1alpha1,name=mengulacluster.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &EngulaCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *EngulaCluster) Default() {
	engulaclusterlog.Info("default", "name", r.Name)

	if r.Spec.Kernel != nil {
		if r.Spec.Kernel.Port != nil {
			r.Spec.Kernel.Port = &DefaultKernelPort
		}
		setDefaultPollPolicy(r.Spec.Kernel.Image)
	}

	if r.Spec.Storage != nil {
		if r.Spec.Storage.Port != nil {
			r.Spec.Storage.Port = &DefaultStoragePort
		}
		setDefaultPollPolicy(r.Spec.Storage.Image)
	}

	if r.Spec.Journal != nil {
		if r.Spec.Journal.Port != nil {
			r.Spec.Journal.Port = &DefaultJournalPort
		}
		setDefaultPollPolicy(r.Spec.Journal.Image)
	}

	if r.Spec.Background != nil {
		if r.Spec.Background.Port != nil {
			r.Spec.Background.Port = &DefaultBackgroundPort
		}
		setDefaultPollPolicy(r.Spec.Background.Image)
	}

	if r.Spec.Engine != nil {
		if r.Spec.Engine.Port != nil {
			r.Spec.Engine.Port = &DefaultEnginePort
		}
		setDefaultPollPolicy(r.Spec.Engine.Image)
	}
}

func setDefaultPollPolicy(image *PodImage) {
	if image != nil && image.PullPolicyName == nil {
		p := v1.PullIfNotPresent
		image.PullPolicyName = &p
	}
}

//+kubebuilder:webhook:path=/validate-cluster-engula-io-v1alpha1-engulacluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=cluster.engula.io,resources=engulaclusters,verbs=create;update,versions=v1alpha1,name=vengulacluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &EngulaCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *EngulaCluster) ValidateCreate() error {
	engulaclusterlog.Info("validate create", "name", r.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *EngulaCluster) ValidateUpdate(old runtime.Object) error {
	engulaclusterlog.Info("validate update", "name", r.Name)
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *EngulaCluster) ValidateDelete() error {
	engulaclusterlog.Info("validate delete", "name", r.Name)
	return nil
}
