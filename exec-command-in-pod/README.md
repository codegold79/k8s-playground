# Executing a command in a pod

How to programmatically exec into a pod and run a command by using the Kubernetes controller-runtime package.

Here is the commandline equivalent that I am demonstrating programmatically:

```bash
    kubectl exec -it pod-name -n pod-namespace -- /bin/bash
```

Reference:

- [1] https://github.com/kubernetes-sigs/kubebuilder/issues/803
