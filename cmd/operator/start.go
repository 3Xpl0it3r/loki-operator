package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/l0calh0st/loki-operator/cmd/operator/options"
	"github.com/l0calh0st/loki-operator/pkg/apis/install"
	crclientset "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned"
	crinformers "github.com/l0calh0st/loki-operator/pkg/client/informers/externalversions"
	"github.com/l0calh0st/loki-operator/pkg/controller"
	"github.com/l0calh0st/loki-operator/pkg/controller/loki"
	"github.com/l0calh0st/loki-operator/pkg/controller/promtail"
	crdregister "github.com/l0calh0st/loki-operator/pkg/crd/install"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	apicorev1 "k8s.io/api/core/v1"
	extensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"
	"net"
	"net/http"

	"os"
	"time"
)

func NewStartCommand(stopCh <-chan struct{}) *cobra.Command {
	opts := options.NewOptions()
	cmd := &cobra.Command{
		Short: "Launch Loki-Operator",
		Long:  "Launch Loki-Operator",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return fmt.Errorf("Options validate failed, %v. ", err)
			}
			if err := opts.Complete(); err != nil {
				return fmt.Errorf("Options Complete failed %v. ", err)
			}
			if err := runCommand(opts, stopCh); err != nil {
				return fmt.Errorf("Run %s failed.: %v ", os.Args[0], err)
			}
			return nil
		},
	}
	fs := cmd.Flags()
	nfs := opts.NamedFlagSets()
	for _, f := range nfs.FlagSets {
		fs.AddFlagSet(f)
	}
	local := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(local)
	nfs.FlagSet("logging").AddGoFlagSet(local)

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), nfs, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), nfs, cols)
	})
	return cmd
}

func runCommand(o *options.Options, stopCh <-chan struct{}) error {
	install.Install(scheme.Scheme)
	var err error
	restConfig, err := buildKubeConfig("", "")
	if err != nil {
		return err
	}
	extClientSet, err := extensionsclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	if err = crdregister.InstallCustomResourceDefineToApiServer(extClientSet); err != nil {
		return err
	}

	kubeClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	crClientSet, err := crclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	crInformers := buildCustomResourceInformerFactory(crClientSet)
	kubeInformers := buildKubeStandardResourceInformerFactory(kubeClientSet)

	register := prometheus.NewRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)
	defer cancel()
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(nil, promhttp.HandlerOpts{}))
	svc := &http.Server{Handler: mux}
	l, err := net.Listen("tcp", o.ListenAddress)
	if err != nil {
		panic(err)
	}
	wg.Go(serve(svc, l))
	// build controller
	promtailControler := promtail.NewController(kubeClientSet, kubeInformers, crClientSet, crInformers, register)
	lokiController := loki.NewController(kubeClientSet, kubeInformers, crClientSet, crInformers, register)
	// run sharedInformer
	crInformers.Start(ctx.Done())
	kubeInformers.Start(ctx.Done())
	// run controller
	wg.Go(runController(ctx, promtailControler))
	wg.Go(runController(ctx, lokiController))
	select {
	case <-stopCh:
		klog.Infof("exited")
	}
	cancel()
	if err = wg.Wait(); err != nil {
		return err
	}
	return nil
}

func runController(ctx context.Context, controller controller.Controller) func() error {
	return func() error {
		if err := controller.Start(ctx); err != nil {
			return err
		}
		return nil
	}
}

func serve(srv *http.Server, listener net.Listener) func() error {
	return func() error {
		//level.Info(logger).Log("msg", "Starting insecure server on "+listener.Addr().String())
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}

func serveTLS(srv *http.Server, listener net.Listener) func() error {
	return func() error {
		//level.Info(logger).Log("msg", "Starting secure server on "+listener.Addr().String())
		if err := srv.ServeTLS(listener, "", ""); err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}

// buildKubeConfig build rest.Config from the following ways
// 1: path of kube_config 2: KUBECONFIG environment 3. ~/.kube/config, as kubeconfig may not in $HOMEDIR/.kube/
func buildKubeConfig(masterUrl, kubeConfig string) (*rest.Config, error) {
	cfgLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	cfgLoadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	cfgLoadingRules.ExplicitPath = kubeConfig
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(cfgLoadingRules, &clientcmd.ConfigOverrides{})
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	if err = rest.SetKubernetesDefaults(config); err != nil {
		return nil, err
	}
	return config, nil
}

// buildCustomResourceInformerFactory build crd informer factory according some options
func buildCustomResourceInformerFactory(crClient crclientset.Interface) crinformers.SharedInformerFactory {
	var factoryOpts []crinformers.SharedInformerOption
	factoryOpts = append(factoryOpts, crinformers.WithNamespace(apicorev1.NamespaceAll))
	factoryOpts = append(factoryOpts, crinformers.WithTweakListOptions(func(listOptions *v1.ListOptions) {
		// todo
	}))
	return crinformers.NewSharedInformerFactoryWithOptions(crClient, 5*time.Second, factoryOpts...)
}

// buildKubeStandardResourceInformerFactory build a kube informer factory according some options
func buildKubeStandardResourceInformerFactory(kubeClient kubernetes.Interface) informers.SharedInformerFactory {
	var factoryOpts []informers.SharedInformerOption
	factoryOpts = append(factoryOpts, informers.WithNamespace(apicorev1.NamespaceAll))
	//factoryOpts = append(factoryOpts, informers.WithCustomResyncConfig(nil))
	factoryOpts = append(factoryOpts, informers.WithTweakListOptions(func(listOptions *v1.ListOptions) {
		// todo
	}))
	return informers.NewSharedInformerFactoryWithOptions(kubeClient, 5*time.Second, factoryOpts...)
}
