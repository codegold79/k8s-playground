package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println("building client config from cluster credentials", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("create clientset from config", err)
		os.Exit(1)
	}

	secrets, err := clientset.CoreV1().Secrets("default").List(context.Background(), meta.ListOptions{})
	if err != nil {
		fmt.Println("list secrets", err)
		os.Exit(1)
	}

	for _, s := range secrets.Items {
		fmt.Println(s.Name)
		fmt.Println(s.Data)

		v, ok := s.Data["app"]
		if ok {
			// Secret data are stored in maps with string keys and slice of bytes values.
			fmt.Printf("data: %v, type: %[1]T, string: %s\n", v, string(v))
		}

		fmt.Println("======")
	}
}
