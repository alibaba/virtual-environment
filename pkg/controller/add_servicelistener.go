package controller

import (
	"alibaba.com/virtual-env-operator/pkg/controller/servicelistener"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, servicelistener.Add)
}
