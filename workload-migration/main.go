/*
Use Velero as a library to back up and restore Kubernetes clusters.
*/
package main

import (
	"context"
	"errors"
	"flag"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// TODO: Have these passed in as flags
const (
	includedNamespaces         = "nginx-example"
	namespace                  = "playground"
	backupName                 = "backup"
	restoreName                = "restore"
	resourceTerminatingTimeout = time.Minute * 10
)

type backupConfig struct {
	name               string
	namespace          string
	includedNamespaces string
	filepath           string
}

type restoreConfig struct {
	name                       string
	namespace                  string
	includedNamespaces         string
	filepath                   string
	backupName                 string
	resourceTerminatingTimeout time.Duration
}

func main() {
	action := flag.String("a", "", "backup or restore")
	filepath := flag.String("p", "", "The location and file name of the backup file or the file to restore. Example: /Users/codegold/Documents/kubebackups/backup.tar.gz.")
	flag.Parse()

	*filepath = path.Clean(*filepath)

	// TODO: make a timeout context.
	ctx := context.TODO()
	log := logrus.New()

	if err := validateCommandArgs(*action, *filepath); err != nil {
		log.WithField("event", "validate commandline arguments").Error(err)
	}

	// TODO: Have a clusterCxn factory that can determine which client type to
	// instantiate.
	clusterCxn, err := connectionComponents(ctx, log)
	if err != nil {
		log.WithField("event", "collect cluster connection components").Error(err)
	}

	switch *action {
	case "backup":
		config := backupConfig{
			name:               backupName,
			namespace:          namespace,
			includedNamespaces: includedNamespaces,
			filepath:           *filepath,
		}

		// TODO: Generate the correct remote/in-cluster client with a connection
		// factory
		if err := clusterCxn.backup(ctx, log, config); err != nil {
			log.WithFields(logrus.Fields{
				"event":              "backup Kubernetes workload",
				"cluster connection": clusterCxn.description,
			})
		}
	case "restore":
		config := restoreConfig{
			name:                       restoreName,
			namespace:                  namespace,
			includedNamespaces:         "",
			filepath:                   *filepath,
			backupName:                 backupName,
			resourceTerminatingTimeout: resourceTerminatingTimeout,
		}

		// TODO: Generate the correct remote/in-cluster client with a connection
		// factory
		if err := clusterCxn.restore(ctx, log, config); err != nil {
			log.WithFields(logrus.Fields{
				"event":              "restore Kubernetes workload",
				"cluster connection": clusterCxn.description,
			})
		}
	}
}

func validateCommandArgs(action, filepath string) error {
	if action != "backup" && action != "restore" {
		return errors.New("action must be either backup or restore")
	}

	if filepath == "" {
		return errors.New("backup file location must not be empty")
	}

	return nil
}
