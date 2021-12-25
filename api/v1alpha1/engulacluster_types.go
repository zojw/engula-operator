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

	// (Optional) Default image for other components in cluster.
	// +optional
	Image *PodImage `json:"image,omitempty"`

	// (Optional) Kernel defines the desired state of Kernel service in cluster.
	// +optional
	Kernel *Kernel `json:"kernel,omitempty"`

	// (Optional) Storage defines the desired state of Storage service in cluster.
	// +optional
	Storage *Storage `json:"storage,omitempty"`

	// (Optional) Journal defines the desired state of Journal service in cluster.
	// +optional
	Journal *Journal `json:"journal,omitempty"`

	// (Optional) Engine defines the desired state of Engine service in cluster.
	// +optional
	Engine *Engine `json:"engine,omitempty"`

	// (Optional) Background defines the desired state of Background service in cluster.
	// +optional
	Background *Background `json:"background,omitempty"`
}

type Kernel struct {
	// Number of kernel replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for kernel
	// use EngulaCluster's image if not specified
	// +optional
	Image *PodImage `json:"image,omitempty"`
	// (Optional) Port use for expose service.
	// Default: 245677
	// +optional
	Port *int32 `json:"port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Disk volumn configuration.
	// +kubebuilder:validation:Required
	Volume Volume `json:"volume"`
}

type Engine struct {
	// Number of engine replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for engine.
	// use EngulaCluster's image if not specified
	// +optional
	Image *PodImage `json:"image,omitempty"`
	// (Optional) Port use for expose serivce.
	// Default: 245678
	// +optional
	Port *int32 `json:"port,omitempty"`
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
	// use EngulaCluster's image if not specified
	// +optional
	Image *PodImage `json:"image,omitempty"`
	// (Optional) Port use for expose serivce.
	// Default: 245679
	// +optional
	Port *int32 `json:"port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Disk volumn configuration.
	// +kubebuilder:validation:Required
	Volume Volume `json:"volume"`
}

type Journal struct {
	// Number of journal replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for journal.
	// use EngulaCluster's image if not specified
	// +optional
	Image *PodImage `json:"image,omitempty"`
	// (Optional) Port use for expose serivce.
	// Default: 245680
	// +optional
	Port *int32 `json:"port,omitempty"`
	// (Optional) Resource limits for container.
	// Default: (not specified)
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Disk volumn configuration.
	// +kubebuilder:validation:Required
	Volume Volume `json:"volume"`
}

type Background struct {
	// Number of background replicas (pods) in the cluster
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// (Optional) Image uses for background.
	// use EngulaCluster's image if not specified
	// +optional
	Image *PodImage `json:"image,omitempty"`
	// (Optional) Port use for expose serivce.
	// Default: 245681
	// +optional
	Port *int32 `json:"port,omitempty"`
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

type Volume struct {
	// (Optional) Directory from the host node's filesystem
	// +optional
	HostPath *corev1.HostPathVolumeSource `json:"hostPath,omitempty"`
	// (Optional) Persistent volume to use
	// +optional
	VolumeClaim *VolumeClaim `json:"pvc,omitempty"`
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

// EngulaClusterStatus defines the observed state of EngulaCluster
type EngulaClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// List of conditions represents the current state of cluster
	Conditions []ClusterCondition `json:"conditions,omitempty"`
}

func (s *EngulaClusterStatus) SetState(typ ClusterConditionType, status metav1.ConditionStatus, now metav1.Time) {
	c := s.findOrCreate(typ)
	if c.Status == status {
		return
	}
	c.Status = status
	c.LastUpdateTime = now
}

func (s *EngulaClusterStatus) findOrCreate(typ ClusterConditionType) *ClusterCondition {
	var idx = -1
	for i, c := range s.Conditions {
		if c.Type == typ {
			idx = i
			break
		}
	}
	if idx >= 0 {
		return &s.Conditions[idx]
	}

	s.Conditions = append(s.Conditions, ClusterCondition{
		Type:           typ,
		Status:         metav1.ConditionUnknown,
		LastUpdateTime: metav1.Now(),
	})
	return &s.Conditions[len(s.Conditions)-1]
}

type ClusterCondition struct {
	// Type of the condition
	// +kubebuilder:validation:Required
	Type ClusterConditionType `json:"type"`
	// Condition status: True, False or Unknown
	// +kubebuilder:validation:Required
	Status metav1.ConditionStatus `json:"status"`
	// The last time for condition updated
	// +kubebuilder:validation:Required
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

type ClusterConditionType string

const (
	//InitializedCondition string
	InitializedCondition ClusterConditionType = "Initialized"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EngulaCluster is the Schema for the engulaclusters API
type EngulaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EngulaClusterSpec   `json:"spec,omitempty"`
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
