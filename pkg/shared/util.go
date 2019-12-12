package shared

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// delete specified instance
func DeleteIns(client client.Client, namespace string, name string, obj runtime.Object) error {
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, obj)
	if err != nil {
		return err
	}
	err = client.Delete(context.TODO(), obj)
	if err != nil {
		return err
	}
	return nil
}

// put postfix after target name if it's not empty
func NameWithPostfix(name string, postfix string) string {
	if len(postfix) > 0 {
		return name + "-" + postfix
	}
	return name
}
