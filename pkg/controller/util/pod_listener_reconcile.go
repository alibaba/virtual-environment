package util

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var log = logf.Log.WithName("controller_podslistener")

func Reconcile(client client.Client, request reconcile.Request, resourceType string) (reconcile.Result, error) {
	reqLogger := log.WithValues("Ref", request.Namespace+":"+resourceType+":"+request.Name)
	shared.Lock.RLock()

	deployment := &appsv1.Deployment{}
	err := client.Get(context.TODO(), request.NamespacedName, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Removing Deployment")
			delete(shared.AvailableDeployments, request.Name)
			shared.Lock.RUnlock()
			shared.ReconcileVirtualEnv(request.Namespace, reqLogger)
			return reconcile.Result{}, nil
		}
		shared.Lock.RUnlock()
		return reconcile.Result{}, err
	}

	reqLogger.Info("Adding Deployment")
	shared.AvailableDeployments[request.Name] = deployment.Spec.Template.Labels

	shared.Lock.RUnlock()

	shared.ReconcileVirtualEnv(request.Namespace, reqLogger)
	return reconcile.Result{}, nil
}
