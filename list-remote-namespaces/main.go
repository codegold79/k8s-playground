package main

import (
	"context"
	"fmt"
	"os"

	core "k8s.io/api/core/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Provide pass at least one cluster filename as an the argument to run this application.")
		os.Exit(1)
	}

	files := os.Args[1:]
	clients := []kclient.Client{}

	for _, f := range files {
		clusterInfo, err := parseYAML(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client, err := remoteClient(clusterInfo)
		if err != nil {
			fmt.Printf("connect to src cluster at %s: %v\n", clusterInfo.Host, err)
			os.Exit(1)
		}

		clients = append(clients, client)
	}

	nsNames, err := listNamespacesInClusters(clients)
	if err != nil {
		fmt.Println(err)
	}

	for _, n := range nsNames {
		fmt.Println(n)
	}
}

func listNamespacesInClusters(clients []kclient.Client) ([]string, error) {
	nsNames := make([]string, 0)

	fmt.Printf("============ Namespaces from %d clusters =============\n", len(clients))
	for _, c := range clients {
		nsList := core.NamespaceList{}
		if err := c.List(context.Background(), &nsList); err != nil {
			return nil, err
		}

		for _, item := range nsList.Items {
			nsNames = append(nsNames, item.Name)
		}

	}
	return nsNames, nil
}
