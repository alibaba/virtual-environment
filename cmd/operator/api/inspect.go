package api

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	"alibaba.com/virtual-env-operator/version"
	"github.com/labstack/echo/v4"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
)

var log = logf.Log.WithName("inspect-api")

func Start(inspectHost string, inspectPort int) {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/status", inspectGlobalVariable)
	e.GET("/version", inspectBuildVersion)

	// Start server
	inspectAddr := inspectHost + ":" + strconv.Itoa(inspectPort)
	go func() {
		e.HideBanner = true
		e.HidePort = true
		log.Info("Starting inspect api", "addr", inspectAddr)
		err := e.Start(inspectAddr)
		if err != nil {
			log.Error(err, "Inspect api cannot listen to "+inspectAddr)
		}
	}()
}

func inspectGlobalVariable(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"VirtualEnv": shared.VirtualEnvIns,
		"Services":   shared.AvailableServices,
		"Labels":     shared.AvailableLabels,
	})
}

func inspectBuildVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Version":   version.Version,
		"BuildTime": version.BuildTime,
	})
}
