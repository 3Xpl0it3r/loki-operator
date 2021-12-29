package loki_test

import (
	crfakeclients "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned/fake"
	crinformers "github.com/l0calh0st/loki-operator/pkg/client/informers/externalversions"
	croperatortesting "github.com/l0calh0st/loki-operator/pkg/k8s/testing"
	operatortesting "github.com/l0calh0st/loki-operator/pkg/k8s/testing"
	"github.com/l0calh0st/loki-operator/pkg/operator"
	"github.com/l0calh0st/loki-operator/pkg/operator/loki"
	"k8s.io/client-go/informers"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/record"
	"time"
)

var (
	noResyncPeriod = func() time.Duration { return 0 }
)

type fakeController struct {
	fixture *operatortesting.Fixture

	// client
	kubeClient *k8sfake.Clientset
	crClient   *crfakeclients.Clientset
	// informers
	kubeInformers informers.SharedInformerFactory
	crInformers   crinformers.SharedInformerFactory
	// recorder
	recorder *record.FakeRecorder
	operator operator.Operator
}

// newFakeController return fakeController used to test operator
func newFakeController(kubeClient *k8sfake.Clientset, crClient *crfakeclients.Clientset) *fakeController {
	fc := &fakeController{
		kubeClient:    kubeClient,
		crClient:      crClient,
		kubeInformers: informers.NewSharedInformerFactory(kubeClient, noResyncPeriod()),
		crInformers:   crinformers.NewSharedInformerFactory(crClient, noResyncPeriod()),
		recorder:      record.NewFakeRecorder(10),
	}

	fc.fixture = operatortesting.NewFixture(fc.kubeInformers, fc.crInformers)
	fc.operator = loki.NewOperator(kubeClient, crClient, fc.crInformers.Lokioperator().V1alpha1().Lokis().Lister(),
		fc.kubeInformers.Core().V1().ConfigMaps().Lister(), fc.kubeInformers.Core().V1().Services().Lister(),
		fc.kubeInformers.Apps().V1().StatefulSets().Lister(), fc.kubeInformers.Apps().V1().Deployments().Lister(), fc.recorder)
	return fc
}

func (fc *fakeController) runController(key string, startInformer bool) error {
	if startInformer {
		stopCh := make(chan struct{})
		close(stopCh)
		fc.kubeInformers.Start(stopCh)
		fc.crInformers.Start(stopCh)
	}
	// run controller
	if err := fc.operator.Reconcile(key); err != nil {
		return err
	}

	// validate customResource actions that recorded by customResourceClient
	if err := croperatortesting.ValidateActions(fc.fixture.GetCustomResourceActions(), informationActionFilter(fc.crClient.Actions())); err != nil {
		return err
	}

	// validate kube object actions that recorded by kube client
	if err := croperatortesting.ValidateActions(fc.fixture.GetKubeActions(), informationActionFilter(fc.kubeClient.Actions())); err != nil {
		return err
	}
	return nil
}

// filterInformerActions return actions according filer
func informationActionFilter(actions []k8stest.Action) []k8stest.Action {
	ret := make([]k8stest.Action, 0)
	for _, action := range actions {
		if len(action.GetNamespace()) == 0 && (action.Matches("list", "promtails") || action.Matches("watch", "promtails") ||
			action.Matches("list", "services") || action.Matches("watch", "services") ||
			action.Matches("list", "configmaps") || action.Matches("watch", "configmaps") ||
			action.Matches("list", "deployments") || action.Matches("watch", "deployments") ||
			action.Matches("list", "statefulsets") || action.Matches("watch", "statefulsets")) {
			continue
		}
		ret = append(ret, action)
	}
	return ret
}
