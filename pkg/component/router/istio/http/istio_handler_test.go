package http

import (
	"alibaba.com/virtual-env-operator/pkg/component/routenent/route"
	"testing"
)

func TestDestinationRuleMatchSubset(t *testing.T) {
	rule := route.destinationRuleMatchSubset("test", "demo")
	if rule.Name != "demo" || rule.Labels["test"] != "demo" {
		t.Fail()
	}
}

func TestGetKeys(t *testing.T) {
	keys := route.getKeys(map[string]bool{"a": true, "c": false, "e": true})
	if len(keys) != 3 {
		t.Fail()
	}
	if keys[0] != "a" || keys[1] != "c" || keys[2] != "e" {
		t.Fail()
	}
}

func TestVirtualServiceMatchRoute(t *testing.T) {
	deployments := map[string]string{"dep1": "a.b.c", "dep2": "a.b", "dep3": "a", "dep4": "a.d", "dep5": "a.d.e.f.g"}
	route, ok := route.virtualServiceMatchRoute("testSvc", deployments, "a.b", "test", ".", "dev")
	if !ok || route.Match[0].Headers["test"].Exact != "a.b" || route.Route[0].Destination.Subset != "a.b" {
		t.Fail()
	}
	route, ok = route.virtualServiceMatchRoute("testSvc", deployments, "a.d.e.f", "test", ".", "dev")
	if !ok || route.Match[0].Headers["test"].Exact != "a.d.e.f" || route.Route[0].Destination.Subset != "a.d" {
		t.Fail()
	}
	route, ok = route.virtualServiceMatchRoute("testSvc", deployments, "b.x", "test", ".", "dev")
	if ok {
		t.Fail()
	}
}

func TestFindLongestString(t *testing.T) {
	if route.findLongestString([]string{"abc", "defgh", "ij", "k"}) != "defgh" {
		t.Fail()
	}
	if route.findLongestString([]string{"abc", "ij", "k"}) != "abc" {
		t.Fail()
	}
	if route.findLongestString([]string{"ij", "k", "abc"}) != "abc" {
		t.Fail()
	}
}

func TestLeveledEqual(t *testing.T) {
	if !route.leveledEqual("top.second.third", "top.second.third", ".") {
		t.Fail()
	}
	if !route.leveledEqual("top.second.third", "top.second", ".") {
		t.Fail()
	}
	if !route.leveledEqual("top.second.third", "top", ".") {
		t.Fail()
	}
	if route.leveledEqual("top.second", "top.second.third", ".") {
		t.Fail()
	}
	if route.leveledEqual("top.second.third", "top.second.", ".") {
		t.Fail()
	}
}
