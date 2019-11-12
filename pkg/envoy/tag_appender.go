package envoy

import (
	protobuftypes "github.com/gogo/protobuf/types"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	networkingv1alpha3api "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnvoyFilter to auto append env tag into HTTP header
func TagAppenderFilter(namespace string, name string, envLabel string, envHeader string) *networkingv1alpha3api.EnvoyFilter {
	return &networkingv1alpha3api.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.EnvoyFilter{
			Filters: []*networkingv1alpha3.EnvoyFilter_Filter{{
				ListenerMatch: &networkingv1alpha3.EnvoyFilter_DeprecatedListenerMatch{
					ListenerType: networkingv1alpha3.EnvoyFilter_DeprecatedListenerMatch_SIDECAR_OUTBOUND,
				},
				FilterName: "envoy.lua",
				FilterType: networkingv1alpha3.EnvoyFilter_Filter_HTTP,
				FilterConfig: &protobuftypes.Struct{
					Fields: map[string]*protobuftypes.Value{
						"inlineCode": {
							Kind: &protobuftypes.Value_StringValue{
								StringValue: luaScript(envLabel, envHeader),
							},
						},
					},
				},
			}},
		},
	}
}

// generate lua script to auto inject env tag from label to header
func luaScript(envLabel string, envHeader string) string {
	return `
	  local envLabel = "` + envLabel + `"
	  local envHeader = "` + envHeader + `"
	  local labels = os.getenv ("ISTIO_METAJSON_LABELS")
	  local beginPos, endPos, curEnv
	  _, beginPos = string.find(labels, '","' .. envLabel .. '":"')
	  if beginPos ~= nil then
		endPos = string.find(labels, '"', beginPos + 1)
		if endPos ~= nil and endPos > beginPos then
		  curEnv = string.sub(labels, beginPos + 1, endPos - 1)
		end
	  end
	  function envoy_on_request(request_handle)
		local env = request_handle:headers()[envHeader]
		if env == nil and curEnv ~= nil then
		  request_handle:headers():add(envHeader, curEnv)
		end
	  end
	`
}
