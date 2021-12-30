package loki

import (
	"context"
	"fmt"
	"github.com/l0calh0st/loki-operator/pkg/operator"
	"github.com/l0calh0st/loki-operator/pkg/operator/loki"
	listerappsv1 "k8s.io/client-go/listers/apps/v1"

	"time"

	"github.com/prometheus/client_golang/prometheus"
	apicorev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	crclientset "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned"
	crinformers "github.com/l0calh0st/loki-operator/pkg/client/informers/externalversions"
	crlisterv1alpha1 "github.com/l0calh0st/loki-operator/pkg/client/listers/lokioperator.l0calh0st.cn/v1alpha1"
	crcontroller "github.com/l0calh0st/loki-operator/pkg/controller"
)


var (
	alwaysReady = func() bool{return true}
	noResyncPeriodFunc = func() time.Duration{ return  0}
)

// controller is implement Controller for Loki resources
type controller struct {
	crcontroller.Base
	register      prometheus.Registerer
	kubeClientSet kubeclientset.Interface
	crClientSet   crclientset.Interface
	queue         workqueue.RateLimitingInterface
	recorder      record.EventRecorder

	statefulSetLister        listerappsv1.StatefulSetLister
	serviceLister    listercorev1.ServiceLister
	configMapLister  listercorev1.ConfigMapLister
	deploymentLister listerappsv1.DeploymentLister
	lokiLister       crlisterv1alpha1.LokiLister
	cacheSynced      []cache.InformerSynced

	operator operator.Operator
}

// NewFakeController return a new fake promTail controller
func NewFakeController(kubeClient kubeclientset.Interface,kubeInformerFactory informers.SharedInformerFactory,
	crClient crclientset.Interface, crInformerFactory crinformers.SharedInformerFactory,operator operator.Operator,
) crcontroller.Controller {
	c := NewController(kubeClient, kubeInformerFactory, crClient, crInformerFactory, nil).(*controller)
	c.operator = operator
	c.cacheSynced = []cache.InformerSynced{alwaysReady}
	c.recorder = record.NewFakeRecorder(10)
	return c
}


// NewController create a new controller for Loki resources
func NewController(kubeClientSet kubeclientset.Interface, kubeInformerFactory informers.SharedInformerFactory, crClientSet crclientset.Interface,
	crInformerFactory crinformers.SharedInformerFactory, reg prometheus.Registerer) crcontroller.Controller {

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.V(2).Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events(apicorev1.NamespaceAll)})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, apicorev1.EventSource{Component: "loki-operator"})


	return newLokiController(kubeClientSet, kubeInformerFactory, crClientSet, crInformerFactory, recorder, reg)
}



// newLokiController is really conv
func newLokiController(kubeClientSet kubeclientset.Interface, kubeInformers informers.SharedInformerFactory, crClientSet crclientset.Interface,
	crInformers crinformers.SharedInformerFactory, recorder record.EventRecorder, reg prometheus.Registerer) *controller {
	c := &controller{
		register:      reg,
		kubeClientSet: kubeClientSet,
		crClientSet:   crClientSet,
		recorder:      recorder,
	}
	c.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	statefulsetInformer := kubeInformers.Apps().V1().StatefulSets()
	c.statefulSetLister = statefulsetInformer.Lister()
	c.cacheSynced = append(c.cacheSynced, statefulsetInformer.Informer().HasSynced)

	deploymentInformer := kubeInformers.Apps().V1().Deployments()
	c.deploymentLister = deploymentInformer.Lister()
	c.cacheSynced = append(c.cacheSynced, deploymentInformer.Informer().HasSynced)

	serviceInformer := kubeInformers.Core().V1().Services()
	c.serviceLister = serviceInformer.Lister()
	c.cacheSynced = append(c.cacheSynced, serviceInformer.Informer().HasSynced)

	configMapInformer := kubeInformers.Core().V1().ConfigMaps()
	c.configMapLister = configMapInformer.Lister()
	c.cacheSynced = append(c.cacheSynced, configMapInformer.Informer().HasSynced)

	lokiInformer := crInformers.Lokioperator().V1alpha1().Lokis()
	c.lokiLister = lokiInformer.Lister()
	lokiInformer.Informer().AddEventHandlerWithResyncPeriod(newLokiEventHandler(c), 5*time.Second)
	c.cacheSynced = append(c.cacheSynced, lokiInformer.Informer().HasSynced)

	c.operator = loki.NewOperator(kubeClientSet, crClientSet,c.lokiLister, c.configMapLister, c.serviceLister, c.statefulSetLister, c.deploymentLister,c.recorder)
	return c
}

func (c *controller) Start(ctx context.Context) error {

	// wait for all involved cached to be synced , before processing items from the queue is started
	if !cache.WaitForCacheSync(ctx.Done(), func() bool {
		for _, hasSyncdFn := range c.cacheSynced {
			if !hasSyncdFn() {
				return false
			}
		}
		return true
	}) {
		return fmt.Errorf("timeout wait for cache to be synced")
	}
	klog.Infof("loki controller has started, begin handler items....")
	go wait.Until(c.runWorker, time.Second, ctx.Done())
	<-ctx.Done()
	return ctx.Err()
}

// runWorker for loop
func (c *controller) runWorker() {
	defer utilruntime.HandleCrash()
	for c.processNextItem() {
	}
}

func (c *controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)
	if err := c.operator.Reconcile(key); err != nil {
		// There was a failure so be sure to report it. This method allows for plugable error handling
		// which can be used for things like cluster-monitoring
		utilruntime.HandleError(fmt.Errorf("failed to reconcile lokioperator %q: %v", key, err))
		return true
	}
	// Successfully processed the key or the key was not found so tell the queue to stop tracking history for your key
	// This will reset things like failure counts for per-items rate limiting
	c.queue.Forget(key)
	return true
}

func (c *controller) Stop() {
	klog.Info("Stopping the loki operator controller")
	c.queue.ShutDown()
}

func (c *controller) enqueueFunc(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("failed to get key for %v: %v", obj, err)
		return
	}
	c.queue.AddRateLimited(key)
}
