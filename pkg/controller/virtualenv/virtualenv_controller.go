package virtualenv

import (
	envv1alpha1 "alibaba.com/virtual-env-operator/pkg/apis/env/v1alpha1"
	"alibaba.com/virtual-env-operator/pkg/status"
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis/istio/common/v1alpha1"
	networkingv1alpha3 "knative.dev/pkg/apis/istio/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
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

	virtualEnv, err := r.fetchVirtualEnvIns(request, reqLogger)
	if virtualEnv == nil {
		status.Lock.Unlock()
		return reconcile.Result{}, err
	}

	reqLogger.Info("Responding VirtualService and DestinationRule")
	for srv, selector := range status.AvailableServices {
		availableLabels := findAllVirtualEnvLabelValues(status.AvailableDeployments, virtualEnv.Spec.VeLabel)
		relatedDeployments := findAllRelatedDeployments(status.AvailableDeployments, selector, virtualEnv.Spec.VeLabel)

		if len(availableLabels) > 0 && len(relatedDeployments) > 0 {
			err = r.reconcileVirtualService(virtualEnv, srv, request, availableLabels, relatedDeployments, reqLogger)
			if err != nil {
				status.Lock.Unlock()
				return reconcile.Result{}, err
			}
			err = r.reconcileDestinationRule(virtualEnv, srv, request, relatedDeployments, reqLogger)
			if err != nil {
				status.Lock.Unlock()
				return reconcile.Result{}, err
			}
		}
	}

	status.Lock.Unlock()
	return reconcile.Result{}, nil
}

// fetch the VirtualEnv instance from request
func (r *ReconcileVirtualEnv) fetchVirtualEnvIns(request reconcile.Request, reqLogger logr.Logger) (*envv1alpha1.VirtualEnv, error) {
	virtualEnv := &envv1alpha1.VirtualEnv{}
	err := r.client.Get(context.TODO(), request.NamespacedName, virtualEnv)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("VirtualEnv resource missing")
			if status.VirtualEnvIns == request.Name {
				reqLogger.Info("VirtualEnv deleted")
				status.VirtualEnvIns = ""
			}
			return nil, nil
		}
		reqLogger.Error(err, "Failed to get VirtualEnv")
		return nil, err
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
	return virtualEnv, err
}

func (r *ReconcileVirtualEnv) reconcileDestinationRule(virtualEnv *envv1alpha1.VirtualEnv, srv string, request reconcile.Request, relatedDeployments map[string]string, reqLogger logr.Logger) error {
	destRule := r.destinationRule(virtualEnv, srv, request.Namespace, relatedDeployments)
	foundDestRule := &networkingv1alpha3.DestinationRule{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: srv, Namespace: request.Namespace}, foundDestRule)
	if err != nil {
		// DestinationRule not exist, create one
		if errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), destRule)
			if err != nil {
				reqLogger.Error(err, "Failed to create DestinationRule "+destRule.Name)
				return err
			}
			reqLogger.Info("DestinationRule " + destRule.Name + " created")
		} else {
			reqLogger.Error(err, "Failed to get DestinationRule")
			return err
		}
	} else if isDifferentDestinationRule(foundDestRule.Spec, destRule.Spec, virtualEnv.Spec.VeLabel) {
		// existing DestinationRule changed
		foundDestRule.Spec = destRule.Spec
		err := r.client.Update(context.TODO(), foundDestRule)
		if err != nil {
			reqLogger.Error(err, "Failed to update DestinationRule status")
			return err
		}
		reqLogger.Info("DestinationRule " + destRule.Name + " changed")
	}
	return nil
}

func (r *ReconcileVirtualEnv) reconcileVirtualService(virtualEnv *envv1alpha1.VirtualEnv, srv string, request reconcile.Request,
	availableLabels []string, relatedDeployments map[string]string, reqLogger logr.Logger) error {
	virtualSrv := r.virtualService(virtualEnv, srv, request.Namespace, availableLabels, relatedDeployments)
	foundVirtualSrv := &networkingv1alpha3.VirtualService{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: srv, Namespace: request.Namespace}, foundVirtualSrv)
	if err != nil {
		// VirtualService not exist, create one
		if errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), virtualSrv)
			if err != nil {
				reqLogger.Error(err, "Failed to create VirtualService "+virtualSrv.Name)
				return err
			}
			reqLogger.Info("VirtualService " + virtualSrv.Name + " created")
		} else {
			reqLogger.Error(err, "Failed to get VirtualService")
			return err
		}
	} else if isDifferentVirtualService(foundVirtualSrv.Spec, virtualSrv.Spec, virtualEnv.Spec.VeHeader) {
		// existing VirtualService changed
		foundVirtualSrv.Spec = virtualSrv.Spec
		err := r.client.Update(context.TODO(), foundVirtualSrv)
		if err != nil {
			reqLogger.Error(err, "Failed to update VirtualService status")
			return err
		}
		reqLogger.Info("VirtualService " + virtualSrv.Name + " changed")
	}
	return nil
}

// check whether DestinationRule is different
func isDifferentDestinationRule(spec1 networkingv1alpha3.DestinationRuleSpec,
	spec2 networkingv1alpha3.DestinationRuleSpec, label string) bool {
	if len(spec1.Subsets) != len(spec2.Subsets) {
		return true
	}
	for _, subset1 := range spec1.Subsets {
		subset2 := findSubsetByName(spec2.Subsets, subset1.Name)
		if subset2 == nil {
			return true
		}
		if subset1.Labels[label] != subset2.Labels[label] {
			return true
		}
	}
	return false
}

// find subset from list
func findSubsetByName(subsets []networkingv1alpha3.Subset, name string) *networkingv1alpha3.Subset {
	for _, subset := range subsets {
		if subset.Name == name {
			return &subset
		}
	}
	return nil
}

