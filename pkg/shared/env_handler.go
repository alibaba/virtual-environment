package shared

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// delete specified instance
func DeleteIns(client client.Client, namespace string, name string, obj runtime.Object) error {
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, obj)
	if err != nil {
		return err
	}
	err = client.Delete(context.TODO(), obj)
	if err != nil {
		return err
	}
	return nil
}

func ReconcileVirtualEnv(namespace string, logger logr.Logger) {
	if ReconcileTriggerLock.TryLock() {
		logger.Info("trigger reconcile VirtualEnvironment")
		go func() {
			time.Sleep(3 * time.Second)
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
	}
}
