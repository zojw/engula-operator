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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KernelSpec defines the desired state of Kernel
type KernelSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Number of kernel replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for kernel
	// use EngulaCluster's image if not specified
	// +optional
	Image *PodImage `json:"image,omitempty"`
	// (Optional) Port use for expose service.
	// Default: 24567
	// +optional
	Port *int32 `json:"port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resource,omitempty"`
	// Disk volumn configuration.
	// +kubebuilder:validation:Required
	Volume VolumeClaim `json:"volume"`
}

// KernelStatus defines the observed state of Kernel
type KernelStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Kernel is the Schema for the kernels API
type Kernel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KernelSpec   `json:"spec,omitempty"`
	Status KernelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KernelList contains a list of Kernel
type KernelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kernel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kernel{}, &KernelList{})
}
