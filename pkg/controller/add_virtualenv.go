package controller

import (
	"alibaba.com/virtual-env-operator/pkg/controller/virtualenv"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, virtualenv.Add)
}
