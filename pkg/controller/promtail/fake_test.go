package promtail_test

import (
	crfakeclients "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned/fake"
	crinformers "github.com/l0calh0st/loki-operator/pkg/client/informers/externalversions"
	"github.com/l0calh0st/loki-operator/pkg/controller"
	"github.com/l0calh0st/loki-operator/pkg/controller/promtail"
	"github.com/l0calh0st/loki-operator/pkg/operator/fake"
	"k8s.io/apimachinery/pkg/watch"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
	"time"
)

var noResyncPeriod = func()time.Duration {return 0}

type fakeController struct {
	controller controller.Controller

	fakeWatcher *watch.FakeWatcher
	crClient *crfakeclients.Clientset
	crInformer crinformers.SharedInformerFactory
}

func newFakeController() *fakeController {
	kubeClient := k8sfake.NewSimpleClientset()
	crClient := crfakeclients.NewSimpleClientset()
	watcher := watch.NewFakeWithChanSize(10, false)
	crClient.PrependWatchReactor("lokis", testing.DefaultWatchReactor(watcher, nil))
	crInformer := crinformers.NewSharedInformerFactory(crClient, noResyncPeriod())
	return &fakeController{
		controller:  promtail.NewFakeController(kubeClient, crClient, crInformer,fake.NewOperator()),
		fakeWatcher: watcher,
		crClient: crClient,
		crInformer: crInformer,
	}
}


