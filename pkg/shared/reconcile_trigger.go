package shared

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// guaranteed time interval between virtual environment reconcile
const reconcileCoolOffSeconds = 5

// mutex to make sure there is only one reconcile trigger candidate
var ReconcileTriggerLock = TriableMutex{}

// mechanism to reduce virtual env reconcile frequency
var ShouldDelayRefresh = AtomBool{}

// virtual env controller
var VirtualEnvController = new(controller.Controller)

// trigger virtual environment reconcile
func ReconcileVirtualEnv(namespace string, logger logr.Logger) {
	// only the first changed resource would trigger a reconcile
	if ReconcileTriggerLock.TryLock() {
		logger.Info("trigger reconcile VirtualEnvironment")
		go func() {
			ShouldDelayRefresh.Set(true)
			// reconcile triggered only after the cooling time of the last resource change event ends
			for ShouldDelayRefresh.Get() {
				ShouldDelayRefresh.Set(false)
				time.Sleep(reconcileCoolOffSeconds * time.Second)
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
		// other resource change events only delay the reconcile time
		ShouldDelayRefresh.Set(true)
	}
}
