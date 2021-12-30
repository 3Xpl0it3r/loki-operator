package install

import (
	"github.com/l0calh0st/loki-operator/pkg/crd"
	"github.com/l0calh0st/loki-operator/pkg/crd/loki"
	"github.com/l0calh0st/loki-operator/pkg/crd/promtail"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	extensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/util/errors"
)

func InstallCustomResourceDefineToApiServer(extClientSet extensionsclientset.Interface) error {
	cs := []*apiextensionsv1.CustomResourceDefinition{}
	cs = append(cs, loki.NewCustomResourceDefine())
	cs = append(cs, promtail.NewCustomResourceDefine())
	for _, c := range cs {
		err := crd.RegisterCRDWithObj(extClientSet, c)
		if err != nil {
			return err
		}
		if err = crd.WaitForCRDEstablished(extClientSet, c.GetName()); err != nil {
			return errors.NewAggregate([]error{err, crd.UnRegisterCRD(extClientSet, c.GetName())})
		}
	}
	return nil
}
