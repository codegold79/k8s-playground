# Connect to a Kubernetes Cluster programmatically using client-go 

## Run program

- Start a Kubernetes cluster
- The example program will use the default kubeconfig at `/home/.kube/config`, however, a different kubeconfig path can be specified as an argument.

Example command to run program:

`go run main.go /path/to/kubeconfig`

## References

- [Access a client using go-client](https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/#go-client)
- [Get a user's home directory](https://golangbyexample.com/get-current-user-home-directory-go/)
- [The difference between using client-go and controller-runtime](https://vivilearns2code.github.io/k8s/2021/03/12/diving-into-controller-runtime.html)