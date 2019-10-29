package status

// virtual env instance name
var VirtualEnvIns string

// service name -> selectors
var AvailableServices map[string]map[string]string

// deployment name -> labels
var AvailableDeployments map[string]map[string]string
