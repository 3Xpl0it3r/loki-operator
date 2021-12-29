package promtail

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)



type jsonSchemePropsType string
const (
	jsonSchemePropsTypeAsInteger jsonSchemePropsType = "integer"
	jsonSchemePropsTypeAsString jsonSchemePropsType = "string"
	jsonSchemePropsTypeAsObject jsonSchemePropsType = "object"
	jsonSchemePropsTypesAsNumber jsonSchemePropsType = "number"
	jsonSchemePropsTypeAsArray jsonSchemePropsType = "array"
)

func NewCustomResourceDefine()*apiextensionsv1.CustomResourceDefinition{
	crd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "promtails"+"." + crapiv1alpha1.SchemeGroupVersion.Group,
		},
		Spec:       apiextensionsv1.CustomResourceDefinitionSpec{
			Group: crapiv1alpha1.SchemeGroupVersion.Group,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:     "promtails",
				Singular:   "promtail",
				Kind:       "Promtail",
				ListKind:   "PromtailList",
			},
			Scope: apiextensionsv1.ResourceScope(apiextensionsv1.NamespaceScoped),
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name: crapiv1alpha1.Version,
					Storage: true,
					Served: true,
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
					Schema: &apiextensionsv1.CustomResourceValidation{OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Type:                   "object",
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"apiVersion": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsString)},
							"kind": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsString)},
							"metadata": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject)},
							"spec": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"replicas": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsInteger)},
								"image": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsString)},
								"configMap": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsString)},
								"config": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"server": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"http_listen_address": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsString)},
									},Description: `The server block configures Promtail’s behavior as an HTTP server`},
									"clients": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject),Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"url": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsString)},
									},Description:  `The clients block configures how Promtail connects to an instance of Loki`},
									"positions": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject),Description: `The positions block configures where Promtail will save a file indicating how far it has read into a file. 
It is needed for when Promtail is restarted to allow it to continue from where it left off.`},
									"scrape_configs": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject)},
								},Description: `Promtail is configured in a YAML file (usually referred to as config.yaml) which contains information on the Promtail server, where positions are stored, and how to scrape logs from files`},
							}, Description: "Spec describes the specification of Promtail applications using kubernetes as a cluster manager"},
							"status": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"phase": apiextensionsv1.JSONSchemaProps{Type: string(jsonSchemePropsTypeAsObject)},
							}, Description: "Status is the current running status of promtail in k8s cluster. This data may be out of date by some window of time."},
						},
						Required: []string{"apiVersion", "kind", "metadata", "spec"},
					}},
				},
			},
		},
		Status:     apiextensionsv1.CustomResourceDefinitionStatus{},
	}
	return crd
}
