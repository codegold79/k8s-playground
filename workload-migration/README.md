# Kubernetes Workload Migration using Velero as a Library

Velero is a back up and restore tool for Kubernetes clusters. Most of the back and restore functions are exported and can be used as external packages.

Here is an example of being able to back up workloads with the minimal Velero imports.

## Instructions

- Set up a Kubernetes cluster
- Before a backup, install an example application on the cluster in its own namespace
- Before a restore, delete the namespace being restored
- Create a directory called /tmp/velero-playground
- Run the application. For now, you need to comment out restore or backup, depending on which one you aren't doing
- A file containing the Velero backup files should be generated at /tmp/velero-playground/backup.tar.gz

## More TODOs not in code comments

- Include server plugins in order to backup cluster-scoped resources
