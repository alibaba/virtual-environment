package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualEnvironmentSpec defines the desired state of VirtualEnvironment
// +k8s:openapi-gen=true
type VirtualEnvironmentSpec struct {
	// Default subset to route when env header matches nothing
	// +kubebuilder:validation:MinLength=1
	DefaultSubset string `json:"defaultSubset,omitempty"`
	// Header to keep env name in trace
	// +kubebuilder:validation:MinLength=1
	EnvHeader string `json:"envHeader,omitempty"`
	// Environment variable to mark env name of deployment
	// +kubebuilder:validation:MinLength=1
	EnvLabel string `json:"envLabel,omitempty"`
	// Symbol to split virtual env levels
	// +kubebuilder:validation:MaxLength=1
	// +kubebuilder:validation:MinLength=1
	EnvSplitter string `json:"envSplitter,omitempty"`
}

// VirtualEnvironmentStatus defines the observed state of VirtualEnvironment
// +k8s:openapi-gen=true
type VirtualEnvironmentStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualEnvironment is the Schema for the virtualenvironments API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=virtualenvironments,scope=Namespaced
type VirtualEnvironment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualEnvironmentSpec   `json:"spec,omitempty"`
	Status VirtualEnvironmentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualEnvironmentList contains a list of VirtualEnvironment
type VirtualEnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualEnvironment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualEnvironment{}, &VirtualEnvironmentList{})
}
