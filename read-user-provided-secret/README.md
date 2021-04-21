# Read a Secret

Example where  user creates a secret and it can be read by a process running in a pod.

Create a secret

```bash
kubectl apply -f secret.yaml
```

Run program and it will show the secret data.

```bash
go build -o readSecret

./readSecret
```
