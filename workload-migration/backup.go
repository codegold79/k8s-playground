package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	velero "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/backup"
	"github.com/vmware-tanzu/velero/pkg/builder"
)

func (cxn clusterConnection) backup(ctx context.Context, log *logrus.Logger) error {
	builtBackup := builder.
		ForBackup(velero.DefaultNamespace, "playground").
		IncludedNamespaces("nginx-example").
		DefaultVolumesToRestic(false).
		Result()

	request := backup.Request{
		Backup: builtBackup,
	}

	// TODO: Make names, namespaces, backup file to restore configurable.
	os.MkdirAll("/tmp/velero-playground", 0644)
	backupFile, err := os.Create("/tmp/velero-playground/backup-nginx.tar.gz")
	if err != nil {
		return fmt.Errorf("create backup file: %w", err)
	}
	defer backupFile.Close()

	backupper, err := backup.NewKubernetesBackupper(
		cxn.veleroClient.VeleroV1(),
		cxn.discoveryHelper,
		cxn.dynamicFactory,
		cxn.podCommandExecutor,
		nil,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("create backupper: %w", err)
	}

	backupper.Backup(log, &request, backupFile, nil, nil)

	return nil
}
