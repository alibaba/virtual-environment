package http

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	networkingv1alpha3 "knative.dev/pkg/apis/istio/v1alpha3"
	"regexp"
)

// replace invalid chars in subset name
func toSubsetName(labelValue string) string {
	re, _ := regexp.Compile("[_.]")
	return re.ReplaceAllString(labelValue, "-")
}

// merge two http route structures
func mergeRoute(a, b *networkingv1alpha3.HTTPRoute) (*networkingv1alpha3.HTTPRoute, error) {
	if a == nil {
		return b, nil
	}
	jb, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	var merged *networkingv1alpha3.HTTPRoute
	err = deepCopy(&merged, a)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jb, &merged)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

// deep copy an object
func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
