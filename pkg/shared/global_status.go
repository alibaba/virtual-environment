package shared

import (
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sync"
)

// mutex to avoid controllers conflict
var Lock = sync.RWMutex{}

// mutex to reduce virtual env reconcile frequency cause by deployment/service change
var ReconcileTriggerLock = TriableMutex{}

// virtual env controller
var VirtualEnvController = new(controller.Controller)

// virtual env instance name
var VirtualEnvIns = ""

// service name -> selectors
var AvailableServices = make(map[string]map[string]string)

// deployment name -> labels
var AvailableDeployments = make(map[string]map[string]string)
