package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	velero "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/builder"
	"github.com/vmware-tanzu/velero/pkg/restore"
)

var defaultRestorePriorities = []string{
	"customresourcedefinitions",
	"namespaces",
	"storageclasses",
	"volumesnapshotclass.snapshot.storage.k8s.io",
	"volumesnapshotcontents.snapshot.storage.k8s.io",
	"volumesnapshots.snapshot.storage.k8s.io",
	"persistentvolumes",
	"persistentvolumeclaims",
	"secrets",
	"configmaps",
	"serviceaccounts",
	"limitranges",
	"pods",
	// we fully qualify replicasets.apps because prior to Kubernetes 1.16, replicasets also
	// existed in the extensions API group, but we back up replicasets from "apps" so we want
	// to ensure that we prioritize restoring from "apps" too, since this is how they're stored
	// in the backup.
	"replicasets.apps",
	"clusters.cluster.x-k8s.io",
	"clusterresourcesets.addons.cluster.x-k8s.io",
}

func (cxn clusterConnection) restore(ctx context.Context, log *logrus.Logger) error {
	builtBackup := builder.
		ForBackup(velero.DefaultNamespace, "playground").
		IncludedNamespaces("nginx-example").
		DefaultVolumesToRestic(false).
		Result()

	// TODO: Make names and namespaces configurable.
	builtRestore := builder.
		ForRestore(velero.DefaultNamespace, "playground").
		Backup("playground").
		Result()

	// TODO: Make backup file name and path configurable
	backupFile, err := os.Open("/tmp/velero-playground/backup-nginx.tar.gz")
	if err != nil {
		return fmt.Errorf("open backup file: %w", err)
	}

	request := restore.Request{
		Log:          log,
		Backup:       builtBackup,
		Restore:      builtRestore,
		BackupReader: backupFile,
	}

	restorer, err := restore.NewKubernetesRestorer(
		cxn.veleroClient.VeleroV1(),
		cxn.discoveryHelper,
		cxn.dynamicFactory,
		defaultRestorePriorities,
		cxn.kubeClient.CoreV1().Namespaces(),
		nil,
		0,
		10*time.Minute,
		log,
		cxn.podCommandExecutor,
		cxn.kubeClient.CoreV1().RESTClient(),
	)
	if err != nil {
		return fmt.Errorf("create restorer: %w", err)
	}

	restorer.Restore(request, nil, nil, nil)

	return nil
}
