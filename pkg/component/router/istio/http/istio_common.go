package http

import "regexp"

// replace invalid chars in subset name
func toSubsetName(labelValue string) string {
	re, _ := regexp.Compile("[_.]")
	return re.ReplaceAllString(labelValue, "-")
}
