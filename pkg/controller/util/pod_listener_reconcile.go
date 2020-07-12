package util

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var log = logf.Log.WithName("controller_podslistener")

func Reconcile(client client.Client, request reconcile.Request, obj runtime.Object,
	getLabels func(interface{}) map[string]string) (reconcile.Result, error) {

	resourceType := obj.GetObjectKind().GroupVersionKind().Kind
	reqLogger := log.WithValues("Ref", request.Namespace+":"+resourceType+":"+request.Name)
	shared.Lock.RLock()

	err := client.Get(context.TODO(), request.NamespacedName, obj)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Removing " + resourceType)
			delete(shared.AvailableLabels, labelMark(resourceType, request.Name))
			shared.Lock.RUnlock()
			shared.ReconcileVirtualEnv(request.Namespace, reqLogger)
			return reconcile.Result{}, nil
		}
		shared.Lock.RUnlock()
		return reconcile.Result{}, err
	}

	reqLogger.Info("Adding " + resourceType)
	shared.AvailableLabels[labelMark(resourceType, request.Name)] = getLabels(obj)

	shared.Lock.RUnlock()

	shared.ReconcileVirtualEnv(request.Namespace, reqLogger)
	return reconcile.Result{}, nil
}

func labelMark(resourceType string, name string) string {
	return resourceType + "#" + name
}
