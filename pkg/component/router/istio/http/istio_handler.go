package http

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis/istio/common/v1alpha1"
	networkingv1alpha3 "knative.dev/pkg/apis/istio/v1alpha3"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// generate istio virtual service instance
func VirtualService(namespace string, svcName string, availableLabels []string,
	relatedDeployments map[string]string, envHeader string, envSplitter string,
	defaultSubset string) *networkingv1alpha3.VirtualService {
	virtualSvc := &networkingv1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.VirtualServiceSpec{
			Hosts: []string{svcName},
			HTTP:  []networkingv1alpha3.HTTPRoute{},
		},
	}
	for _, port := range shared.AvailableServicePorts[svcName] {
		for _, label := range availableLabels {
			matchRoute, ok := virtualServiceMatchRoute(svcName, relatedDeployments,
				label, envHeader, envSplitter, port, defaultSubset)
			if ok {
				virtualSvc.Spec.HTTP = append(virtualSvc.Spec.HTTP, matchRoute)
			}
		}
		virtualSvc.Spec.HTTP = append(virtualSvc.Spec.HTTP, defaultRoute(svcName, port, toSubsetName(defaultSubset)))
	}
	return virtualSvc
}

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

// delete VirtualService
func DeleteVirtualService(client client.Client, namespace string, name string) error {
	return shared.DeleteIns(client, namespace, name, &networkingv1alpha3.VirtualService{})
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

// check whether VirtualService is different
func IsDifferentVirtualService(spec1 *networkingv1alpha3.VirtualServiceSpec, spec2 *networkingv1alpha3.VirtualServiceSpec,
	header string) bool {
	if len(spec1.HTTP) != len(spec2.HTTP) {
		return true
	}
	for _, route1 := range spec1.HTTP {
		if !findMatchRoute(spec2.HTTP, &route1, header) {
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

// check whether HTTPRoute exist in list
func findMatchRoute(routes []networkingv1alpha3.HTTPRoute, target *networkingv1alpha3.HTTPRoute, header string) bool {
	for _, route := range routes {
		if isRouteEqual(&route, target, header) {
			return true
		}
	}
	return false
}

// compare whether route rule is equal
func isRouteEqual(route *networkingv1alpha3.HTTPRoute, target *networkingv1alpha3.HTTPRoute, header string) bool {
	if route.Match == nil || target.Match == nil {
		return route.Match == nil && target.Match == nil && isDestinationEqual(route, target)
	} else if len(route.Match) == 0 || len(target.Match) == 0 {
		return len(route.Match) == 0 && len(target.Match) == 0 && isDestinationEqual(route, target)
	} else {
		return route.Match[0].Headers[header] == target.Match[0].Headers[header] && isDestinationEqual(route, target)
	}
}

// compare whether route destination is equal
func isDestinationEqual(route *networkingv1alpha3.HTTPRoute, target *networkingv1alpha3.HTTPRoute) bool {
	return route.Route[0].Destination.Subset == target.Route[0].Destination.Subset
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

// calculate and generate http route instance
func virtualServiceMatchRoute(serviceName string, relatedDeployments map[string]string, labelVal string, headerKey string,
	splitter string, port uint32, defaultSubset string) (networkingv1alpha3.HTTPRoute, bool) {
	var possibleRoutes []string
	for _, v := range relatedDeployments {
		if leveledEqual(labelVal, v, splitter) {
			possibleRoutes = append(possibleRoutes, v)
		}
	}
	if len(possibleRoutes) > 0 {
		var subsetName = toSubsetName(findLongestString(possibleRoutes))
		if defaultSubset != subsetName {
			return matchRoute(serviceName, headerKey, labelVal, port, subsetName), true
		}
	}
	return networkingv1alpha3.HTTPRoute{}, false
}

// replace invalid chars in subset name
func toSubsetName(labelValue string) string {
	re, _ := regexp.Compile("[_.]")
	return re.ReplaceAllString(labelValue, "-")
}

// generate istio route
func generateHttpRoute(serviceName string, port uint32, subsetName string) []networkingv1alpha3.HTTPRouteDestination {
	return []networkingv1alpha3.HTTPRouteDestination{{
		Destination: networkingv1alpha3.Destination{
			Host:   serviceName,
			Subset: subsetName,
			Port:   networkingv1alpha3.PortSelector{Number: port},
		},
		Weight: 100,
	}}
}

// generate default http route instance
func defaultRoute(name string, port uint32, defaultSubset string) networkingv1alpha3.HTTPRoute {
	return networkingv1alpha3.HTTPRoute{
		Route: generateHttpRoute(name, port, defaultSubset),
		Match: []networkingv1alpha3.HTTPMatchRequest{{
			Port: port,
		}},
	}
}

// generate istio virtual service http route instance
func matchRoute(serviceName string, headerKey string, labelVal string, port uint32,
	subsetName string) networkingv1alpha3.HTTPRoute {
	return networkingv1alpha3.HTTPRoute{
		Route: generateHttpRoute(serviceName, port, subsetName),
		Match: []networkingv1alpha3.HTTPMatchRequest{{
			Headers: map[string]v1alpha1.StringMatch{
				headerKey: {
					Exact: labelVal,
				},
			},
			Port: port,
		}},
	}
}

// get the longest string in list
func findLongestString(strings []string) string {
	mostLongStr := ""
	for _, str := range strings {
		if len(str) > len(mostLongStr) {
			mostLongStr = str
		}
	}
	return mostLongStr
}

// check whether source string match target string at any level
func leveledEqual(source string, target string, splitter string) bool {
	for {
		if source == target {
			return true
		}
		if strings.Contains(source, splitter) {
			source = source[0:strings.LastIndex(source, splitter)]
		} else {
			return false
		}
	}
}
