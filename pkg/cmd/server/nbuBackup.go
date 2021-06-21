package server

import (
	velerov1apis "github.com/adi-bhardwaj/velero-modified/pkg/apis/velero/v1"
	"github.com/adi-bhardwaj/velero-modified/pkg/backup"
	"github.com/adi-bhardwaj/velero-modified/pkg/client"
	"github.com/adi-bhardwaj/velero-modified/pkg/controller"
	"github.com/adi-bhardwaj/velero-modified/pkg/persistence"
	"github.com/adi-bhardwaj/velero-modified/pkg/plugin/clientmgmt"
	"github.com/adi-bhardwaj/velero-modified/pkg/podexec"
	"github.com/sirupsen/logrus"
)

func ProcessNetBackupBackupRequest(factory client.Factory, vBackup *velerov1apis.Backup,
	mountPath string, log *logrus.Logger) error {
	serverConfig := getDefaultServerConfig()
	s, err := newServer(factory, serverConfig, true, log)
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

	kbClient, err := factory.KubebuilderClient()
	if err != nil {
		log.WithError(err).Error("failed to get kubebuilder client")
		return err
	}

	backupController := controller.NewBackupController(
		s.sharedInformerFactory.Velero().V1().Backups(),
		s.veleroClient.VeleroV1(),
		s.discoveryHelper,
		backupper,
		s.logger,
		s.logLevel,
		newPluginManager,
		backupTracker,
		kbClient,
		//s.mgr.GetClient(),
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

	return controller.ProcessNetBackupBackup(backupController, vBackup, mountPath, factory)
}
