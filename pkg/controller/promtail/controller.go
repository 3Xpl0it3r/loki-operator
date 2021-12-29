package promtail

import (
	"context"
	"fmt"
	"github.com/l0calh0st/loki-operator/pkg/operator"
	"github.com/l0calh0st/loki-operator/pkg/operator/promtail"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	apicorev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	listterappsv1 "k8s.io/client-go/listers/apps/v1"
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
	alwaysReady = func() bool { return true }
)

// controller is implement Controller for Loki resources
type controller struct {
	crcontroller.Base
	register      prometheus.Registerer
	kubeClientSet kubeclientset.Interface
	crClientSet   crclientset.Interface
	queue         workqueue.RateLimitingInterface

	// maintains those lister for defining custom handler
	daemonsetLister listterappsv1.DaemonSetLister
	serviceLister   listercorev1.ServiceLister
	promtailLister  crlisterv1alpha1.PromtailLister
	configMapLister listercorev1.ConfigMapLister
	cacheSynced     []cache.InformerSynced

	recorder record.EventRecorder

	operator operator.Operator
}

// NewController create a new controller for Loki resources
func NewController(kubeClientSet kubeclientset.Interface, kubeInformers informers.SharedInformerFactory, crClientSet crclientset.Interface,
	crInformers crinformers.SharedInformerFactory, reg prometheus.Registerer) crcontroller.Controller {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.V(2).Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events(apicorev1.NamespaceAll)})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, apicorev1.EventSource{Component: "loki-operator"})

	return newPromtailController(kubeClientSet, kubeInformers, crClientSet, crInformers, recorder, reg)
}

// NewFakeController return a new fake promTail controller
func NewFakeController(kubeClient kubeclientset.Interface, crClient crclientset.Interface, crInformers crinformers.SharedInformerFactory,operator operator.Operator) crcontroller.Controller {
	c :=  &controller{
		kubeClientSet: kubeClient,
		crClientSet:   crClient,
		register:      nil,
		recorder:      record.NewFakeRecorder(10),
		operator:      operator,
		cacheSynced:   []cache.InformerSynced{alwaysReady},
		queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}
	promtailInformer := crInformers.Lokioperator().V1alpha1().Promtails()
	promtailInformer.Informer().AddEventHandler(newPromtailEventHandler(c))

	return c
}

// newLokiController is really conv
func newPromtailController(kubeClientSet kubeclientset.Interface, kubeInformers informers.SharedInformerFactory, crClientSet crclientset.Interface,
	crInformers crinformers.SharedInformerFactory, recorder record.EventRecorder, reg prometheus.Registerer) *controller {
	c := &controller{
		register:      reg,
		kubeClientSet: kubeClientSet,
		crClientSet:   crClientSet,
		recorder:      recorder,
	}
	c.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	daemonsetInformer := kubeInformers.Apps().V1().DaemonSets()
	c.daemonsetLister = daemonsetInformer.Lister()
	c.cacheSynced = append(c.cacheSynced, daemonsetInformer.Informer().HasSynced)

	serviceInformer := kubeInformers.Core().V1().Services()
	c.serviceLister = serviceInformer.Lister()
	c.cacheSynced = append(c.cacheSynced, serviceInformer.Informer().HasSynced)

	configMapInformer := kubeInformers.Core().V1().ConfigMaps()
	c.configMapLister = configMapInformer.Lister()
	c.cacheSynced = append(c.cacheSynced, configMapInformer.Informer().HasSynced)

	promtailInformer := crInformers.Lokioperator().V1alpha1().Promtails()
	c.promtailLister = promtailInformer.Lister()
	promtailInformer.Informer().AddEventHandlerWithResyncPeriod(newPromtailEventHandler(c), time.Second)
	c.cacheSynced = append(c.cacheSynced, promtailInformer.Informer().HasSynced)

	c.operator = promtail.NewOperator(kubeClientSet, crClientSet, c.promtailLister, c.configMapLister, c.serviceLister, c.daemonsetLister, c.recorder)

	return c
}

// Start begin start controller, wait all informer synced and begin do reconcile
func (c *controller) Start(ctx context.Context) error {
	// wait for all involved cached to be synced , before processing items from the queue is started
	if !cache.WaitForNamedCacheSync("loki controller", ctx.Done(), func() bool {
		for _, cacheSynced := range c.cacheSynced {
			if !cacheSynced() {
				return false
			}
		}
		return true
	}) {
		return fmt.Errorf("timeout wait for cache to be synced")
	}

	klog.Info("promtail controller started. begin handler items......")
	go wait.Until(c.runWorker, time.Second, ctx.Done())
	return nil
}

// runWorker for loop
func (c *controller) runWorker() {
	defer utilruntime.HandleCrash()
	for c.processNextItem() {
		// for loop that never existed, until get key from workqueue failed
	}
}

// get item from workqueue to process
func (c *controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)
	klog.Infof("Starting process key: %q", key)
	if err := c.operator.Reconcile(key); err != nil {
		// There was a failure so be sure to report it. This method allows for plugable error handling
		// which can be used for things like cluster-monitoring
		c.queue.AddRateLimited(key)
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
	klog.Infof("add object %v\n", obj)
	if err != nil {
		klog.Errorf("failed to get key for %v: %v", obj, err)
		return
	}
	c.queue.AddRateLimited(key)
}
