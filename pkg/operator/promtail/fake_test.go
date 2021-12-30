package promtail_test

import (
	"time"

	"k8s.io/client-go/informers"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/record"

	crfakeclients "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned/fake"
	crinformers "github.com/l0calh0st/loki-operator/pkg/client/informers/externalversions"
	croperatortesting "github.com/l0calh0st/loki-operator/pkg/k8s/testing"
	"github.com/l0calh0st/loki-operator/pkg/operator"
	"github.com/l0calh0st/loki-operator/pkg/operator/promtail"
)

var (
	noResyncPeriod                = func() time.Duration { return 0 }
	defaultFakeRecorderBufferSize = 10
)

type fakeController struct {
	fixture *croperatortesting.Fixture

	// client
	kubeClient *k8sfake.Clientset
	crClient   *crfakeclients.Clientset
	// informers
	kubeInformerFactory informers.SharedInformerFactory
	crInformerFactory   crinformers.SharedInformerFactory
	// recorder
	recorder *record.FakeRecorder
	operator operator.Operator
}

// newFakeController return fakeController used to test operator
func newFakeController(kubeClient *k8sfake.Clientset, crClient *crfakeclients.Clientset) *fakeController {
	fc := &fakeController{
		kubeClient:          kubeClient,
		crClient:            crClient,
		kubeInformerFactory: informers.NewSharedInformerFactory(kubeClient, noResyncPeriod()),
		crInformerFactory:   crinformers.NewSharedInformerFactory(crClient, noResyncPeriod()),
		recorder:            record.NewFakeRecorder(defaultFakeRecorderBufferSize),
	}

	fc.fixture = croperatortesting.NewFixture(fc.kubeInformerFactory, fc.crInformerFactory)
	fc.operator = promtail.NewOperator(kubeClient, crClient, fc.crInformerFactory.Lokioperator().V1alpha1().Promtails().Lister(),
		fc.kubeInformerFactory.Core().V1().ConfigMaps().Lister(), fc.kubeInformerFactory.Core().V1().Services().Lister(),
		fc.kubeInformerFactory.Apps().V1().DaemonSets().Lister(), fc.recorder)
	return fc
}

func (fc *fakeController) runController(key interface{}, startInformer bool) error {
	if startInformer {
		stopCh := make(chan struct{})
		close(stopCh)
		fc.kubeInformerFactory.Start(stopCh)
		fc.crInformerFactory.Start(stopCh)
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
			action.Matches("list", "daemonsets") || action.Matches("watch", "daemonsets")) {
			continue
		}
		ret = append(ret, action)
	}
	return ret
}
