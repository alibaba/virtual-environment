package main

import (
	"fmt"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	tlsDir               = `/run/secrets/tls`
	tlsCertFile          = `tls.crt`
	tlsKeyFile           = `tls.key`
	envVarName           = "VIRTUAL_ENVIRONMENT_TAG"
	sidecarContainerName = "istio-proxy"
)

var (
	podResource = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
	buildTime   string
)

// injectEnvironmentTag read the environment tag from pod label, and save to the sidecar container as an environment
// variable named `VIRTUAL_ENVIRONMENT_TAG`
func injectEnvironmentTag(req *v1beta1.AdmissionRequest) ([]PatchOperation, error) {
	logDebug("handling admission request for", req.Name)

	// This handler should only get called on Pod objects as per the MutatingWebhookConfiguration in the YAML file.
	// However, if (for whatever reason) this gets invoked on an object of a different kind, issue a log message but
	// let the object request pass through otherwise.
	if req.Resource != podResource {
		logError("expect resource to be", podResource)
		return nil, nil
	}

	// Parse the Pod object.
	raw := req.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &pod); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	// Retrieve the environment label name from pod label
	envLabels := os.Getenv(CONF_ENV_LABEL)
	if envLabels == "" {
		logFatal("cannot determine env label !!")
	}
	// Retrieve the environment tag from pod label
	envLabelList := strings.Split(envLabels, ",")
	envTag := ""
	for _, label := range envLabelList {
		if value, ok := pod.Labels[label]; ok {
			envTag = value
			break
		}
	}
	if envTag == "" {
		logError("no environment tag found on pod", getPodName(pod))
		return nil, nil
	}

	sidecarContainerIndex := -1
	for i, container := range pod.Spec.Containers {
		if container.Name == sidecarContainerName {
			sidecarContainerIndex = i
		}
	}
	if sidecarContainerIndex < 0 {
		logError("no sidecar container found on pod", getPodName(pod))
		return nil, nil
	}

	envVarIndex := -1
	for i, envVar := range pod.Spec.Containers[sidecarContainerIndex].Env {
		if envVar.Name == envVarName {
			envVarIndex = i
		}
	}

	// Create patch operations to apply environment tag
	var patches []PatchOperation
	if envVarIndex < 0 {
		patches = append(patches, PatchOperation{
			Op:    "add",
			Path:  fmt.Sprintf("/spec/containers/%d/env/0", sidecarContainerIndex),
			Value: corev1.EnvVar{Name: envVarName, Value: envTag},
		})
	} else {
		patches = append(patches, PatchOperation{
			Op:    "replace",
			Path:  fmt.Sprintf("/spec/containers/%d/env/%d/value", sidecarContainerIndex, envVarIndex),
			Value: envTag,
		})
	}

	logInfo("marked", getPodName(pod), "as", envTag)
	return patches, nil
}

func getPodName(pod corev1.Pod) string {
	if pod.Name == "" {
		return pod.GenerateName
	}
	return pod.Name
}

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	initLogger()
	logInfo("sidecar environment tag injector starting")
	logInfo("build version: " + buildTime)
	logInfo("environment labels: " + os.Getenv(CONF_ENV_LABEL))
	logInfo("log level: " + os.Getenv(CONF_LOG_LEVEL))

	mux := http.NewServeMux()
	mux.Handle("/inject", admitFuncHandler(injectEnvironmentTag))
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    ":8443",
		Handler: mux,
	}
	logFatal(server.ListenAndServeTLS(certPath, keyPath))
}
