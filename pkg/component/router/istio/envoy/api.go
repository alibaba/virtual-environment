package envoy

import (
	networkingv1alpha3api "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// Package-wide variables from generator "register".
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha3"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

const (
	// Package-wide consts from generator "register".
	GroupName = "networking.istio.io"
)

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&networkingv1alpha3api.EnvoyFilter{},
		&networkingv1alpha3api.EnvoyFilterList{},
	)
	v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
