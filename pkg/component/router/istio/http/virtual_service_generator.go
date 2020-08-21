package http

import (
	envv1alpha2 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha2"
	"alibaba.com/virtual-env-operator/pkg/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis/istio/common/v1alpha1"
	networkingv1alpha3 "knative.dev/pkg/apis/istio/v1alpha3"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// generate istio virtual service instance
func VirtualService(namespace string, svcName string, availableLabels []string, relatedDeployments []string,
	spec envv1alpha2.VirtualEnvironmentSpec) *networkingv1alpha3.VirtualService {
	envHeaderName := spec.EnvHeader.Name
	envHeaderAliases := spec.EnvHeader.Aliases
	envSplitter := spec.EnvLabel.Splitter
	defaultSubset := spec.EnvLabel.DefaultSubset
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
	serviceInfo := shared.AvailableServices[svcName]
	if len(serviceInfo.Gateways) > 0 {
		virtualSvc.Spec.Gateways = serviceInfo.Gateways
	}
	if len(serviceInfo.Hosts) > 0 {
		virtualSvc.Spec.Hosts = serviceInfo.Hosts
	}
	for _, port := range serviceInfo.Ports {
		for _, label := range availableLabels {
			matchRoutes, ok := virtualServiceMatchRoute(svcName, relatedDeployments, label, envHeaderName,
				envHeaderAliases, envSplitter, port, toSubsetName(defaultSubset), len(serviceInfo.Ports))
			if ok {
				for _, matchRoute := range matchRoutes {
					virtualSvc.Spec.HTTP = append(virtualSvc.Spec.HTTP, matchRoute)
				}
			}
		}
		virtualSvc.Spec.HTTP = append(virtualSvc.Spec.HTTP,
			defaultRoute(svcName, port, toSubsetName(defaultSubset), len(serviceInfo.Ports)))
	}
	return virtualSvc
}

// delete VirtualService
func DeleteVirtualService(client client.Client, namespace string, name string) error {
	return shared.DeleteIns(client, namespace, name, &networkingv1alpha3.VirtualService{})
}

// check whether VirtualService is different
func IsDifferentVirtualService(spec1 *networkingv1alpha3.VirtualServiceSpec, spec2 *networkingv1alpha3.VirtualServiceSpec,
	header string) bool {
	if !reflect.DeepEqual(spec1.Gateways, spec2.Gateways) {
		return true
	}
	if !reflect.DeepEqual(spec1.Hosts, spec2.Hosts) {
		return true
	}
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
	return route.Route[0].Destination.Subset == target.Route[0].Destination.Subset &&
		route.Route[0].Destination.Port.Number == target.Route[0].Destination.Port.Number
}

// calculate and generate http route instance
func virtualServiceMatchRoute(serviceName string, relatedDeployments []string, labelVal string,
	headerKeyName string, headerKeyAlias []envv1alpha2.EnvHeaderAliasSpec, splitter string, port uint32,
	defaultSubset string, totalPortCount int) ([]networkingv1alpha3.HTTPRoute, bool) {
	possibleRoutes := getPossibleRoutes(relatedDeployments, labelVal, splitter)
	if len(possibleRoutes) > 0 {
		var subsetName = toSubsetName(findLongestString(possibleRoutes))
		if defaultSubset != subsetName {
			var routes = []networkingv1alpha3.HTTPRoute{
				matchRouteExact(serviceName, headerKeyName, labelVal, port, subsetName, totalPortCount),
			}
			if headerKeyAlias != nil {
				for _, alias := range headerKeyAlias {
					var route networkingv1alpha3.HTTPRoute
					if alias.Pattern != "" && alias.Placeholder != "" {
						regexLabelVal := strings.ReplaceAll(alias.Pattern, alias.Placeholder, labelVal)
						route = matchRouteRegex(serviceName, alias.Name, regexLabelVal, port, subsetName, totalPortCount)
					} else {
						route = matchRouteExact(serviceName, alias.Name, labelVal, port, subsetName, totalPortCount)
					}
					routes = append(routes, route)
				}
			}
			return routes, true
		}
	}
	return []networkingv1alpha3.HTTPRoute{}, false
}

// fetch all route rule for specified label
func getPossibleRoutes(relatedDeployments []string, labelVal string, splitter string) []string {
	var possibleRoutes []string
	for _, v := range relatedDeployments {
		if leveledEqual(labelVal, v, splitter) {
			possibleRoutes = append(possibleRoutes, v)
		}
	}
	return possibleRoutes
}

// generate default http route instance
func defaultRoute(name string, port uint32, defaultSubset string, totalPortCount int) networkingv1alpha3.HTTPRoute {
	route := networkingv1alpha3.HTTPRoute{
		Route: generateHttpRoute(name, port, defaultSubset),
	}
	if totalPortCount > 1 {
		route.Match = []networkingv1alpha3.HTTPMatchRequest{{
			Port: port,
		}}
	}
	return route
}

// generate istio virtual service http route instance with regex match
func matchRouteRegex(serviceName string, headerKey string, labelVal string, port uint32,
	subsetName string, totalPortCount int) networkingv1alpha3.HTTPRoute {
	route := matchRoute(serviceName, port, subsetName, totalPortCount)
	route.Match[0].Headers[headerKey] = v1alpha1.StringMatch{ Regex: labelVal }
	return route
}

// generate istio virtual service http route instance with exact match
func matchRouteExact(serviceName string, headerKey string, labelVal string, port uint32,
	subsetName string, totalPortCount int) networkingv1alpha3.HTTPRoute {
	route := matchRoute(serviceName, port, subsetName, totalPortCount)
	route.Match[0].Headers[headerKey] = v1alpha1.StringMatch{ Exact: labelVal }
	return route
}

// generate istio virtual service http route instance
func matchRoute(serviceName string, port uint32, subsetName string, totalPortCount int) networkingv1alpha3.HTTPRoute {
	route := networkingv1alpha3.HTTPRoute{
		Route: generateHttpRoute(serviceName, port, subsetName),
		Match: []networkingv1alpha3.HTTPMatchRequest{{
			Headers: map[string]v1alpha1.StringMatch{},
		}},
	}
	if totalPortCount > 1 {
		route.Match[0].Port = port
	}
	return route
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
