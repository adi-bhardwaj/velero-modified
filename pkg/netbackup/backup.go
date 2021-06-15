package netbackup

import (
	velerov1apis "github.com/adi-bhardwaj/velero-modified/pkg/apis/velero/v1"
	"github.com/adi-bhardwaj/velero-modified/pkg/client"
	"github.com/adi-bhardwaj/velero-modified/pkg/cmd/server"
	"github.com/sirupsen/logrus"
)

const (
	baseApplicationName = "nbukops"
)

func ProcessBackup(backup *velerov1apis.Backup, mountPath string, log *logrus.Logger) error {
	f , err := getFactoryForRequest()
	if err != nil {
		log.WithError(err).Error("failed to load the velero config")
		return err
	}

	return server.ProcessNetBackupBackupRequest(f, backup, mountPath, log)
}

func getFactoryForRequest() (client.Factory, error) {
	config, err := client.LoadConfig()
	if err != nil {
		return nil, err
	}

	f := client.NewFactory(baseApplicationName, config)
	return f, nil
}
