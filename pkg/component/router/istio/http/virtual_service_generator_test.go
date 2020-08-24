package http

import (
	"strconv"
	"testing"
)

func TestVirtualServiceMatchRoute(t *testing.T) {
	deployments := []string{"a.b.c", "a.b", "a", "a.d", "a.d.e.f.g"}
	routes, ok := virtualServiceMatchRoute("testSvc", deployments, "a.b", "test",
		nil, ".", 0, "dev", 1)
	route := routes[0]
	if !ok || route.Match[0].Headers["test"].Exact != "a.b" || route.Route[0].Destination.Subset != "a-b" {
		t.Errorf("Match failed 1 : " + strconv.FormatBool(ok) + ", match: " + route.Match[0].Headers["test"].Exact +
			", subset: " + route.Route[0].Destination.Subset)
	}
	routes, ok = virtualServiceMatchRoute("testSvc", deployments, "a.d.e.f", "test",
		nil, ".", 0, "dev", 1)
	route = routes[0]
	if !ok || route.Match[0].Headers["test"].Exact != "a.d.e.f" || route.Route[0].Destination.Subset != "a-d" {
		t.Errorf("Match failed 2 : " + strconv.FormatBool(ok) + ", match: " + route.Match[0].Headers["test"].Exact +
			", subset: " + route.Route[0].Destination.Subset)
	}
	routes, ok = virtualServiceMatchRoute("testSvc", deployments, "b.x", "test",
		nil, ".", 0, "dev", 1)
	if ok {
		t.Errorf("Match failed 3")
	}
}

func TestFindLongestString(t *testing.T) {
	if findLongestString([]string{"abc", "defgh", "ij", "k"}) != "defgh" {
		t.Fail()
	}
	if findLongestString([]string{"abc", "ij", "k"}) != "abc" {
		t.Fail()
	}
	if findLongestString([]string{"ij", "k", "abc"}) != "abc" {
		t.Fail()
	}
}

func TestLeveledEqual(t *testing.T) {
	if !leveledEqual("top.second.third", "top.second.third", ".") {
		t.Fail()
	}
	if !leveledEqual("top.second.third", "top.second", ".") {
		t.Fail()
	}
	if !leveledEqual("top.second.third", "top", ".") {
		t.Fail()
	}
	if leveledEqual("top.second", "top.second.third", ".") {
		t.Fail()
	}
	if leveledEqual("top.second.third", "top.second.", ".") {
		t.Fail()
	}
}
