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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

const (
	// guaranteed time interval between virtual environment reconcile
	reconcileCoolOffSeconds = 5
	// default values for virtual environment configures
	defaultEnvHeader                 = "X-Virtual-Env"
	defaultEnvLabel                  = "virtual-env"
	defaultEnvSplitter               = "."
	defaultEnvHeaderAliasPlaceholder = "(@)"
)

var (
	// logger
	logger logr.Logger = log.WithName("reconcile")
	// global virtual environment object
	globalVirtualEnvironment *ReconcileVirtualEnv
	// mutex to make sure there is only one reconcile trigger candidate
	ReconcileTriggerLock = shared.TriableMutex{}
	// mechanism to reduce virtual env reconcile frequency
	ShouldDelayRefresh = shared.AtomBool{}
)

// trigger virtual environment reconcile
func TriggerReconcile() {
	// only the first changed resource would trigger a reconcile
	if ReconcileTriggerLock.TryLock() {
		logger.Info("Trigger reconcile VirtualEnvironment")
		go func() {
			ShouldDelayRefresh.Set(true)
			// reconcile triggered only after the cooling time of the last resource change event ends
			for ShouldDelayRefresh.Get() {
				ShouldDelayRefresh.Set(false)
				time.Sleep(reconcileCoolOffSeconds * time.Second)
			}
			if shared.VirtualEnvIns != nil {
				logger.Info("Send reconcile signal")
				_, err := reconcileVirtualEnvironment(shared.VirtualEnvIns)
				if err != nil {
					logger.Error(err, "failed to reconcile VirtualEnvironment")
				}
			}
			ReconcileTriggerLock.Unlock()
		}()
	} else {
		// other resource change events only delay the reconcile time
		ShouldDelayRefresh.Set(true)
	}
}

// core logic of virtual environment operator
func reconcileVirtualEnvironment(nn *types.NamespacedName) (reconcile.Result, error) {
	shared.Lock.Lock()
	c := globalVirtualEnvironment.client
	s := globalVirtualEnvironment.scheme

	virtualEnv, err := fetchVirtualEnvIns(*nn)
	if err != nil && !shared.IsVirtualEnvChanged(err) {
		shared.Lock.Unlock()
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		} else {
			return reconcile.Result{}, err
		}
	}

	logger.Info("Receive reconcile request with name: " + nn.Name)
	err = checkTagAppender(virtualEnv, *nn, shared.IsVirtualEnvChanged(err))
	for svc, service := range shared.AvailableServices {
		selector := service.Selectors
		availableLabels := parser.FindAllVirtualEnvLabelValues(shared.AvailableLabels, virtualEnv.Spec.EnvLabel.Name)
		relatedDeployments := parser.FindAllRelatedLabels(shared.AvailableLabels, selector, virtualEnv.Spec.EnvLabel.Name)
		if len(availableLabels) > 0 && len(relatedDeployments) > 0 {
			// update mesh controller panel configure
			err = router.GetDefaultRoute().GenerateRoute(c, s, virtualEnv, nn.Namespace, svc,
				availableLabels, relatedDeployments)
		}
	}

	shared.Lock.Unlock()
	return reconcile.Result{}, err
}

// create or delete tag appender according to virtual environment configure
func checkTagAppender(virtualEnv *envv1alpha2.VirtualEnvironment, nn types.NamespacedName, isVirtualEnvChanged bool) error {
	c := globalVirtualEnvironment.client
	s := globalVirtualEnvironment.scheme
	tagAppenderStatus := router.GetDefaultRoute().CheckTagAppender(c, virtualEnv, nn.Namespace, nn.Name)
	if virtualEnv.Spec.EnvHeader.AutoInject {
		if isVirtualEnvChanged || common.IsTagAppenderNeedUpdate(tagAppenderStatus) {
			_ = router.GetDefaultRoute().DeleteTagAppender(c, nn.Namespace, nn.Name)
			return router.GetDefaultRoute().CreateTagAppender(c, s, virtualEnv, nn.Namespace, nn.Name)
		}
	} else {
		if isVirtualEnvChanged || common.IsTagAppenderExist(tagAppenderStatus) {
			return router.GetDefaultRoute().DeleteTagAppender(c, nn.Namespace, nn.Name)
		}
	}
	return nil
}

// fetch the VirtualEnv instance from request
func fetchVirtualEnvIns(nn types.NamespacedName) (*envv1alpha2.VirtualEnvironment, error) {
	virtualEnv := &envv1alpha2.VirtualEnvironment{}
	err := globalVirtualEnvironment.client.Get(context.TODO(), nn, virtualEnv)
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
			deleteVirtualEnv(nn.Namespace, shared.VirtualEnvIns.Name)
		}
		shared.VirtualEnvIns = &types.NamespacedName{Namespace: nn.Namespace, Name: nn.Name}
		logger.Info("VirtualEnv added", "Spec", virtualEnv.Spec)
		return virtualEnv, shared.VirtualEnvChangeDetected{}
	}
	return virtualEnv, nil
}

// delete specified virtual env instance
func deleteVirtualEnv(namespace string, name string) {
	err := shared.DeleteIns(globalVirtualEnvironment.client, namespace, name, &envv1alpha2.VirtualEnvironment{})
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
