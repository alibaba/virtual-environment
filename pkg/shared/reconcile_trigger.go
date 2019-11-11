package shared

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// mutex to reduce virtual env reconcile frequency cause by deployment/service change
var ReconcileTriggerLock = TriableMutex{}

// another mechanism to reduce virtual env reconcile frequency
var ShouldDelayRefresh = AtomBool{}

// virtual env controller
var VirtualEnvController = new(controller.Controller)

// trigger virtual environment reconcile
func ReconcileVirtualEnv(namespace string, logger logr.Logger) {
	if ReconcileTriggerLock.TryLock() {
		logger.Info("trigger reconcile VirtualEnvironment")
		go func() {
			ShouldDelayRefresh.Set(true)
			for ShouldDelayRefresh.Get() {
				ShouldDelayRefresh.Set(false)
				time.Sleep(3 * time.Second)
			}
			if VirtualEnvIns != "" {
				_, err := (*VirtualEnvController).Reconcile(reconcile.Request{
					NamespacedName: types.NamespacedName{Name: VirtualEnvIns, Namespace: namespace},
				})
				if err != nil {
					logger.Error(err, "failed to reconcile VirtualEnvironment")
				}
			}
			ReconcileTriggerLock.Unlock()
		}()
	} else {
		ShouldDelayRefresh.Set(true)
	}
}
