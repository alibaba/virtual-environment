package virtualenv

import (
	envv1alpha1 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha1"
	"alibaba.com/virtual-env-operator/pkg/shared"
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	networkingv1alpha3 "knative.dev/pkg/apis/istio/v1alpha3"
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

	shared.Lock.Lock()

	virtualEnv, err := r.fetchVirtualEnvIns(request, reqLogger)
	if err != nil {
		shared.Lock.Unlock()
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		} else {
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Responding VirtualService and DestinationRule")
	for svc, selector := range shared.AvailableServices {
		availableLabels := shared.FindAllVirtualEnvLabelValues(shared.AvailableDeployments, virtualEnv.Spec.VeLabel)
		relatedDeployments := shared.FindAllRelatedDeployments(shared.AvailableDeployments, selector, virtualEnv.Spec.VeLabel)

		if len(availableLabels) > 0 && len(relatedDeployments) > 0 {
			err = r.reconcileVirtualService(virtualEnv, svc, request, availableLabels, relatedDeployments, reqLogger)
			if err != nil {
				shared.Lock.Unlock()
				return reconcile.Result{}, err
			}
			err = r.reconcileDestinationRule(virtualEnv, svc, request, relatedDeployments, reqLogger)
			if err != nil {
				shared.Lock.Unlock()
				return reconcile.Result{}, err
			}
		}
	}

	shared.Lock.Unlock()
	return reconcile.Result{}, nil
}

// fetch the VirtualEnv instance from request
func (r *ReconcileVirtualEnv) fetchVirtualEnvIns(request reconcile.Request, logger logr.Logger) (*envv1alpha1.VirtualEnv, error) {
	virtualEnv := &envv1alpha1.VirtualEnv{}
	err := r.client.Get(context.TODO(), request.NamespacedName, virtualEnv)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("VirtualEnv resource missing")
			if shared.VirtualEnvIns == request.Name {
				shared.VirtualEnvIns = ""
				logger.Info("VirtualEnv record removed")
			}
			return nil, err
		}
		logger.Error(err, "Failed to get VirtualEnv")
		return nil, err
	}
	if shared.VirtualEnvIns != request.Name {
		if shared.VirtualEnvIns != "" {
			logger.Info("New VirtualEnv resource detected, deleting " + shared.VirtualEnvIns)
			r.deleteVirtualEnv(request.Namespace, shared.VirtualEnvIns, logger)
		}
		shared.VirtualEnvIns = request.Name
		logger.Info("VirtualEnv added", "VeLabel", virtualEnv.Spec.VeLabel,
			"VeHeader", virtualEnv.Spec.VeHeader, "VeSplitter", virtualEnv.Spec.VeSplitter)
	}
	return virtualEnv, err
}

// delete specified virtual env instance
func (r *ReconcileVirtualEnv) deleteVirtualEnv(namespace string, name string, logger logr.Logger) {
	err := shared.DeleteIns(r.client, namespace, name, &envv1alpha1.VirtualEnv{})
	if err != nil {
		logger.Error(err, "failed to remove VirtualEnv instance "+name)
	} else {
		logger.Info("VirtualEnv deleted")
	}
}

// reconcile virtual service according to related deployments and available labels
func (r *ReconcileVirtualEnv) reconcileVirtualService(virtualEnv *envv1alpha1.VirtualEnv, svc string, request reconcile.Request,
	availableLabels []string, relatedDeployments map[string]string, logger logr.Logger) error {
	virtualSvc := shared.VirtualService(svc, request.Namespace, availableLabels, relatedDeployments,
		virtualEnv.Spec.VeHeader, virtualEnv.Spec.VeSplitter)
	// Set VirtualEnv instance as the owner and controller
	err := controllerutil.SetControllerReference(virtualEnv, virtualSvc, r.scheme)
	if err != nil {
		return err
	}
	foundVirtualSvc := &networkingv1alpha3.VirtualService{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc, Namespace: request.Namespace}, foundVirtualSvc)
	if err != nil {
		// VirtualService not exist, create one
		if errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), virtualSvc)
			if err != nil {
				logger.Error(err, "Failed to create VirtualService "+virtualSvc.Name)
				return err
			}
			logger.Info("VirtualService " + virtualSvc.Name + " created")
		} else {
			logger.Error(err, "Failed to get VirtualService")
			return err
		}
	} else if shared.IsDifferentVirtualService(foundVirtualSvc.Spec, virtualSvc.Spec, virtualEnv.Spec.VeHeader) {
		// existing VirtualService changed
		foundVirtualSvc.Spec = virtualSvc.Spec
		err := r.client.Update(context.TODO(), foundVirtualSvc)
		if err != nil {
			logger.Error(err, "Failed to update VirtualService status")
			return err
		}
		logger.Info("VirtualService " + virtualSvc.Name + " changed")
	}
	return nil
}

// reconcile destination rule according to related deployments
func (r *ReconcileVirtualEnv) reconcileDestinationRule(virtualEnv *envv1alpha1.VirtualEnv, svc string, request reconcile.Request,
	relatedDeployments map[string]string, logger logr.Logger) error {
	destRule := shared.DestinationRule(svc, request.Namespace, relatedDeployments, virtualEnv.Spec.VeLabel)
	// Set VirtualEnv instance as the owner and controller
	err := controllerutil.SetControllerReference(virtualEnv, destRule, r.scheme)
	if err != nil {
		return err
	}
	foundDestRule := &networkingv1alpha3.DestinationRule{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc, Namespace: request.Namespace}, foundDestRule)
	if err != nil {
		// DestinationRule not exist, create one
		if errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), destRule)
			if err != nil {
				logger.Error(err, "Failed to create DestinationRule "+destRule.Name)
				return err
			}
			logger.Info("DestinationRule " + destRule.Name + " created")
		} else {
			logger.Error(err, "Failed to get DestinationRule")
			return err
		}
	} else if shared.IsDifferentDestinationRule(foundDestRule.Spec, destRule.Spec, virtualEnv.Spec.VeLabel) {
		// existing DestinationRule changed
		foundDestRule.Spec = destRule.Spec
		err := r.client.Update(context.TODO(), foundDestRule)
		if err != nil {
			logger.Error(err, "Failed to update DestinationRule status")
			return err
		}
		logger.Info("DestinationRule " + destRule.Name + " changed")
	}
	return nil
}
