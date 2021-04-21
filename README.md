# k8s-playground

## Kubernetes Practice

Directories and Topics

- [list-remote-namespaces](list-remote-namespaces) List out namespaces in other Kubernetes clusters.
  - Practice with service account tokens (needed for a pod to have access to other clusters)
  - Use the go-client REST config to connect
  - Useful references
    - [konveyor/mig-controller](https://github.com/konveyor/mig-controller/blob/5ab8db550971b69cb772e2700e8213022b51bdfd/pkg/apis/migration/v1alpha1/migcluster_types.go#L258)
