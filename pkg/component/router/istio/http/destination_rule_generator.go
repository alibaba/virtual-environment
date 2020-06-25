package http

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	networkingv1alpha3 "knative.dev/pkg/apis/istio/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// generate istio destination rule instance
func DestinationRule(namespace string, svcName string, relatedDeployments map[string]string,
	envLabel string) *networkingv1alpha3.DestinationRule {
	destRule := &networkingv1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.DestinationRuleSpec{
			Host:    svcName,
			Subsets: []networkingv1alpha3.Subset{},
		},
	}
	for _, label := range relatedDeployments {
		destRule.Spec.Subsets = append(destRule.Spec.Subsets, destinationRuleMatchSubset(envLabel, label))
	}
	return destRule
}

// delete DestinationRule
func DeleteDestinationRule(client client.Client, namespace string, name string) error {
	return shared.DeleteIns(client, namespace, name, &networkingv1alpha3.DestinationRule{})
}

// check whether DestinationRule is different
func IsDifferentDestinationRule(spec1 *networkingv1alpha3.DestinationRuleSpec,
	spec2 *networkingv1alpha3.DestinationRuleSpec, label string) bool {
	if len(spec1.Subsets) != len(spec2.Subsets) {
		return true
	}
	for _, subset1 := range spec1.Subsets {
		subset2 := findSubsetByName(spec2.Subsets, subset1.Name)
		if subset2 == nil {
			return true
		}
		if subset1.Labels[label] != subset2.Labels[label] {
			return true
		}
	}
	return false
}

// find subset from list
func findSubsetByName(subsets []networkingv1alpha3.Subset, name string) *networkingv1alpha3.Subset {
	for _, subset := range subsets {
		if subset.Name == name {
			return &subset
		}
	}
	return nil
}

// generate istio destination rule subset instance
func destinationRuleMatchSubset(labelKey string, labelValue string) networkingv1alpha3.Subset {
	return networkingv1alpha3.Subset{
		Name: toSubsetName(labelValue),
		Labels: map[string]string{
			labelKey: labelValue,
		},
	}
}
