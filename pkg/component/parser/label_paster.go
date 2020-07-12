package parser

// return map of deployment name to virtual label value
func FindAllRelatedLabels(availableLabels map[string]map[string]string, selector map[string]string,
	envLabel string) map[string]string {
	relatedDeployments := make(map[string]string)
	for dep, labels := range availableLabels {
		match := true
		for k, v := range selector {
			if labels[k] != v {
				match = false
				break
			}
		}
		if _, exist := labels[envLabel]; match && exist {
			relatedDeployments[dep] = labels[envLabel]
		}
	}
	return relatedDeployments
}

// list all possible values in deployment virtual env label
func FindAllVirtualEnvLabelValues(availableLabels map[string]map[string]string, envLabel string) []string {
	labelSet := make(map[string]bool)
	for _, labels := range availableLabels {
		labelVal, exist := labels[envLabel]
		if exist {
			labelSet[labelVal] = true
		}
	}
	return GetKeys(labelSet)
}

// get all keys of a map as array
func GetKeys(kv map[string]bool) []string {
	keys := make([]string, 0, len(kv))
	for k, _ := range kv {
		keys = append(keys, k)
	}
	return keys
}
