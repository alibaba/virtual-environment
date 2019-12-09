package virtualenv

import (
	envv1alpha1 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha1"
	"alibaba.com/virtual-env-operator/pkg/envoy"
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

const defaultEnvHeader = "X-Virtual-Env"
const defaultEnvLabel = "virtual-env"
const defaultEnvSplitter = "."

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
	err = c.Watch(&source.Kind{Type: &envv1alpha1.VirtualEnvironment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource VirtualService & DestinationRule, requeue their owner to VirtualEnv
	err = c.Watch(&source.Kind{Type: &networkingv1alpha3.VirtualService{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &envv1alpha1.VirtualEnvironment{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &networkingv1alpha3.DestinationRule{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &envv1alpha1.VirtualEnvironment{},
	})
	if err != nil {
		return err
	}

	shared.VirtualEnvController = &c
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

	reqLogger.Info("Reconciling VirtualEnvironment")
	for svc, selector := range shared.AvailableServices {
		availableLabels := shared.FindAllVirtualEnvLabelValues(shared.AvailableDeployments, virtualEnv.Spec.EnvLabel.Name)
		relatedDeployments := shared.FindAllRelatedDeployments(shared.AvailableDeployments, selector, virtualEnv.Spec.EnvLabel.Name)

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
func (r *ReconcileVirtualEnv) fetchVirtualEnvIns(request reconcile.Request, logger logr.Logger) (*envv1alpha1.VirtualEnvironment, error) {
	virtualEnv := &envv1alpha1.VirtualEnvironment{}
	err := r.client.Get(context.TODO(), request.NamespacedName, virtualEnv)
	if err != nil {
		if errors.IsNotFound(err) {
			// virtual environment removed or haven't created yet
			logger.Info("VirtualEnv resource missing")
			if shared.VirtualEnvIns == request.Name {
				shared.VirtualEnvIns = ""
				logger.Info("VirtualEnv record removed")
			}
			return nil, err
		}
		logger.Error(err, "Failed to get VirtualEnvironment")
		return nil, err
	}
	r.handleDefaultConfig(virtualEnv)
	if shared.VirtualEnvIns != request.Name {
		// new virtual environment found
		if shared.VirtualEnvIns != "" {
			// there is an old virtual environment exist
			logger.Info("New VirtualEnv resource detected, deleting " + shared.VirtualEnvIns)
			r.deleteVirtualEnv(request.Namespace, shared.VirtualEnvIns, logger)
		}
		shared.VirtualEnvIns = request.Name
		if virtualEnv.Spec.EnvHeader.AutoInject {
			tagAppenderName := shared.NameWithPostfix(request.Name)
			err = r.createTagAppender(request.Namespace, tagAppenderName, virtualEnv, logger)
			if err != nil {
				logger.Error(err, "failed to create TagAppender instance "+tagAppenderName)
				return virtualEnv, err
			}
		}
		logger.Info("VirtualEnv added", "Spec", virtualEnv.Spec)
	}
	return virtualEnv, err
}

// delete specified virtual env instance
func (r *ReconcileVirtualEnv) deleteVirtualEnv(namespace string, name string, logger logr.Logger) {
	err := shared.DeleteIns(r.client, namespace, name, &envv1alpha1.VirtualEnvironment{})
	if err != nil {
		logger.Error(err, "Failed to remove VirtualEnv instance "+name)
	} else {
		logger.Info("VirtualEnv deleted")
	}
}

// create tag auto appender filter instance
func (r *ReconcileVirtualEnv) createTagAppender(namespace string, name string, virtualEnv *envv1alpha1.VirtualEnvironment,
	logger logr.Logger) error {
	err := envoy.DeleteTagAppenderIfExist(r.client, namespace, name)
	if err != nil {
		logger.Error(err, "Failed to remove old TagAppender instance")
		return err
	}
	tagAppender := envoy.TagAppenderFilter(namespace, name, virtualEnv.Spec.EnvLabel.Name, virtualEnv.Spec.EnvHeader.Name)
	// set VirtualEnv instance as the owner and controller
	err = controllerutil.SetControllerReference(virtualEnv, tagAppender, r.scheme)
	if err == nil {
		err = r.client.Create(context.TODO(), tagAppender)
		if err == nil {
			logger.Info("TagAppender created")
		}
	}
	return err
}

// handle empty virtual env configure item with default value
func (r *ReconcileVirtualEnv) handleDefaultConfig(virtualEnv *envv1alpha1.VirtualEnvironment) {
	if virtualEnv.Spec.EnvHeader.Name == "" {
		virtualEnv.Spec.EnvHeader.Name = defaultEnvHeader
	}
	if virtualEnv.Spec.EnvLabel.Name == "" {
		virtualEnv.Spec.EnvLabel.Name = defaultEnvLabel
	}
	if virtualEnv.Spec.EnvLabel.Splitter == "" {
		virtualEnv.Spec.EnvLabel.Splitter = defaultEnvSplitter
	}
	shared.InsNamePostfix = virtualEnv.Spec.InstancePostfix
}

// reconcile virtual service according to related deployments and available labels
func (r *ReconcileVirtualEnv) reconcileVirtualService(virtualEnv *envv1alpha1.VirtualEnvironment, svc string, request reconcile.Request,
	availableLabels []string, relatedDeployments map[string]string, logger logr.Logger) error {
	virtualSvc := shared.VirtualService(request.Namespace, svc, availableLabels, relatedDeployments,
		virtualEnv.Spec.EnvHeader.Name, virtualEnv.Spec.EnvLabel.Splitter, virtualEnv.Spec.EnvLabel.DefaultSubset)
	foundVirtualSvc := &networkingv1alpha3.VirtualService{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: shared.NameWithPostfix(svc), Namespace: request.Namespace}, foundVirtualSvc)
	if err != nil {
		// VirtualService not exist, create one
		if errors.IsNotFound(err) {
			// set VirtualEnv instance as the owner and controller
			err = controllerutil.SetControllerReference(virtualEnv, virtualSvc, r.scheme)
			if err != nil {
				logger.Error(err, "Failed to set owner of "+virtualSvc.Name)
				return err
			}
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
	} else if shared.IsDifferentVirtualService(&foundVirtualSvc.Spec, &virtualSvc.Spec, virtualEnv.Spec.EnvHeader.Name) {
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
func (r *ReconcileVirtualEnv) reconcileDestinationRule(virtualEnv *envv1alpha1.VirtualEnvironment, svc string,
	request reconcile.Request, relatedDeployments map[string]string, logger logr.Logger) error {
	destRule := shared.DestinationRule(request.Namespace, svc, relatedDeployments, virtualEnv.Spec.EnvLabel.Name)
	foundDestRule := &networkingv1alpha3.DestinationRule{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: shared.NameWithPostfix(svc), Namespace: request.Namespace}, foundDestRule)
	if err != nil {
		// DestinationRule not exist, create one
		if errors.IsNotFound(err) {
			// set VirtualEnv instance as the owner and controller
			err = controllerutil.SetControllerReference(virtualEnv, destRule, r.scheme)
			if err != nil {
				logger.Error(err, "Failed to set owner of "+destRule.Name)
				return err
			}
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
	} else if shared.IsDifferentDestinationRule(&foundDestRule.Spec, &destRule.Spec, virtualEnv.Spec.EnvLabel.Name) {
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
