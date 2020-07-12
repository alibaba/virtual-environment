package util

func LabelMark(resourceType string, name string) string {
	return resourceType + "#" + name
}
