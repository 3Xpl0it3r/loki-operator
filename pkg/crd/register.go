package crd

import (
	"context"
	"fmt"
	"os"
	"syscall"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	extensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// RegisterCRD register CustomResourceDefine according crd filename into k8s apiserver
func RegisterCRDWithFile(namespace string, extClientSet extensionsclientset.Interface, filename string) error {
	crd := new(apiextensionsv1.CustomResourceDefinition)

	fp, err := os.OpenFile(filename, syscall.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open filename %s failed , %v", filename, err)
	}
	decoder := yaml.NewYAMLToJSONDecoder(fp)
	if err = decoder.Decode(crd); err != nil {
		return fmt.Errorf("unmarshal yaml to crd failed: %v", err)
	}
	crd.SetNamespace(namespace)

	return RegisterCRDWithObj(extClientSet, crd)
}

func RegisterCRDWithObj(extClientSet extensionsclientset.Interface, crdObj *apiextensionsv1.CustomResourceDefinition) error {
	if _, err := extClientSet.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crdObj, metav1.CreateOptions{}); err != nil {
		if k8serror.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

// UnRegisterCRD remove custom resource define from k8s apiserver
func UnRegisterCRD(extClientSet extensionsclientset.Interface, crdName string) error {
	return extClientSet.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), crdName, metav1.DeleteOptions{})
}
