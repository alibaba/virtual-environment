package api

import (
	"alibaba.com/virtual-env-operator/pkg/shared"
	"github.com/labstack/echo/v4"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
)

var log = logf.Log.WithName("inspect-api")

func Start(inspectHost string, inspectPort int) {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/inspect/deployment", inspectDeployment)
	e.GET("/inspect/service", inspectService)
	e.GET("/inspect/global", inspectGlobalVariable)
	e.POST("/trigger", triggerReconcile)

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

func inspectDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, shared.AvailableDeployments)
}

func inspectService(c echo.Context) error {
	return c.JSON(http.StatusOK, shared.AvailableServices)
}

func inspectGlobalVariable(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"VirtualEnvIns":  shared.VirtualEnvIns,
		"InsNamePostfix": shared.InsNamePostfix,
	})
}

func triggerReconcile(c echo.Context) error {
	name := shared.VirtualEnvIns
	namespace := os.Getenv("WATCH_NAMESPACE")
	if len(name) > 0 && len(namespace) > 0 {
		_, err := (*shared.VirtualEnvController).Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{Name: name, Namespace: namespace},
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "virtual environment reconciled")
		}
	}
	return c.String(http.StatusOK, "virtual environment not present")
}
