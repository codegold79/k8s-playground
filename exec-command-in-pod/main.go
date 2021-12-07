package main

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const (
	ns        = "codegold79"
	pod       = "k8s"
	container = "playground"

	saTokenKey     = "CREDENTIAL_SA_TOKEN"
	clusterHostKey = "CREDENTIAL_CLUSTER_HOST"
)

func main() {
	creds, err := readClusterCredentials()
	if err != nil {
		fmt.Println("obtain cluster credentials:", err)
		return
	}

	config, err := restConfig(creds)
	if err != nil {
		fmt.Println("build REST config for client:", err)
		return
	}

	client, err := podRESTClient(config)
	if err != nil {
		fmt.Println("retrieve pod REST client:", err)
		return
	}

	execReq := client.
		Post().
		Namespace(ns).
		Resource("pods").
		Name(pod).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   []string{"ls", "var"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
		}, runtime.NewParameterCodec(scheme.Scheme))

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", execReq.URL())
	if err != nil {
		fmt.Println("create remote command executor:", err)
		return
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		fmt.Println("execute command in pod:", err)
	}
}

type clusterCredentials struct {
	saToken     string
	clusterHost string
}

func readClusterCredentials() (clusterCredentials, error) {
	var creds clusterCredentials

	host, ok := os.LookupEnv(clusterHostKey)
	if !ok {
		return clusterCredentials{}, fmt.Errorf("environment variable %q for cluster host is required", clusterHostKey)
	}
	creds.clusterHost = host

	token, ok := os.LookupEnv(saTokenKey)
	if !ok {
		return clusterCredentials{}, fmt.Errorf("environment variable %q for service account token is required", saTokenKey)
	}
	creds.saToken = token

	return creds, nil
}
func podRESTClient(config *rest.Config) (rest.Interface, error) {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}

	restClient, err := apiutil.RESTClientForGVK(gvk, false, config, serializer.NewCodecFactory(scheme.Scheme))
	if err != nil {
		return nil, fmt.Errorf("create REST client: %w", err)
	}

	return restClient, nil
}

func restConfig(creds clusterCredentials) (*rest.Config, error) {
	return &rest.Config{
		Host:            creds.clusterHost,
		BearerToken:     creds.saToken,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		Burst:           1000,
		QPS:             100,
	}, nil
}
