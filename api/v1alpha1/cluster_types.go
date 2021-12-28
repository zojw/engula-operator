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

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the name of EngulaCluster.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// (Optional) Default image for other components in cluster.
	// +optional
	Image *PodImage `json:"image,omitempty"`

	// (Optional) Kernel defines the desired state of Kernel service in cluster.
	// +optional
	Kernel *KernelSpec `json:"kernel,omitempty"`

	// (Optional) Storage defines the desired state of Storage service in cluster.
	// +optional
	Storage *StorageSpec `json:"storage,omitempty"`

	// (Optional) Journal defines the desired state of Journal service in cluster.
	// +optional
	Journal *JournalSpec `json:"journal,omitempty"`

	// (Optional) Engine defines the desired state of Engine service in cluster.
	// +optional
	Engine *EngineSpec `json:"engine,omitempty"`

	// (Optional) Background defines the desired state of Background service in cluster.
	// +optional
	Background *BackgroundSpec `json:"background,omitempty"`
}

type PodImage struct {
	// Container image with version.
	// For instance: engula/engula:v0.3
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// (Optional) PullPolicy for the image, which defaults to IfNotPresent.
	// Default: IfNotPresent
	// +optional
	PullPolicyName *corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

// VolumeClaim wraps a persistent volume claim (PVC) to use with the container.
// Only one of the fields should set
type VolumeClaim struct {
	// (Optional) PVC to request a new persistent volume
	// +optional
	PersistentVolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"spec,omitempty"`
	// (Optional) Existing PVC in the same namespace
	// +optional
	PersistentVolumeSource corev1.PersistentVolumeClaimVolumeSource `json:"source,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
