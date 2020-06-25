package http

import (
	"testing"
)

func TestDestinationRuleMatchSubset(t *testing.T) {
	rule := destinationRuleMatchSubset("test", "demo")
	if rule.Name != "demo" || rule.Labels["test"] != "demo" {
		t.Fail()
	}
}