// check whether VirtualService is different
func isDifferentVirtualService(spec1 networkingv1alpha3.VirtualServiceSpec, spec2 networkingv1alpha3.VirtualServiceSpec, header string) bool {
	if len(spec1.HTTP) != len(spec2.HTTP) {
		return true
	}
	for _, route1 := range spec1.HTTP {
		if route1.Match == nil {
			continue
		}
		if !findMatchRoute(spec2.HTTP, &route1, header) {
			return true
		}
	}
	return false
}

// check whether HTTPRoute exist in list
func findMatchRoute(routes []networkingv1alpha3.HTTPRoute, target *networkingv1alpha3.HTTPRoute, header string) bool {
	for _, route := range routes {
		if route.Match == nil {
			continue
		}
		if route.Route[0].Destination.Subset == target.Route[0].Destination.Subset &&
			route.Match[0].Headers[header] == target.Match[0].Headers[header] {
			return true
		}
	}
	return false
}

// generate istio virtual service instance
func (r *ReconcileVirtualEnv) virtualService(e *envv1alpha1.VirtualEnv, name string, namespace string,
	availableLabels []string, relatedDeployments map[string]string) *networkingv1alpha3.VirtualService {
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
						Host: name,
					},
				}},
			}},
		},
	}
	for _, label := range availableLabels {
		matchRoute, ok := virtualServiceMatchRoute(name, relatedDeployments, label, e.Spec.VeHeader, e.Spec.VeSplitter)
		if ok {
			virtualSrv.Spec.HTTP = append(virtualSrv.Spec.HTTP, matchRoute)
		}
	}
	// Set VirtualEnv instance as the owner and controller
	controllerutil.SetControllerReference(e, virtualSrv, r.scheme)
	return virtualSrv
}

// generate istio destination rule instance
func (r *ReconcileVirtualEnv) destinationRule(e *envv1alpha1.VirtualEnv, name string, namespace string,
	relatedDeployments map[string]string) *networkingv1alpha3.DestinationRule {
	destRule := &networkingv1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.DestinationRuleSpec{
			Host:    name,
			Subsets: []networkingv1alpha3.Subset{},
		},
	}
	for dep, label := range relatedDeployments {
		destRule.Spec.Subsets = append(destRule.Spec.Subsets, destinationRuleMatchSubset(dep, e.Spec.VeLabel, label))
	}
	// Set VirtualEnv instance as the owner and controller
	controllerutil.SetControllerReference(e, destRule, r.scheme)
	return destRule
}

// generate istio destination rule subset instance
func destinationRuleMatchSubset(name string, labelKey string, labelValue string) networkingv1alpha3.Subset {
	return networkingv1alpha3.Subset{
		Name: name,
		Labels: map[string]string{
			labelKey: labelValue,
		},
	}
}

func deleteVirtualEnv(namespace string, virtualEnv string) {
	//TODO: delete virtual env instance
}

// return map of deployment name to virtual label value
func findAllRelatedDeployments(deployments map[string]map[string]string, selector map[string]string, velabel string) map[string]string {
	relatedDeployments := make(map[string]string)
	for dep, labels := range deployments {
		match := true
		for k, v := range selector {
			if labels[k] != v {
				match = false
				break
			}
		}
		if _, exist := labels[velabel]; match && exist {
			relatedDeployments[dep] = labels[velabel]
		}
	}
	return relatedDeployments
}

// list all possible values in deployment virtual env label
func findAllVirtualEnvLabelValues(deployments map[string]map[string]string, velabel string) []string {
	labelSet := make(map[string]bool)
	for _, labels := range deployments {
		labelVal, exist := labels[velabel]
		if exist {
			labelSet[labelVal] = true
		}
	}
	return getKeys(labelSet)
}

// get all keys of a map as array
func getKeys(kv map[string]bool) []string {
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	return keys
}

// calculate and generate http route instance
func virtualServiceMatchRoute(serviceName string, relatedDeployments map[string]string, labelVal string, headerKey string,
	splitter string) (networkingv1alpha3.HTTPRoute, bool) {
	var possibleRoutes []string
	for k, v := range relatedDeployments {
		if leveledEqual(v, labelVal, splitter) {
			possibleRoutes = append(possibleRoutes, k)
		}
	}
	if len(possibleRoutes) > 0 {
		return matchRoute(serviceName, headerKey, labelVal, findLongestString(possibleRoutes)), true
	}
	return networkingv1alpha3.HTTPRoute{}, false
}

// generate istio virtual service http route instance
func matchRoute(serviceName string, headerKey string, labelVal string, matchedLabel string) networkingv1alpha3.HTTPRoute {
	return networkingv1alpha3.HTTPRoute{
		Route: []networkingv1alpha3.HTTPRouteDestination{{
			Destination: networkingv1alpha3.Destination{
				Host:   serviceName,
				Subset: matchedLabel,
			},
		}},
		Match: []networkingv1alpha3.HTTPMatchRequest{{
			Headers: map[string]v1alpha1.StringMatch{
				headerKey: {
					Exact: labelVal,
				},
			},
		}},
	}
}

// get the longest string in list
func findLongestString(strings []string) string {
	mostLongStr := ""
	for _, str := range strings {
		if len(str) > len(mostLongStr) {
			mostLongStr = str
		}
	}
	return mostLongStr
}

// check whether source string match target string at any level
func leveledEqual(source string, target string, splitter string) bool {
	for {
		if source == target {
			return true
		}
		if strings.Contains(source, splitter) {
			source = source[0:strings.LastIndex(source, splitter)]
		} else {
			return false
		}
	}
}
