package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	core "k8s.io/api/core/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	httpsProxy := flag.String("httpsproxy", "", "Proxy to use with kubeconfig file")
	kubeConfigPath := flag.String("kubeconfig", "", "Location of the kubeconfig file")
	flag.Parse()

	if httpsProxy != nil {
		setProxy(*httpsProxy)
	}

	var (
		kubeconfig []byte
		err        error
	)
	if kubeConfigPath != nil {
		kubeconfig, err = ioutil.ReadFile(*kubeConfigPath)
	}

	client, err := remoteClient(kubeconfig)
	if err != nil {
		fmt.Printf("connect to src cluster using kubeconfig at %s: %v\n", kubeConfigPath, err)
		os.Exit(1)
	}

	nsNames, err := listNamespacesInCluster(client)
	if err != nil {
		fmt.Println(err)
	}

	for _, n := range nsNames {
		fmt.Println(n)
	}
}

func listNamespacesInCluster(client kclient.Client) ([]string, error) {
	nsNames := make([]string, 0)

	fmt.Println("============ Namespaces in cluster =============")
	nsList := core.NamespaceList{}
	if err := client.List(context.Background(), &nsList); err != nil {
		return nil, err
	}

	for _, item := range nsList.Items {
		nsNames = append(nsNames, item.Name)
	}

	return nsNames, nil
}
