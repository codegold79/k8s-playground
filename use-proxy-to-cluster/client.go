package main

import (
	"net/http"
	"net/url"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func restConfig(kubeconfig []byte, proxy string) (*rest.Config, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		transport := rt.(*http.Transport)
		proxyURL, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
		return transport
	})

	return config, nil
}

func remoteClient(kubeconfig []byte, proxy string) (kclient.Client, error) {
	config, err := restConfig(kubeconfig, proxy)
	if err != nil {
		return nil, err
	}

	client, err := kclient.New(config, kclient.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, err
	}

	return client, nil
}
