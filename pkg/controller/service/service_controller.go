package service

import (
	"alibaba.com/virtual-env-operator/pkg/component/router"
	"alibaba.com/virtual-env-operator/pkg/controller/virtualenv"
	"alibaba.com/virtual-env-operator/pkg/shared"
	"alibaba.com/virtual-env-operator/pkg/shared/logger"
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
)

const (
	annotationGateways = "kt-virtual-environment/gateways"
	annotationHosts    = "kt-virtual-environment/hosts"
	annotationRule     = "kt-virtual-environment/rule"
)

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
	if request.Name == "kubernetes" {
		// Ignore kubernetes service
		return reconcile.Result{}, nil
	}

	shared.Lock.RLock()

	service := &corev1.Service{}
	err := r.client.Get(context.TODO(), request.NamespacedName, service)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Removing Service", request.Name)
			delete(shared.AvailableServices, request.Name)
			// delete related virtual service and destination rule
			_ = router.GetDefaultRoute().CleanupRoute(r.client, request.Namespace, request.Name)
			shared.Lock.RUnlock()
			virtualenv.TriggerReconcile()
			return reconcile.Result{}, nil
		}
		shared.Lock.RUnlock()
		return reconcile.Result{}, err
	}

	logger.Info("Adding Service", request.Name)
	serviceInfo := shared.ServiceInfo{}
	// save service selectors
	serviceInfo.Selectors = service.Spec.Selector
	// save service ports
	serviceInfo.Ports = map[string]uint32{}
	for _, port := range service.Spec.Ports {
		serviceInfo.Ports[port.Name] = uint32(port.Port)
	}
	// fetch service gateways and hosts from annotation
	if value, ok := service.Annotations[annotationGateways]; ok {
		serviceInfo.Gateways = strings.Split(value, ",")
	}
	if value, ok := service.Annotations[annotationHosts]; ok {
		serviceInfo.Hosts = strings.Split(value, ",")
	}
	if value, ok := service.Annotations[annotationRule]; ok {
		serviceInfo.CustomRule = value
	}
	shared.AvailableServices[request.Name] = serviceInfo

	shared.Lock.RUnlock()

	virtualenv.TriggerReconcile()
	return reconcile.Result{}, nil
}
