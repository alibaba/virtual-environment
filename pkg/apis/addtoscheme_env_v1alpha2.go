package apis

import (
	"alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha2"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha2.SchemeBuilder.AddToScheme)
}
