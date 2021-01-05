package crd

import (
	"context"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)




// RegisterCRD registers given custom resource definition into the kubernetes api.
func RegisterCRD(clientSet apiextensions.Interface, crd *apiextensionsv1.CustomResourceDefinition) error {


	crdInterface := clientSet.ApiextensionsV1().CustomResourceDefinitions()
	_, err := crdInterface.Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

// UnregisterCRD removes custom resource definition from the kubernetes api.
func UnregisterCRD(clientSet apiextensions.Interface, crdName string) error {
	crdInterface := clientSet.ApiextensionsV1().CustomResourceDefinitions()
	return crdInterface.Delete(context.TODO(), crdName, metav1.DeleteOptions{})
}




