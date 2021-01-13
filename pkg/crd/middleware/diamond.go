package middleware

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"l0calh0st.cn/clickpaas-operator/pkg/crd"
	"strings"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"fmt"
)


// CreateDiamondCRD registers a Diamond custom resource in kubernetes api.
func CreateDiamondCRD(extensionClient apiextensions.Interface) error {
	crdName := strings.ToLower(fmt.Sprintf("%s.%s", crdv1alpha1.DiamondPlural, crdv1alpha1.SchemeGroupVersion.Group))
	clusterCrd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: crdName,
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: crdv1alpha1.SchemeGroupVersion.Group,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:     crdv1alpha1.DiamondPlural,
				Singular:   crdv1alpha1.DiamondSingular,
				Kind:       crdv1alpha1.DiamondKind,
				ShortNames: []string{crdv1alpha1.DiamondShort},
			},
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:  crdv1alpha1.MiddlewareResourceVersion,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							Type: ValidateTypeAsObject,
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"spec": apiextensionsv1.JSONSchemaProps{
									Type: ValidateTypeAsObject,
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"replicas":        apiextensionsv1.JSONSchemaProps{Type: ValidateTypeAsInt},
										"image":           apiextensionsv1.JSONSchemaProps{Type: ValidateTypeAsString},
										"imagePullPolicy": apiextensionsv1.JSONSchemaProps{Type: ValidateTypeAsString},
										"config": apiextensionsv1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]apiextensionsv1.JSONSchemaProps{
												"host": {Type: ValidateTypeAsString},
												"user": {Type: ValidateTypeAsString},
												"password": {Type: ValidateTypeAsString},
												"port": {Type: ValidateTypeAsInt},
											},
										},
									},
								},
								"status": apiextensionsv1.JSONSchemaProps{
									Type: ValidateTypeAsObject,
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"initialDb": apiextensionsv1.JSONSchemaProps{Type: ValidDateTypeAsBoolean},
									},
								},
							},
						},
					},
					Served: true,
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
				},
			},
			Scope: apiextensionsv1.ResourceScope(apiextensionsv1.NamespaceScoped),
		},
		Status: apiextensionsv1.CustomResourceDefinitionStatus{
			Conditions:     nil,
			AcceptedNames:  apiextensionsv1.CustomResourceDefinitionNames{},
			StoredVersions: nil,
		},
	}
	err := crd.RegisterCRD(extensionClient, clusterCrd)

	if err != nil {
		return err
	}

	err = crd.WaitForCRDEstablished(extensionClient, crdName)
	if err != nil {
		return errors.NewAggregate([]error{err, crd.UnregisterCRD(extensionClient, crdName)})
	}
	return nil
}


