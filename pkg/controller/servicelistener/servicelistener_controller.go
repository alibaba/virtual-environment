package servicelistener

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_servicelistener")

// Add creates a new ServiceListener Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileServiceListener{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("servicelistener-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileServiceListener implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileServiceListener{}

// ReconcileServiceListener reconciles a ServiceListener object
type ReconcileServiceListener struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ServiceListener object and makes changes based on the state read
// and what is in the ServiceListener.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileServiceListener) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Namespace", request.Namespace, "Name", request.Name)

	if request.Name == "kubernetes" {
		// Ignore kubernetes service
		return reconcile.Result{}, nil
	}

	shared.Lock.RLock()

	service := &corev1.Service{}
	err := r.client.Get(context.TODO(), request.NamespacedName, service)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Removing Service")
			delete(shared.AvailableServices, request.Name)
			// delete related virtual service and destination rule
			cachedInstanceName := shared.NameWithPostfix(request.Name, shared.InsNamePostfix)
			shared.DeleteVirtualService(r.client, request.Namespace, cachedInstanceName, reqLogger)
			shared.DeleteDestinationRule(r.client, request.Namespace, cachedInstanceName, reqLogger)
			shared.Lock.RUnlock()
			return reconcile.Result{}, nil
		}
		shared.Lock.RUnlock()
		return reconcile.Result{}, err
	}

	reqLogger.Info("Adding Service")
	shared.AvailableServices[request.Name] = service.Spec.Selector

	shared.Lock.RUnlock()

	shared.ReconcileVirtualEnv(request.Namespace, reqLogger)
	return reconcile.Result{}, nil
}
