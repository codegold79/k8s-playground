package main

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	configmapName = "vmsizes"
	namespaceName = "default"
)

type specs struct {
	VCPU   int `yaml:"vCPU"`
	MemGiB int `yaml:"MemGiB"`
}

type vmSizes struct {
	Tags     map[string]specs
	Category string
}

func main() {
	config, err := RestKubeConfig()
	if err != nil {
		fmt.Printf("kubeconfig is required: %v\n", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("building Kubernetes clientset: %v\n", err)
		return
	}

	configMap, err := clientset.CoreV1().ConfigMaps(namespaceName).Get(context.TODO(), configmapName, meta.GetOptions{})
	if err != nil {
		fmt.Printf("reading %s configmap in %s namespace: %v\n", configmapName, namespaceName, err)
		return
	}
	for k, v := range configMap.Data {
		if k == "vmsizes.yaml" {
			var sizes vmSizes
			if err := yaml.Unmarshal([]byte(v), &sizes); err != nil {
				fmt.Printf("parse yaml data in configmap: %v\n", err)
			}
			fmt.Printf("sizes: %#v\n", sizes)
			return
		}
	}

	fmt.Println("Expected vmsizes configmap data was not found")
}
