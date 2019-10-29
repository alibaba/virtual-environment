package status

// virtual env instance name
var VirtualEnvIns string = ""

// service name -> selectors
var AvailableServices = make(map[string]map[string]string)

// deployment name -> labels
var AvailableDeployments = make(map[string]map[string]string)
