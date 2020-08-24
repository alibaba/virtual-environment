package http

import (
	"knative.dev/pkg/apis/istio/common/v1alpha1"
	networkingv1alpha3 "knative.dev/pkg/apis/istio/v1alpha3"
	"testing"
)

func TestMergeRoute(t *testing.T) {
	a := networkingv1alpha3.HTTPRoute{
		Match: []networkingv1alpha3.HTTPMatchRequest{{
			URI: &v1alpha1.StringMatch{
				Exact: "/abc",
			},
		}},
		Rewrite: &networkingv1alpha3.HTTPRewrite{
			URI: "/def",
		},
	}
	b := networkingv1alpha3.HTTPRoute{
		Match: []networkingv1alpha3.HTTPMatchRequest{{
			Headers: map[string]v1alpha1.StringMatch{
				"xxx": {Exact: "yyy"},
			},
		}},
		Route: []networkingv1alpha3.HTTPRouteDestination{{
			Destination: networkingv1alpha3.Destination{
				Host: "app",
			},
		}},
	}

	c, err := mergeRoute(&a, &b)
	if err != nil {
		t.Errorf("merge failed")
	}
	if c.Rewrite.URI != "/def" {
		t.Errorf("failed to merge rewrite uri: " + c.Rewrite.URI)
	}
	if c.Match[0].URI.Exact != "/abc" {
		t.Errorf("failed to merge match uri: " + c.Match[0].URI.Exact)
	}
	if c.Match[0].Headers["xxx"].Exact != "yyy" {
		t.Errorf("failed to merge match header: " + c.Rewrite.URI)
	}
}
