package shared

import (
	"sync"
)

type ServiceInfo struct {
	Selectors map[string]string
	Ports     []uint32
	Gateways  []string
	Hosts     []string
}

// mutex to avoid controllers conflict
var Lock = sync.RWMutex{}

// virtual env instance name
var VirtualEnvIns = ""

// service name -> service info
var AvailableServices = make(map[string]ServiceInfo)

// deployment name -> labels
var AvailableDeployments = make(map[string]map[string]string)
