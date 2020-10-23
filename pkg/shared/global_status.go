package shared

import (
	"k8s.io/apimachinery/pkg/types"
	"sync"
)

type ServiceInfo struct {
	Selectors  map[string]string
	Ports      map[string]uint32
	Gateways   []string
	Hosts      []string
	CustomRule string
}

// mutex to avoid controllers conflict
var Lock = sync.RWMutex{}

// virtual env instance name
var VirtualEnvIns *types.NamespacedName = nil

// service name -> service info
var AvailableServices = make(map[string]ServiceInfo)

// "resource type # name" -> labels
var AvailableLabels = make(map[string]map[string]string)
