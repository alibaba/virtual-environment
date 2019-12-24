package router

import (
	"alibaba.com/virtual-env-operator/pkg/component/router/istio"
)

var pool = make(map[string]IsolationRouter)

func init() {
	pool["IstioHttp"] = &istio.HttpRouter{}
}

func GetRoute(name string) IsolationRouter {
	return pool[name]
}

// This method will be removed in future
func GetDefaultRoute() IsolationRouter {
	return GetRoute("IstioHttp")
}
