package loki

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type jsonSchemePropsType string

const (
	jsonSchemePropsTypeAsInteger jsonSchemePropsType = "integer"
	jsonSchemePropsTypeAsString  jsonSchemePropsType = "string"
	jsonSchemePropsTypeAsObject  jsonSchemePropsType = "object"
	jsonSchemePropsTypesAsNumber jsonSchemePropsType = "number"
	jsonSchemePropsTypeAsArray   jsonSchemePropsType = "array"
)

func NewCustomResourceDefine() *apiextensionsv1.CustomResourceDefinition {
	crd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "lokis" + "." + crapiv1alpha1.SchemeGroupVersion.Group,
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: crapiv1alpha1.SchemeGroupVersion.Group,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:   "lokis",
				Singular: "loki",
				Kind:     "Loki",
				ListKind: "LokiList",
			},
			Scope: apiextensionsv1.ResourceScope(apiextensionsv1.NamespaceScoped),
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    crapiv1alpha1.Version,
					Storage: true,
					Served:  true,
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
					Schema: &apiextensionsv1.CustomResourceValidation{OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Type: "object",
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"apiVersion": {Type: string(jsonSchemePropsTypeAsString)},
							"kind":       {Type: string(jsonSchemePropsTypeAsString)},
							"metadata":   {Type: string(jsonSchemePropsTypeAsObject)},
							"spec": {Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"replicas":  {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired loki instance in k8s. Default 1"},
								"image":     {Type: string(jsonSchemePropsTypeAsString), Description: "The image of loki used."},
								"configMap": {Type: string(jsonSchemePropsTypeAsString), Description: "Name of configmap. If you special the configmap, the loki will use configmap as it's config or it will use internal default config"},
								"config": {Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
								}},
								"deployMode": {Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"monolithic": {Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"all": {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired loki instance in k8s. Default 1"},
									}},
									"sample": {Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"read":  {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired read target in loki stack. Default 1"},
										"write": {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired write target in loki stack. Default 1"},
									}},
									"microservice": {Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"ingester":        {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired ingester pod in loki stack. Default 1"},
										"distributor":     {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired distributor pod in loki stack. Default 1"},
										"query-frontent":  {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired query-frontent pod in loki stack. Default 1"},
										"query-scheduler": {Type: string(jsonSchemePropsTypeAsInteger), Description: "Nunber of desired query-scheduler pod in loki stack, Default 1"},
										"querier":         {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired querier pod in loki stack. Default 1"},
										"index-gateway":   {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired  index-gateway pod in loki stack. Default 1"},
										"ruler":           {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of desired ruler pod in loki stack. Default 1"},
										"compactor":       {Type: string(jsonSchemePropsTypeAsInteger), Description: "Number of compactor pod in loki stack. Default 1"},
									}, Description: `The microservices deployment mode instantiates components of Loki as distinct processesConfig replicas of each target`},
																							}},
														}, Description: "Spec describes the specification of Loki applications using kubernetes as a cluster manager"},
							"status": {Type: string(jsonSchemePropsTypeAsObject), Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"phase": {Type: string(jsonSchemePropsTypeAsObject)},
							}, Description: "Status if the current running status of Loki in k8s. This data may be out of date by some window of time"},
						},
						Required: []string{"apiVersion", "kind", "metadata", "spec"},
					}},
				},
			},
		},
		Status: apiextensionsv1.CustomResourceDefinitionStatus{},
	}
	return crd
}
