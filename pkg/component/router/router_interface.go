package router

import (
	envv1alpha1 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

type IsolationRouter interface {
	GenerateRoute(client client.Client, scheme *runtime.Scheme, virtualEnv *envv1alpha1.VirtualEnvironment,
		namespace string, svcName string, availableLabels []string, relatedDeployments map[string]string) error

	CleanupRoute(client client.Client, namespace string, name string) error

	RegisterReconcileWatcher(c controller.Controller) error

	DeleteTagAppender(client client.Client, namespace string, name string) error

	CreateTagAppender(client client.Client, scheme *runtime.Scheme, virtualEnv *envv1alpha1.VirtualEnvironment,
		namespace string, name string) error
}
