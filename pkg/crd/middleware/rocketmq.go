package middleware

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"fmt"
	"l0calh0st.cn/clickpaas-operator/pkg/crd"
	"strings"
)


// CreateRocketmqCRD registers a Rocketmq custom resource in kubernetes api.
func CreateRocketmqCRD(extensionClient apiextensions.Interface) error {
	crdName := strings.ToLower(fmt.Sprintf("%s.%s", crdv1alpha1.RocketmqPlural, crdv1alpha1.SchemeGroupVersion.Group))
	clusterCrd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: crdName,
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: crdv1alpha1.SchemeGroupVersion.Group,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:     crdv1alpha1.RocketmqPlural,
				Singular:   crdv1alpha1.RocketmqSingular,
				Kind:       crdv1alpha1.RocketmqKind,
				ShortNames: []string{crdv1alpha1.RocketmqShort},
			},
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:  crdv1alpha1.MiddlewareResourceVersion,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"spec": apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"replicas":        apiextensionsv1.JSONSchemaProps{Type: "integer"},
										"image":           apiextensionsv1.JSONSchemaProps{Type: "string"},
										"imagePullPolicy": apiextensionsv1.JSONSchemaProps{Type: "string"},
									},
								},
							},
						},
					},
					Served: true,
				},
			},
			Scope: apiextensionsv1.ResourceScope(apiextensionsv1.NamespaceScoped),
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


