package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	var kubeconfigPath string
	var err error

	if len(os.Args) == 1 {
		kubeconfigPath, err = defaultKubeconfigPath()
		if err != nil {
			log.Fatalf("obtaining default kubeconfig: %v", err)
		}
	}

	if len(os.Args) > 1 {
		// Use provided kubeconfig.
		kubeconfigPath = os.Args[1]
	}

	fmt.Printf("Connecting to cluster using kubeconfig located at %s\n", kubeconfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalf("obtaining REST config from kubeconfig: %v", err)
	}

	// client-go way
	if err := clientsetExample(config); err != nil {
		log.Fatalf("demonstrating with clientset example: %v", err)
	}

	// controller-runtime/kubebuilder way
	if err := clusterClientExample(config); err != nil {
		log.Fatalf("demonstrating with cluster client example: %v", err)
	}
}

func defaultKubeconfigPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(user.HomeDir, ".kube", "config"), nil
}

func clientsetExample(config *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("obtaining clientset from REST config: %w", err)
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), meta.ListOptions{})
	if err != nil {
		return fmt.Errorf("obtaining pod list: %w", err)
	}
	fmt.Printf("Using client-go, %d pods were found in the cluster\n", len(pods.Items))

	return nil
}

func clusterClientExample(config *rest.Config) error {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	if err := core.AddToScheme(scheme); err != nil {
		return fmt.Errorf("add core to scheme: %w", err)
	}

	options := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
	}
	mgr, err := ctrl.NewManager(config, options)
	if err != nil {
		return fmt.Errorf("create new manager: %w", err)
	}
	go func() error {
		if err := mgr.Start(ctx); err != nil {
			return fmt.Errorf("start manager: %w", err)
		}
		return nil
	}()

	for {
		if mgr.GetCache().WaitForCacheSync(ctx) {
			break
		}
	}

	clnt := mgr.GetClient()
	var podList core.PodList
	if err := clnt.List(ctx, &podList); err != nil {
		return fmt.Errorf("list pods: %w", err)
	}
	fmt.Printf("Using controller-runtime, %d pods were found in the cluster\n", len(podList.Items))

	return nil
}
