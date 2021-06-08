/*
	Use Velero as a library to back up and restore to Kubernetes clusters.
*/
package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/client"
	"github.com/vmware-tanzu/velero/pkg/discovery"
	clientset "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned"
	"github.com/vmware-tanzu/velero/pkg/podexec"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// TODO: make a timeout context.
	ctx := context.TODO()
	log := logrus.New()

	// TODO: Have a clusterCxn factory that can determine which clients to instantiate.
	clusterCxn, err := connectionComponents(ctx, log)
	if err != nil {
		log.WithField("event", "collect cluster connection components").Error(err)
	}

	// TODO: Use user input to determine which function to run.
	// TODO: Determine ahead of time which clusterCxn (remote, in-cluster, etc) to send
	if err := clusterCxn.backup(ctx, log); err != nil {
		log.WithFields(logrus.Fields{
			"event":              "backup Kubernetes workload",
			"cluster connection": clusterCxn.description,
		})
	}

	// TODO: Use user input to determine which function to run.
	// TODO: Determine ahead of time which clusterCxn (remote, in-cluster, etc) to send
	// if err := clusterCxn.restore(ctx, log); err != nil {
	// 	log.WithFields(logrus.Fields{
	// 		"event":              "restore Kubernetes workload",
	// 		"cluster connection": clusterCxn.description,
	// 	})
	// }
}

type clusterConnection struct {
	description        string
	veleroConfig       client.VeleroConfig
	veleroClient       clientset.Interface
	discoveryHelper    discovery.Helper
	dynamicFactory     client.DynamicFactory
	kubeClientConfig   *rest.Config
	kubeClient         kubernetes.Interface
	podCommandExecutor podexec.PodCommandExecutor
}

func connectionComponents(ctx context.Context, log *logrus.Logger) (clusterConnection, error) {
	veleroConfig := client.VeleroConfig{}

	factory := client.NewFactory("playground", veleroConfig)

	veleroClient, err := factory.Client()
	if err != nil {
		return clusterConnection{}, fmt.Errorf("create velero client: %w", err)
	}

	discoveryClient := veleroClient.Discovery()

	dynamicClient, err := factory.DynamicClient()
	if err != nil {
		return clusterConnection{}, fmt.Errorf("create dynamic client: %w", err)
	}

	discoveryHelper, err := discovery.NewHelper(discoveryClient, log)
	if err != nil {
		return clusterConnection{}, fmt.Errorf("create discovery helper: %w", err)
	}

	dynamicFactory := client.NewDynamicFactory(dynamicClient)

	kubeClient, err := factory.KubeClient()
	if err != nil {
		return clusterConnection{}, fmt.Errorf("create kubeclient: %w", err)
	}

	kubeClientConfig, err := factory.ClientConfig()
	if err != nil {
		return clusterConnection{}, fmt.Errorf("create kubeclient config: %w", err)
	}

	podCommandExecutor := podexec.NewPodCommandExecutor(
		kubeClientConfig,
		kubeClient.CoreV1().RESTClient(),
	)

	return clusterConnection{
		description:        "in-cluster Velero",
		veleroConfig:       veleroConfig,
		veleroClient:       veleroClient,
		discoveryHelper:    discoveryHelper,
		dynamicFactory:     dynamicFactory,
		kubeClientConfig:   kubeClientConfig,
		kubeClient:         kubeClient,
		podCommandExecutor: podCommandExecutor,
	}, nil
}
