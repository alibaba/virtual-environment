package virtualenv

import (
	envv1alpha2 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha2"
	"alibaba.com/virtual-env-operator/pkg/component/parser"
	"alibaba.com/virtual-env-operator/pkg/component/router"
	"alibaba.com/virtual-env-operator/pkg/component/router/common"
	"alibaba.com/virtual-env-operator/pkg/shared"
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("virtualenv-listener")

const defaultEnvHeader = "X-Virtual-Env"
const defaultEnvLabel = "virtual-env"
const defaultEnvSplitter = "."
const defaultEnvHeaderAliasPlaceholder = "(@)"

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
	err = c.Watch(&source.Kind{Type: &envv1alpha2.VirtualEnvironment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to generated resource
	err = router.GetDefaultRoute().RegisterReconcileWatcher(c)
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
	logger := log.WithValues("Ref", "[VirtualEnv]"+request.Name)

	// if not invoked by reconcile signal, just send a trigger
	if shared.VirtualEnvIns == nil {
		r.internalReconcile(request.NamespacedName, logger)
	} else if request.Namespace == shared.ExecuteReconcileSignal {
		r.internalReconcile(*shared.VirtualEnvIns, logger)
	}

	logger.Info("Push reconcile request to trigger")
	shared.TriggerReconcile("[VirtualEnv]" + request.Name)
	return reconcile.Result{}, nil
}

// core logic of virtual environment operator
func (r *ReconcileVirtualEnv) internalReconcile(nn types.NamespacedName, logger logr.Logger) (reconcile.Result, error) {
	shared.Lock.Lock()

	virtualEnv, err := r.fetchVirtualEnvIns(nn, logger)
	if err != nil && !shared.IsVirtualEnvChanged(err) {
		shared.Lock.Unlock()
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		} else {
			return reconcile.Result{}, err
		}
	}

	logger.Info("Receive reconcile request with name: " + nn.Name)
	err = r.checkTagAppender(virtualEnv, nn, shared.IsVirtualEnvChanged(err))
	for svc, service := range shared.AvailableServices {
		selector := service.Selectors
		availableLabels := parser.FindAllVirtualEnvLabelValues(shared.AvailableLabels, virtualEnv.Spec.EnvLabel.Name)
		relatedDeployments := parser.FindAllRelatedLabels(shared.AvailableLabels, selector, virtualEnv.Spec.EnvLabel.Name)
		if len(availableLabels) > 0 && len(relatedDeployments) > 0 {
			// update mesh controller panel configure
			err = router.GetDefaultRoute().GenerateRoute(r.client, r.scheme, virtualEnv, nn.Namespace, svc,
				availableLabels, relatedDeployments)
		}
	}

	shared.Lock.Unlock()
	return reconcile.Result{}, err
}

// create or delete tag appender according to virtual environment configure
func (r *ReconcileVirtualEnv) checkTagAppender(virtualEnv *envv1alpha2.VirtualEnvironment, nn types.NamespacedName,
	isVirtualEnvChanged bool) error {
	tagAppenderStatus := router.GetDefaultRoute().CheckTagAppender(r.client, virtualEnv, nn.Namespace, nn.Name)
	if virtualEnv.Spec.EnvHeader.AutoInject {
		if isVirtualEnvChanged || common.IsTagAppenderNeedUpdate(tagAppenderStatus) {
			_ = router.GetDefaultRoute().DeleteTagAppender(r.client, nn.Namespace, nn.Name)
			return router.GetDefaultRoute().CreateTagAppender(r.client, r.scheme, virtualEnv, nn.Namespace, nn.Name)
		}
	} else {
		if isVirtualEnvChanged || common.IsTagAppenderExist(tagAppenderStatus) {
			return router.GetDefaultRoute().DeleteTagAppender(r.client, nn.Namespace, nn.Name)
		}
	}
	return nil
}

// fetch the VirtualEnv instance from request
func (r *ReconcileVirtualEnv) fetchVirtualEnvIns(nn types.NamespacedName, logger logr.Logger) (*envv1alpha2.VirtualEnvironment, error) {
	virtualEnv := &envv1alpha2.VirtualEnvironment{}
	err := r.client.Get(context.TODO(), nn, virtualEnv)
	if err != nil {
		if errors.IsNotFound(err) {
			// virtual environment removed or haven't created yet
			logger.Info("VirtualEnv resource missing")
			if shared.VirtualEnvIns != nil && shared.VirtualEnvIns.Name == nn.Name {
				shared.VirtualEnvIns = nil
				logger.Info("VirtualEnv record removed")
			}
			return nil, err
		}
		logger.Error(err, "Failed to get VirtualEnvironment")
		return nil, err
	}
	r.handleDefaultConfig(virtualEnv)
	if shared.VirtualEnvIns == nil || shared.VirtualEnvIns.Name != nn.Name {
		// new virtual environment found
		if shared.VirtualEnvIns != nil {
			// there is an old virtual environment exist
			logger.Info("New VirtualEnv resource detected, deleting " + shared.VirtualEnvIns.Name)
			r.deleteVirtualEnv(nn.Namespace, shared.VirtualEnvIns.Name, logger)
		}
		shared.VirtualEnvIns = &types.NamespacedName{Namespace: nn.Namespace, Name: nn.Name}
		logger.Info("VirtualEnv added", "Spec", virtualEnv.Spec)
		return virtualEnv, shared.VirtualEnvChangeDetected{}
	}
	return virtualEnv, nil
}

// delete specified virtual env instance
func (r *ReconcileVirtualEnv) deleteVirtualEnv(namespace string, name string, logger logr.Logger) {
	err := shared.DeleteIns(r.client, namespace, name, &envv1alpha2.VirtualEnvironment{})
	if err != nil {
		logger.Error(err, "Failed to remove VirtualEnv instance "+name)
	} else {
		logger.Info("VirtualEnv deleted")
	}
}

// handle empty virtual env configure item with default value
func (r *ReconcileVirtualEnv) handleDefaultConfig(virtualEnv *envv1alpha2.VirtualEnvironment) {
	if virtualEnv.Spec.EnvHeader.Name == "" {
		virtualEnv.Spec.EnvHeader.Name = defaultEnvHeader
	}
	if virtualEnv.Spec.EnvHeader.Aliases != nil {
		for _, alias := range virtualEnv.Spec.EnvHeader.Aliases {
			if alias.Placeholder == "" {
				alias.Placeholder = defaultEnvHeaderAliasPlaceholder
			}
		}
	}
	if virtualEnv.Spec.EnvLabel.Name == "" {
		virtualEnv.Spec.EnvLabel.Name = defaultEnvLabel
	}
	if virtualEnv.Spec.EnvLabel.Splitter == "" {
		virtualEnv.Spec.EnvLabel.Splitter = defaultEnvSplitter
	}
}
