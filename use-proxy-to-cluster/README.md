# Connect to a remote Kubernetes cluster using client-go via a proxy and with a kubeconfig file

This example also makes use of flags to pass in the kubeconfig file and proxy information.

Once you have a kubeconfig file and the optionally, the proxy information, you can build and run the program to list namespaces in the remote cluster.

```bash
go build -o connectWithProxy

./connectWithProxy -httpsproxy=http://<address>:<port> -kubeconfig=example-config
```
