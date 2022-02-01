# Parse ConfigMap YAML Data into Go Structure

Objective

- create Kubernetes client
- create configmap with data that we structured last time
- read ConfigMap
- parse yaml into Go struct
- print structured data

Procedure

- stand up a Kubernetes cluster
- create kubeconfig file and note the path
- create a file named "vmsizes.yaml" with the following contents:

```yaml
tags:
  small:
    vCPU: 2
    MemGiB: 2
  medium:
    vCPU: 2
    MemGiB: 4
  large:
    vCPU: 2
    MemGiB: 8
  xlarge:
    vCPU: 4
    MemGiB: 16
  2xlarge:
    vCPU: 8
    MemGiB: 32
category: "sizes" # Optional.
```

- create a configmap using data from the file you created:

```bash
kubectl create configmap vmsizes --from-file=vmsizes.yaml
```

- run program with

```bash
    go run *.go -kubeconfig /path/to/kubeconfig
```

Credits

- Client-go example (https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go)