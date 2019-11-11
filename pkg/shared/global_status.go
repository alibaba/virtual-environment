package shared

import (
	"sync"
)

// mutex to avoid controllers conflict
var Lock = sync.RWMutex{}

// virtual env instance name
var VirtualEnvIns = ""

// service name -> selectors
var AvailableServices = make(map[string]map[string]string)

// deployment name -> labels
var AvailableDeployments = make(map[string]map[string]string)
