package deploymentlistener

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	"context"
	appsv1 "k8s.io/api/apps/v1"
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

var log = logf.Log.WithName("controller_deploymentlistener")

// Add creates a new DeploymentListener Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDeploymentListener{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("deploymentlistener-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Deployment
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDeploymentListener implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDeploymentListener{}

// ReconcileDeploymentListener reconciles a DeploymentListener object
type ReconcileDeploymentListener struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a DeploymentListener object and makes changes based on the state read
// and what is in the DeploymentListener.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDeploymentListener) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Namespace", request.Namespace, "Name", request.Name)

	shared.Lock.RLock()

	deployment := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), request.NamespacedName, deployment)
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
