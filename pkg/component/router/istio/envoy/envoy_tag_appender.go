package envoy

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	"context"
	"github.com/gogo/protobuf/jsonpb"
	pbtypes "github.com/gogo/protobuf/types"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	networkingv1alpha3api "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// delete EnvoyFilter instance if it already exist
func DeleteTagAppenderIfExist(client client.Client, namespace string, name string) error {
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, &networkingv1alpha3api.EnvoyFilter{})
	if err == nil {
		return shared.DeleteIns(client, namespace, name, &networkingv1alpha3api.EnvoyFilter{})
	}
	return nil
}

// check whether EnvoyFilter is different
func IsDifferentTagAppender(tagAppender *networkingv1alpha3api.EnvoyFilter, envLabel string, envHeader string) bool {
	return tagAppender.ObjectMeta.Labels["envLabel"] != envLabel || tagAppender.ObjectMeta.Labels["envHeader"] != envHeader
}

// generate EnvoyFilter to auto append env tag into HTTP header
func TagAppenderFilter(namespace string, name string, envLabel string, envHeader string) (*networkingv1alpha3api.EnvoyFilter, error) {
	patch, err := buildPatchStruct(envHeader)
	if patch == nil || patch.Fields == nil || err != nil {
		// unmarshal failed
		return nil, err
	}
	return &networkingv1alpha3api.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"envLabel":  envLabel,
				"envHeader": envHeader,
			},
		},
		Spec: networkingv1alpha3.EnvoyFilter{
			ConfigPatches: []*networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{{
				ApplyTo: networkingv1alpha3.EnvoyFilter_HTTP_FILTER,
				Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
					Context: networkingv1alpha3.EnvoyFilter_SIDECAR_OUTBOUND,
					ObjectTypes: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
						Listener: &networkingv1alpha3.EnvoyFilter_ListenerMatch{
							FilterChain: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
								Filter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
									Name: "envoy.http_connection_manager",
								},
							},
						},
					},
				},
				Patch: &networkingv1alpha3.EnvoyFilter_Patch{
					Operation: networkingv1alpha3.EnvoyFilter_Patch_INSERT_BEFORE,
					Value:     patch,
				},
			}},
		},
	}, nil
}

func buildPatchStruct(envHeader string) (*pbtypes.Struct, error) {
	config := `{
        "name": "virtual.environment.lua",
        "typed_config": {
            "@type": "type.googleapis.com/envoy.config.filter.http.lua.v2.Lua",
            "inline_code": "` + toOneLine(luaScript(envHeader)) + `"
        }
    }`
	unmarshalledConfig := &pbtypes.Struct{}
	err := jsonpb.Unmarshal(strings.NewReader(config), unmarshalledConfig)
	return unmarshalledConfig, err
}

func toOneLine(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, "\n", "\\n"), "\"", "\\\"")
}

// generate lua script to auto inject env tag from label to header
// Note: istio-proxy start with "--proxyLogLevel warning", only logWarn or logError will print
func luaScript(envHeader string) string {
	return strings.Trim(`
local curEnv = os.getenv("VIRTUAL_ENVIRONMENT_TAG")
function envoy_on_request(req)
  local env = req:headers():get("`+envHeader+`")
  if env == nil and curEnv ~= nil then
    req:headers():add("`+envHeader+`", curEnv)
  end
  local diagnose = req:headers():get("KT-ENV-DIAGNOSE")
  if diagnose ~= nil then
    if diagnose == "DEBUG" then
      req:logWarn("--- REQUEST HEADERS ---")
      for key, value in pairs(req:headers()) do
        req:logWarn(key .. " : " .. value)
      end
      req:logWarn("--- END OF DIAGNOSE ---")
    else
      if env ~= nil then
        req:logWarn("Env mark '`+envHeader+`' found as '" .. env .. "'")
      elseif curEnv ~= nil then
        req:logWarn("Env mark '`+envHeader+`' not found, added as '" .. curEnv .. "'")
      else
        req:logWarn("Env mark '`+envHeader+`' not found and current env unknown")
      end
    end
  end
end`, " \n")
}
