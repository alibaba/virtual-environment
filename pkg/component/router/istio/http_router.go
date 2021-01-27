package istio

import (
	envv1alpha2 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha2"
	"alibaba.com/virtual-env-operator/pkg/component/router/common"
	"alibaba.com/virtual-env-operator/pkg/component/router/istio/envoy"
	"alibaba.com/virtual-env-operator/pkg/component/router/istio/http"
	"alibaba.com/virtual-env-operator/pkg/shared/logger"
	"context"
	networkingv1alpha3api "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	networkingv1alpha3 "knative.dev/pkg/apis/istio/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type HttpRouter struct {
}

// generate virtual services and destination rules
func (r *HttpRouter) GenerateRoute(client client.Client, scheme *runtime.Scheme, virtualEnv *envv1alpha2.VirtualEnvironment,
	namespace string, svcName string, availableLabels []string, relatedDeployments []string) error {
	shouldBeDeleted, err := r.reconcileVirtualService(client, scheme, virtualEnv, namespace, svcName, availableLabels, relatedDeployments)
	if err != nil {
		return err
	}
	return r.reconcileDestinationRule(client, scheme, virtualEnv, namespace, svcName, relatedDeployments, shouldBeDeleted)
}

// clean up virtual services and destination rules
func (r *HttpRouter) CleanupRoute(client client.Client, namespace string, svcName string) error {
	err := http.DeleteVirtualService(client, namespace, svcName)
	if err != nil {
		logger.Error(err, "Failed to remove VirtualService instance "+svcName)
	} else {
		logger.Info("VirtualService deleted " + svcName)
	}
	err = http.DeleteDestinationRule(client, namespace, svcName)
	if err != nil {
		logger.Error(err, "Failed to remove DestinationRule instance "+svcName)
	} else {
		logger.Info("DestinationRule deleted " + svcName)
	}
	return nil
}

// watch for changes to secondary resource VirtualService & DestinationRule, requeue their owner to VirtualEnv
func (r *HttpRouter) RegisterReconcileWatcher(c controller.Controller) error {
	err := c.Watch(&source.Kind{Type: &networkingv1alpha3.VirtualService{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &envv1alpha2.VirtualEnvironment{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &networkingv1alpha3.DestinationRule{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &envv1alpha2.VirtualEnvironment{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &networkingv1alpha3api.EnvoyFilter{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &envv1alpha2.VirtualEnvironment{},
	})
	if err != nil {
		return err
	}
	return nil
}

// look for envoy filter instance in namespace
func (r *HttpRouter) CheckTagAppender(client client.Client, virtualEnv *envv1alpha2.VirtualEnvironment,
	namespace string, name string) common.TagAppenderStatus {
	envoyFilter := &networkingv1alpha3api.EnvoyFilter{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, envoyFilter)
	if errors.IsNotFound(err) {
		return common.NotExist
	} else if err == nil {
		if envoy.IsDifferentTagAppender(envoyFilter, virtualEnv.Spec.EnvLabel.Name, virtualEnv.Spec.EnvHeader.Name) {
			return common.Outdated
		}
		return common.UpToDate
	}
	return common.Unknown
}

// create envoy filter to automatically append tag to service
func (r *HttpRouter) CreateTagAppender(client client.Client, scheme *runtime.Scheme, virtualEnv *envv1alpha2.VirtualEnvironment,
	namespace string, name string) error {
	tagAppender, err := envoy.TagAppenderFilter(namespace, name, virtualEnv.Spec.EnvLabel.Name, virtualEnv.Spec.EnvHeader.Name)
	if tagAppender == nil {
		logger.Error(err, "Failed to initialize TagAppender instance", namespace+":"+name)
		return err
	}
	// set VirtualEnv instance as the owner and controller
	err = controllerutil.SetControllerReference(virtualEnv, tagAppender, scheme)
	if err == nil {
		err = client.Create(context.TODO(), tagAppender)
	}
	if err != nil {
		logger.Error(err, "Failed to create TagAppender instance "+namespace+":"+name)
	} else {
		logger.Info("TagAppender created " + namespace + ":" + name)
	}
	return err
}

// delete auto tag appender envoy filter
func (r *HttpRouter) DeleteTagAppender(client client.Client, namespace string, name string) error {
	err := envoy.DeleteTagAppenderIfExist(client, namespace, name)
	if err == nil {
		logger.Info("TagAppender deleted " + namespace + ":" + name)
	}
	return err
}

// reconcile virtual service according to related pod-resources and available labels
func (r *HttpRouter) reconcileVirtualService(client client.Client, scheme *runtime.Scheme, virtualEnv *envv1alpha2.VirtualEnvironment,
	namespace string, svcName string, availableLabels []string, relatedDeployments []string) (bool, error) {
	virtualSvc := http.VirtualService(namespace, svcName, availableLabels, relatedDeployments, virtualEnv.Spec)
	// when no http port or less than 2 destination available, and no custom gateway configured,
	// virtual service instance should be removed
	shouldBeDeleted := len(virtualSvc.Spec.HTTP) < 2 && len(virtualSvc.Spec.Gateways) < 2
	foundVirtualSvc := &networkingv1alpha3.VirtualService{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: svcName, Namespace: namespace}, foundVirtualSvc)
	if err != nil {
		// VirtualService not exist, create one
		if shouldBeDeleted {
			return shouldBeDeleted, nil
		} else if errors.IsNotFound(err) {
			err = r.createVirtualService(client, scheme, virtualEnv, virtualSvc)
			if err != nil {
				logger.Error(err, "Failed to create new VirtualService")
				return shouldBeDeleted, err
			}
		} else {
			logger.Error(err, "Failed to get VirtualService")
			return shouldBeDeleted, err
		}
	} else if shouldBeDeleted {
		// VirtualService should be remove
		logger.Info("Deleting VirtualService instance " + virtualSvc.Name)
		err := client.Delete(context.TODO(), foundVirtualSvc)
		if err != nil {
			logger.Error(err, "Failed to delete VirtualService instance")
			return shouldBeDeleted, err
		}
	} else if http.IsDifferentVirtualService(&foundVirtualSvc.Spec, &virtualSvc.Spec) {
		// existing VirtualService changed
		foundVirtualSvc.Spec = virtualSvc.Spec
		err := client.Update(context.TODO(), foundVirtualSvc)
		if err != nil {
			logger.Error(err, "Failed to update VirtualService status")
			return shouldBeDeleted, err
		}
		logger.Info("VirtualService " + virtualSvc.Name + " changed")
	}
	return shouldBeDeleted, nil
}

// reconcile destination rule according to related pod-resources
func (r *HttpRouter) reconcileDestinationRule(client client.Client, scheme *runtime.Scheme, virtualEnv *envv1alpha2.VirtualEnvironment,
	namespace string, svcName string, relatedDeployments []string, shouldBeDeleted bool) error {
	destRule := http.DestinationRule(namespace, svcName, relatedDeployments, virtualEnv.Spec.EnvLabel.Name)
	foundDestRule := &networkingv1alpha3.DestinationRule{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: svcName, Namespace: namespace}, foundDestRule)
	if err != nil {
		// DestinationRule not exist, create one
		if shouldBeDeleted {
			return nil
		} else if errors.IsNotFound(err) {
			err = r.createDestinationRule(client, scheme, virtualEnv, destRule)
			if err != nil {
				logger.Error(err, "Failed to create new DestinationRule")
				return err
			}
		} else {
			logger.Error(err, "Failed to get DestinationRule")
			return err
		}
	} else if shouldBeDeleted {
		// DestinationRule should be remove
		logger.Info("Deleting DestinationRule instance " + destRule.Name)
		err := client.Delete(context.TODO(), foundDestRule)
		if err != nil {
			logger.Error(err, "Failed to delete DestinationRule instance")
			return err
		}
	} else if http.IsDifferentDestinationRule(&foundDestRule.Spec, &destRule.Spec, virtualEnv.Spec.EnvLabel.Name) {
		// existing DestinationRule changed
		foundDestRule.Spec = destRule.Spec
		err := client.Update(context.TODO(), foundDestRule)
		if err != nil {
			logger.Error(err, "Failed to update DestinationRule status")
			return err
		}
		logger.Info("DestinationRule " + destRule.Name + " changed")
	}
	return nil
}

// create virtual service instance
func (r *HttpRouter) createVirtualService(client client.Client, scheme *runtime.Scheme,
	virtualEnv *envv1alpha2.VirtualEnvironment, virtualSvc *networkingv1alpha3.VirtualService) error {
	// set VirtualEnv instance as the owner and controller
	err := controllerutil.SetControllerReference(virtualEnv, virtualSvc, scheme)
	if err != nil {
		logger.Error(err, "Failed to set owner of "+virtualSvc.Name)
		return err
	}
	err = client.Create(context.TODO(), virtualSvc)
	if err != nil {
		logger.Error(err, "Failed to create VirtualService "+virtualSvc.Name)
		return err
	}
	logger.Info("VirtualService " + virtualSvc.Name + " created")
	return nil
}

// create destination rule instance
func (r *HttpRouter) createDestinationRule(client client.Client, scheme *runtime.Scheme,
	virtualEnv *envv1alpha2.VirtualEnvironment, destRule *networkingv1alpha3.DestinationRule) error {
	// set VirtualEnv instance as the owner and controller
	err := controllerutil.SetControllerReference(virtualEnv, destRule, scheme)
	if err != nil {
		logger.Error(err, "Failed to set owner of "+destRule.Name)
		return err
	}
	err = client.Create(context.TODO(), destRule)
	if err != nil {
		logger.Error(err, "Failed to create DestinationRule "+destRule.Name)
		return err
	}
	logger.Info("DestinationRule " + destRule.Name + " created")
	return nil
}
