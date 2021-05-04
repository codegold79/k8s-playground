package main

import (
	"os"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func setProxy(url string) {
	os.Setenv("HTTPS_PROXY", url)
}

func restConfig(kubeconfig []byte) (*rest.Config, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	// The below gives the error, "using a custom transport with TLS certificate
	// options or the insecure flag is not allowed". As a workaround, I created
	// setProxy function to set it via environmental variable. I believe the
	// better way would be to use WrapTransport func(rt http.RoundTripper)
	// http.RoundTripper. The hint was provided at @GalloCedrone at
	// https://stackoverflow.com/questions/52218669/use-http-proxy-for-kubernetes-go-client.
	// TODO: figure out how to use WrapTransport.
	// proxyURL, err := url.Parse(info.Proxy) if err != nil {
	//  return nil, err
	// } config.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	return config, nil
}

func remoteClient(kubeconfig []byte) (kclient.Client, error) {
	config, err := restConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := kclient.New(config, kclient.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, err
	}

	return client, nil
}
