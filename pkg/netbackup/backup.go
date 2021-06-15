package netbackup

import (
	"github.com/sirupsen/logrus"
	pkgbackup "github.com/adi-bhardwaj/velero-modified/pkg/backup"
	"github.com/adi-bhardwaj/velero-modified/pkg/client"
	"github.com/adi-bhardwaj/velero-modified/pkg/cmd/server"
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
