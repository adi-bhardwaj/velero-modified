package netbackup

import (
	"github.com/sirupsen/logrus"
	pkgbackup "github.com/vmware-tanzu/velero/pkg/backup"
	"github.com/vmware-tanzu/velero/pkg/client"
	"github.com/vmware-tanzu/velero/pkg/cmd/server"
)

const (
	baseApplicationName = "nbukops"
)

func ProcessBackup(backup *pkgbackup.Request, mountPath string, log *logrus.Logger) error {
	config, err := client.LoadConfig()
	if err != nil {
		log.WithError(err).Error("failed to load the velero config")
		return err
	}

	f := client.NewFactory(baseApplicationName, config)

	return server.ProcessNetBackupBackupRequest(f, backup, mountPath, log)
}
