package shared

import (
	"sync"
)

type ServiceInfo struct {
	Selectors  map[string]string
	Ports      []uint32
	Gateways   []string
	Hosts      []string
	CustomRule string
}

// mutex to avoid controllers conflict
var Lock = sync.RWMutex{}

// virtual env instance name
var VirtualEnvIns = ""

// service name -> service info
var AvailableServices = make(map[string]ServiceInfo)

// "resource type # name" -> labels
var AvailableLabels = make(map[string]map[string]string)
