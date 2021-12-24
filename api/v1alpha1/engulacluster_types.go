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

// EngulaClusterSpec defines the desired state of EngulaCluster
type EngulaClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the name of EngulaCluster.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// (Optional) Kernel defines the desired state of Kernel service in cluster.
	// +optional
	Kernel Kernel `json:"kernel"`

	// (Optional) Storage defines the desired state of Storage service in cluster.
	// +optional
	Storage Storage `json:"storage"`

	// (Optional) Journal defines the desired state of Journal service in cluster.
	// +optional
	Journal Journal `json:"journal"`
}

type Kernel struct {
	// Number of kernel replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for kernel.
	// +optional
	Image *PodImage `json:"image"`
	// (Optional) Port use for expose serivce.
	// Default: 245677
	// +optional
	Port *int32 `json:"Port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type Engine struct {
	// Number of engine replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for engine.
	// +optional
	Image *PodImage `json:"image"`
	// (Optional) Port use for expose serivce.
	// Default: 245678
	// +optional
	Port *int32 `json:"Port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type Storage struct {
	// Number of storage replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for storage.
	// +optional
	Image *PodImage `json:"image"`
	// (Optional) Port use for expose serivce.
	// Default: 245679
	// +optional
	Port *int32 `json:"Port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type Journal struct {
	// Number of journal replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for journal.
	// +optional
	Image *PodImage `json:"image"`
	// (Optional) Port use for expose serivce.
	// Default: 245680
	// +optional
	Port *int32 `json:"Port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type Background struct {
	// Number of background replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for background.
	// +optional
	Image *PodImage `json:"image"`
	// (Optional) Port use for expose serivce.
	// Default: 245681
	// +optional
	Port *int32 `json:"Port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type PodImage struct {
	// Container image with version.
	// For instance: engula/engula:v0.3
	// +required
	Name string `json:"name"`
	// (Optional) PullPolicy for the image, which defaults to IfNotPresent.
	// Default: IfNotPresent
	// +optional
	PullPolicyName *corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

// EngulaClusterStatus defines the observed state of EngulaCluster
type EngulaClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EngulaCluster is the Schema for the engulaclusters API
type EngulaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EngulaClusterSpec   `json:"sÏpec,omitempty"`
	Status EngulaClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EngulaClusterList contains a list of EngulaCluster
type EngulaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EngulaCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EngulaCluster{}, &EngulaClusterList{})
}
