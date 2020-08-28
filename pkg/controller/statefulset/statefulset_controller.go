package statefulset

import (
	"alibaba.com/virtual-env-operator/pkg/controller/util"
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

var log = logf.Log.WithName("statefulset-listener")

// Add creates a new StatefulSet Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStatefulSet{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("statefulset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to StatefulSet
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileStatefulSet implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileStatefulSet{}

// ReconcileStatefulSet reconciles a StatefulSet object
type ReconcileStatefulSet struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a StatefulSet object and makes changes based on the state read
// and what is in the StatefulSet.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileStatefulSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Ref", "[StatefulSet]"+request.Name)

	shared.Lock.RLock()

	statefulset := &appsv1.StatefulSet{}
	err := r.client.Get(context.TODO(), request.NamespacedName, statefulset)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Removing StatefulSet")
			delete(shared.AvailableLabels, util.LabelMark("StatefulSet", request.Name))
			shared.Lock.RUnlock()
			shared.TriggerReconcile("[-StatefulSet]" + request.Name)
			return reconcile.Result{}, nil
		}
		shared.Lock.RUnlock()
		return reconcile.Result{}, err
	}

	logger.Info("Adding StatefulSet")
	shared.AvailableLabels[util.LabelMark("StatefulSet", request.Name)] = statefulset.Spec.Template.Labels

	shared.Lock.RUnlock()

	shared.TriggerReconcile("[+StatefulSet]" + request.Name)
	return reconcile.Result{}, nil
}
