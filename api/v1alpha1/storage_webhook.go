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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var storagelog = logf.Log.WithName("storage-resource")

func (r *Storage) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-cluster-engula-io-v1alpha1-storage,mutating=true,failurePolicy=fail,sideEffects=None,groups=cluster.engula.io,resources=storages,verbs=create;update,versions=v1alpha1,name=mstorage.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Storage{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Storage) Default() {
	storagelog.Info("default", "name", r.Name)
	r.Spec.Default()
}

func (s *StorageSpec) Default() {
	if s != nil {
		if s.Port == nil {
			s.Port = &DefaultStoragePort
		}
		setDefaultPollPolicy(s.Image)
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-cluster-engula-io-v1alpha1-storage,mutating=false,failurePolicy=fail,sideEffects=None,groups=cluster.engula.io,resources=storages,verbs=create;update,versions=v1alpha1,name=vstorage.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Storage{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Storage) ValidateCreate() error {
	storagelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Storage) ValidateUpdate(old runtime.Object) error {
	storagelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Storage) ValidateDelete() error {
	storagelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}