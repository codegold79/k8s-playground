package main

import (
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func restConfig(clusterInfo RemoteClusterInfo) (*rest.Config, error) {
	return &rest.Config{
		Host:            clusterInfo.Host,
		BearerToken:     clusterInfo.ServiceAccount.Token,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		Burst:           1000,
		QPS:             100,
	}, nil
}

func remoteClient(info RemoteClusterInfo) (kclient.Client, error) {
	config, err := restConfig(info)
	if err != nil {
		return nil, err
	}

	client, err := kclient.New(config, kclient.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, err
	}

	return client, nil
}
