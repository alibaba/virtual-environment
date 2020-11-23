package http

import (
	"knative.dev/pkg/apis/istio/common/v1alpha1"
	"strconv"
	"testing"
)

func TestVirtualServiceMatchRoute(t *testing.T) {
	deployments := []string{"a.b.c", "a.b", "a", "a.d", "a.d.e.f.g"}
	routes, ok := virtualServiceMatchRoute("testSvc", deployments, "a.b", "test",
		nil, ".", 0, "dev")
	route := routes[0]
	if !ok || route.Match[0].Headers["test"].Exact != "a.b" || route.Route[0].Destination.Subset != "a-b" {
		t.Errorf("Match failed 1 : " + strconv.FormatBool(ok) + ", match: " + route.Match[0].Headers["test"].Exact +
			", subset: " + route.Route[0].Destination.Subset)
	}
	routes, ok = virtualServiceMatchRoute("testSvc", deployments, "a.d.e.f", "test",
		nil, ".", 0, "dev")
	route = routes[0]
	if !ok || route.Match[0].Headers["test"].Exact != "a.d.e.f" || route.Route[0].Destination.Subset != "a-d" {
		t.Errorf("Match failed 2 : " + strconv.FormatBool(ok) + ", match: " + route.Match[0].Headers["test"].Exact +
			", subset: " + route.Route[0].Destination.Subset)
	}
	routes, ok = virtualServiceMatchRoute("testSvc", deployments, "b.x", "test",
		nil, ".", 0, "dev")
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

func TestMatchHeadersEqual(t *testing.T) {
	headers1 := map[string]v1alpha1.StringMatch{"ali-env-mark": {Exact: "dev.proj2"},
		"COOKIE": {Regex: ".*;VirtualEnv=dev.proj2;.*"}}
	headers2 := map[string]v1alpha1.StringMatch{"COOKIE": {Regex: ".*;VirtualEnv=dev.proj2;.*"},
		"ali-env-mark": {Exact: "dev.proj2"}}
	headers3 := map[string]v1alpha1.StringMatch{"ali-env-mark": {Exact: "dev.proj2"}}
	headers4 := map[string]v1alpha1.StringMatch{"COOKIE": {Regex: ".*;VirtualEnv=dev.proj3;.*"},
		"ali-env-mark": {Exact: "dev.proj2"}}
	headers5 := map[string]v1alpha1.StringMatch{"COOKIE": {Regex: ".*;VirtualEnv=dev.proj2;.*"},
		"ali-env-mark": {Exact: "dev.proj3"}}
	if !isMatchHeadersEqual(headers1, headers2) {
		t.Error("header1 and header2 should equal !")
	}
	if isMatchHeadersEqual(headers1, headers3) {
		t.Error("header1 and header3 not should equal !")
	}
	if isMatchHeadersEqual(headers1, headers4) {
		t.Error("header1 and header4 not should equal !")
	}
	if isMatchHeadersEqual(headers1, headers5) {
		t.Error("header1 and header5 not should equal !")
	}
}
