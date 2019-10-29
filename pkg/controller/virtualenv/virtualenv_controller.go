package virtualenv

import (
	envv1alpha1 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha1"
	"alibaba.com/virtual-env-operator/pkg/status"
	"context"
	networkingv1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_virtualenv")

// Add creates a new VirtualEnv Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileVirtualEnv{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("virtualenv-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource VirtualEnv
	err = c.Watch(&source.Kind{Type: &envv1alpha1.VirtualEnv{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource VirtualService & DestinationRule, requeue their owner to VirtualEnv
	err = c.Watch(&source.Kind{Type: &networkingv1alpha3.VirtualService{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &envv1alpha1.VirtualEnv{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &networkingv1alpha3.DestinationRule{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &envv1alpha1.VirtualEnv{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileVirtualEnv implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileVirtualEnv{}

// ReconcileVirtualEnv reconciles a VirtualEnv object
type ReconcileVirtualEnv struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a VirtualEnv object and makes changes based on the state read
// and what is in the VirtualEnv.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileVirtualEnv) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Namespace", request.Namespace, "Name", request.Name)
	reqLogger.Info("Reconciling VirtualEnv")

	status.Lock.Lock()

	// Fetch the VirtualEnv instance
	virtualEnv := &envv1alpha1.VirtualEnv{}
	err := r.client.Get(context.TODO(), request.NamespacedName, virtualEnv)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("VirtualEnv resource missing")
			if status.VirtualEnvIns == request.Name {
				reqLogger.Info("VirtualEnv deleted")
				status.VirtualEnvIns = ""
			}
			status.Lock.Unlock()
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get VirtualEnv")
		status.Lock.Unlock()
		return reconcile.Result{}, err
	}

	if status.VirtualEnvIns != request.Name {
		if status.VirtualEnvIns != "" {
			reqLogger.Info("New VirtualEnv resource detected, deleting " + status.VirtualEnvIns)
			deleteVirtualEnv(request.Namespace, status.VirtualEnvIns)
		}
		reqLogger.Info("VirtualEnv added", "VeLabel", virtualEnv.Spec.VeLabel,
			"VeHeader", virtualEnv.Spec.VeHeader, "VeSplitter", virtualEnv.Spec.VeSplitter)
		status.VirtualEnvIns = request.Name
	}
	reqLogger.Info("Responding VirtualEnv")

	for srv, selector := range status.AvailableServices {
		virtualSrv := &networkingv1alpha3.VirtualService{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: srv, Namespace: request.Namespace}, virtualSrv)
		if err != nil {
			// VirtualService not exist, create one
			if errors.IsNotFound(err) {
				virtualSrv = r.virtualService(virtualEnv, srv, request.Namespace, selector)
				reqLogger.Info("Creating VirtualService " + virtualSrv.Name)
				err = r.client.Create(context.TODO(), virtualSrv)
				if err != nil {
					reqLogger.Error(err, "Failed to create VirtualService "+virtualSrv.Name)
					status.Lock.Unlock()
					return reconcile.Result{}, err
				}
			} else {
				reqLogger.Error(err, "Failed to get VirtualService")
				status.Lock.Unlock()
				return reconcile.Result{}, err
			}
		} else {
			// VirtualService already exist, TODO: check and update
			//err := r.client.Status().Update(context.TODO(), virtualSrv)
			//if err != nil {
			//	reqLogger.Error(err, "Failed to update VirtualService status")
			//	status.Lock.Unlock()
			//	return reconcile.Result{}, err
			//}
		}

		destRule := &networkingv1alpha3.DestinationRule{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: srv, Namespace: request.Namespace}, destRule)
		if err != nil {
			// DestinationRule not exist, create one
			if errors.IsNotFound(err) {
				destRule = r.destinationRule(virtualEnv, srv, request.Namespace, selector)
				reqLogger.Info("Creating DestinationRule " + destRule.Name)
				err = r.client.Create(context.TODO(), destRule)
				if err != nil {
					reqLogger.Error(err, "Failed to create DestinationRule "+destRule.Name)
					status.Lock.Unlock()
					return reconcile.Result{}, err
				}
			} else {
				reqLogger.Error(err, "Failed to get DestinationRule")
				status.Lock.Unlock()
				return reconcile.Result{}, err
			}
		} else {
			// DestinationRule already exist, TODO: check and update
			//err := r.client.Status().Update(context.TODO(), destRule)
			//if err != nil {
			//	reqLogger.Error(err, "Failed to update DestinationRule status")
			//	status.Lock.Unlock()
			//	return reconcile.Result{}, err
			//}
		}
	}

	status.Lock.Unlock()
	return reconcile.Result{}, nil
}

func (r *ReconcileVirtualEnv) virtualService(e *envv1alpha1.VirtualEnv, name string, namespace string,
	selector map[string]string) *networkingv1alpha3.VirtualService {
	virtualSrv := &networkingv1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.VirtualServiceSpec{
			Hosts: []string{name},
			HTTP: []networkingv1alpha3.HTTPRoute{{
				Route: []networkingv1alpha3.HTTPRouteDestination{{
					Destination: networkingv1alpha3.Destination{
						Host:   name,
						Subset: "default",
					},
				}},
			}},
		},
	}
	// Set VirtualEnv instance as the owner and controller
	controllerutil.SetControllerReference(e, virtualSrv, r.scheme)
	return virtualSrv
}

func (r *ReconcileVirtualEnv) destinationRule(e *envv1alpha1.VirtualEnv, name string, namespace string,
	selector map[string]string) *networkingv1alpha3.DestinationRule {
	destRule := &networkingv1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.DestinationRuleSpec{
			Host: name,
			Subsets: []networkingv1alpha3.Subset{{
				Name:   "default",
				Labels: map[string]string{e.Spec.VeLabel: "default"},
			}},
		},
	}
	// Set VirtualEnv instance as the owner and controller
	controllerutil.SetControllerReference(e, destRule, r.scheme)
	return destRule
}

func deleteVirtualEnv(namespace string, virtualEnv string) {
	//TODO: delete virtual env instance
}
