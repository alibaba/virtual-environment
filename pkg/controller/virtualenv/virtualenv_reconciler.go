package virtualenv

import (
	envv1alpha2 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha2"
	"alibaba.com/virtual-env-operator/pkg/component/parser"
	"alibaba.com/virtual-env-operator/pkg/component/router"
	"alibaba.com/virtual-env-operator/pkg/component/router/common"
	"alibaba.com/virtual-env-operator/pkg/shared"
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const defaultEnvHeader = "X-Virtual-Env"
const defaultEnvLabel = "virtual-env"
const defaultEnvSplitter = "."
const defaultEnvHeaderAliasPlaceholder = "(@)"

// core logic of virtual environment operator
func ReconcileVirtualEnvironment(client client.Client, scheme *runtime.Scheme,
	nn types.NamespacedName, logger logr.Logger) (reconcile.Result, error) {
	shared.Lock.Lock()

	virtualEnv, err := fetchVirtualEnvIns(client, nn, logger)
	if err != nil && !shared.IsVirtualEnvChanged(err) {
		shared.Lock.Unlock()
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		} else {
			return reconcile.Result{}, err
		}
	}

	logger.Info("Receive reconcile request with name: " + nn.Name)
	err = checkTagAppender(client, scheme, virtualEnv, nn, shared.IsVirtualEnvChanged(err))
	for svc, service := range shared.AvailableServices {
		selector := service.Selectors
		availableLabels := parser.FindAllVirtualEnvLabelValues(shared.AvailableLabels, virtualEnv.Spec.EnvLabel.Name)
		relatedDeployments := parser.FindAllRelatedLabels(shared.AvailableLabels, selector, virtualEnv.Spec.EnvLabel.Name)
		if len(availableLabels) > 0 && len(relatedDeployments) > 0 {
			// update mesh controller panel configure
			err = router.GetDefaultRoute().GenerateRoute(client, scheme, virtualEnv, nn.Namespace, svc,
				availableLabels, relatedDeployments)
		}
	}

	shared.Lock.Unlock()
	return reconcile.Result{}, err
}

// create or delete tag appender according to virtual environment configure
func checkTagAppender(client client.Client, scheme *runtime.Scheme, virtualEnv *envv1alpha2.VirtualEnvironment,
	nn types.NamespacedName, isVirtualEnvChanged bool) error {
	tagAppenderStatus := router.GetDefaultRoute().CheckTagAppender(client, virtualEnv, nn.Namespace, nn.Name)
	if virtualEnv.Spec.EnvHeader.AutoInject {
		if isVirtualEnvChanged || common.IsTagAppenderNeedUpdate(tagAppenderStatus) {
			_ = router.GetDefaultRoute().DeleteTagAppender(client, nn.Namespace, nn.Name)
			return router.GetDefaultRoute().CreateTagAppender(client, scheme, virtualEnv, nn.Namespace, nn.Name)
		}
	} else {
		if isVirtualEnvChanged || common.IsTagAppenderExist(tagAppenderStatus) {
			return router.GetDefaultRoute().DeleteTagAppender(client, nn.Namespace, nn.Name)
		}
	}
	return nil
}

// fetch the VirtualEnv instance from request
func fetchVirtualEnvIns(client client.Client, nn types.NamespacedName,
	logger logr.Logger) (*envv1alpha2.VirtualEnvironment, error) {
	virtualEnv := &envv1alpha2.VirtualEnvironment{}
	err := client.Get(context.TODO(), nn, virtualEnv)
	if err != nil {
		if errors.IsNotFound(err) {
			// virtual environment removed or haven't created yet
			logger.Info("VirtualEnv resource missing")
			if shared.VirtualEnvIns != nil && shared.VirtualEnvIns.Name == nn.Name {
				shared.VirtualEnvIns = nil
				logger.Info("VirtualEnv record removed")
			}
			return nil, err
		}
		logger.Error(err, "Failed to get VirtualEnvironment")
		return nil, err
	}
	handleDefaultConfig(virtualEnv)
	if shared.VirtualEnvIns == nil || shared.VirtualEnvIns.Name != nn.Name {
		// new virtual environment found
		if shared.VirtualEnvIns != nil {
			// there is an old virtual environment exist
			logger.Info("New VirtualEnv resource detected, deleting " + shared.VirtualEnvIns.Name)
			deleteVirtualEnv(client, nn.Namespace, shared.VirtualEnvIns.Name, logger)
		}
		shared.VirtualEnvIns = &types.NamespacedName{Namespace: nn.Namespace, Name: nn.Name}
		logger.Info("VirtualEnv added", "Spec", virtualEnv.Spec)
		return virtualEnv, shared.VirtualEnvChangeDetected{}
	}
	return virtualEnv, nil
}

// delete specified virtual env instance
func deleteVirtualEnv(client client.Client, namespace string, name string, logger logr.Logger) {
	err := shared.DeleteIns(client, namespace, name, &envv1alpha2.VirtualEnvironment{})
	if err != nil {
		logger.Error(err, "Failed to remove VirtualEnv instance "+name)
	} else {
		logger.Info("VirtualEnv deleted")
	}
}

// handle empty virtual env configure item with default value
func handleDefaultConfig(virtualEnv *envv1alpha2.VirtualEnvironment) {
	if virtualEnv.Spec.EnvHeader.Name == "" {
		virtualEnv.Spec.EnvHeader.Name = defaultEnvHeader
	}
	if virtualEnv.Spec.EnvHeader.Aliases != nil {
		for _, alias := range virtualEnv.Spec.EnvHeader.Aliases {
			if alias.Placeholder == "" {
				alias.Placeholder = defaultEnvHeaderAliasPlaceholder
			}
		}
	}
	if virtualEnv.Spec.EnvLabel.Name == "" {
		virtualEnv.Spec.EnvLabel.Name = defaultEnvLabel
	}
	if virtualEnv.Spec.EnvLabel.Splitter == "" {
		virtualEnv.Spec.EnvLabel.Splitter = defaultEnvSplitter
	}
}
