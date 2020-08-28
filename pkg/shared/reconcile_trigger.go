package shared

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var log = logf.Log.WithName("reconcile-trigger")

const (
	// guaranteed time interval between virtual environment reconcile
	reconcileCoolOffSeconds = 5
	// a special namespace name means to execute reconcile
	ExecuteReconcileSignal = "VIRTUAL_ENVIRONMENT_RECONCILE_SIGNAL"
)

var (
	// mutex to make sure there is only one reconcile trigger candidate
	ReconcileTriggerLock = TriableMutex{}
	// mechanism to reduce virtual env reconcile frequency
	ShouldDelayRefresh = AtomBool{}
	// virtual env controller
	VirtualEnvController = new(controller.Controller)
)

// trigger virtual environment reconcile
func TriggerReconcile(referenceName string) {
	logger := log.WithValues("Ref", referenceName)
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
					NamespacedName: types.NamespacedName{Name: VirtualEnvIns, Namespace: ExecuteReconcileSignal},
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
