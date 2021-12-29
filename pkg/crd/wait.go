package crd

import (
	"context"
	"k8s.io/klog/v2"
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	extensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForCRDEstablished is check the result if crd can be acceped and be serverd,
// for  apiextensions-apiserver will validate the name of crd, and print the result to status
func WaitForCRDEstablished(extClientSet extensionsclientset.Interface, crdName string) error {
	return wait.Poll(1250*time.Millisecond, 10*time.Second, func() (done bool, err error) {
		crd, err := extClientSet.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), crdName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1.NamesAccepted:
				//  检查资源定义名字是否满足一致性要求，并且是不是存在冲突
				if cond.Status == apiextensionsv1.ConditionFalse {
					klog.Error("CRD Name conflict")
				}
			case apiextensionsv1.Established:
				// ApiServer是否开始可以为定义的资源提供服务
				if cond.Status == apiextensionsv1.ConditionTrue {
					return true, nil
				}
			}
		}
		return false, err
	})
}
