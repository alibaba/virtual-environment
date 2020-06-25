package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualEnvironmentSpec defines the desired state of VirtualEnvironment
// +k8s:openapi-gen=true
type VirtualEnvironmentSpec struct {
	// Pod label to mark virtual environment name
	EnvLabel EnvLabelSpec `json:"envLabel,omitempty"`
	// Header to keep env name in trace
	EnvHeader EnvHeaderSpec `json:"envHeader,omitempty"`
}

type EnvLabelSpec struct {
	// Name of the label
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// Symbol to split virtual env levels
	// +kubebuilder:validation:MaxLength=1
	// +kubebuilder:validation:MinLength=1
	Splitter string `json:"splitter,omitempty"`
	// Default subset to route when env header matches nothing
	// +kubebuilder:validation:MinLength=1
	DefaultSubset string `json:"defaultSubset,omitempty"`
}

type EnvHeaderSpec struct {
	// Name of the header
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// Other Names which also can be used for match env
	Aliases []EnvHeaderAliasSpec `json:"aliases,omitempty"`
	// Whether auto inject env header via sidecar
	AutoInject bool `json:"autoInject,omitempty"`
}

type EnvHeaderAliasSpec struct {
	// Alias name of the header
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// Regular expression to extract env tag from header value
	Pattern string `json:"pattern,omitempty"`
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
