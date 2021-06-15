package server

import (
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/backup"
	"github.com/vmware-tanzu/velero/pkg/client"
	"github.com/vmware-tanzu/velero/pkg/controller"
	"github.com/vmware-tanzu/velero/pkg/persistence"
	"github.com/vmware-tanzu/velero/pkg/plugin/clientmgmt"
	"github.com/vmware-tanzu/velero/pkg/podexec"
)

func ProcessNetBackupBackupRequest(factory client.Factory, vBackup *backup.Request,
	mountPath string, log *logrus.Logger) error {
	serverConfig := getDefaultServerConfig()
	s, err := newServer(factory, serverConfig, log)
	if err != nil {
		log.WithError(err).Error("failed to create new velero server for request processing")
		return err
	}

	if err = s.initDiscoveryHelper(); err != nil {
		log.WithError(err).Error("failed to initialize the discovery helper")
		return err
	}

	backupTracker := controller.NewBackupTracker()
	backupper, err := backup.NewKubernetesBackupper(
		s.veleroClient.VeleroV1(),
		s.discoveryHelper,
		client.NewDynamicFactory(s.dynamicClient),
		podexec.NewPodCommandExecutor(s.kubeClientConfig, s.kubeClient.CoreV1().RESTClient()),
		s.resticManager,
		s.config.podVolumeOperationTimeout,
		s.config.defaultVolumesToRestic,
	)
	if err != nil {
		log.WithError(err).Error("failed to initialize the backupper for the backup request")
		return err
	}

	newPluginManager := func(logger logrus.FieldLogger) clientmgmt.Manager {
		return clientmgmt.NewManager(logger, s.logLevel, s.pluginRegistry)
	}

	backupStoreGetter := persistence.NewObjectBackupStoreGetter(s.credentialFileStore)

	csiVSLister, csiVSCLister := s.getCSISnapshotListers()

	backupController := controller.NewBackupController(
		s.sharedInformerFactory.Velero().V1().Backups(),
		s.veleroClient.VeleroV1(),
		s.discoveryHelper,
		backupper,
		s.logger,
		s.logLevel,
		newPluginManager,
		backupTracker,
		s.mgr.GetClient(),
		s.config.defaultBackupLocation,
		s.config.defaultVolumesToRestic,
		s.config.defaultBackupTTL,
		s.sharedInformerFactory.Velero().V1().VolumeSnapshotLocations().Lister(),
		s.config.defaultVolumeSnapshotLocations,
		s.metrics,
		s.config.formatFlag.Parse(),
		csiVSLister,
		csiVSCLister,
		backupStoreGetter,
	)

	return controller.ProcessNetBackupBackup(backupController, vBackup, mountPath)
}
