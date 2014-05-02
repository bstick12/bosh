package infrastructure

import (
	bosherr "bosh/errors"
	boshdpresolv "bosh/infrastructure/devicepathresolver"
	boshplatform "bosh/platform"
	boshsettings "bosh/settings"
)

type openstackInfrastructure struct {
	metadataService    MetadataService
	registry           Registry
	platform           boshplatform.Platform
	devicePathResolver boshdpresolv.DevicePathResolver
}

func NewOpenstackInfrastructure(
	metadataService MetadataService,
	registry Registry,
	platform boshplatform.Platform,
	devicePathResolver boshdpresolv.DevicePathResolver,
) (inf openstackInfrastructure) {
	inf.metadataService = metadataService
	inf.registry = registry
	inf.platform = platform
	inf.devicePathResolver = devicePathResolver
	return
}

func (inf openstackInfrastructure) GetDevicePathResolver() boshdpresolv.DevicePathResolver {
	return inf.devicePathResolver
}

func (inf openstackInfrastructure) SetupSsh(username string) error {
	publicKey, err := inf.metadataService.GetPublicKey()
	if err != nil {
		return bosherr.WrapError(err, "Error getting public key")
	}

	return inf.platform.SetupSsh(publicKey, username)
}

func (inf openstackInfrastructure) GetSettings() (boshsettings.Settings, error) {
	settings, err := inf.registry.GetSettings()
	if err != nil {
		return settings, bosherr.WrapError(err, "Getting settings from registry")
	}

	return settings, nil
}

func (inf openstackInfrastructure) SetupNetworking(networks boshsettings.Networks) (err error) {
	return inf.platform.SetupDhcp(networks)
}

func (inf openstackInfrastructure) GetEphemeralDiskPath(devicePath string) (realPath string, found bool) {
	return inf.platform.NormalizeDiskPath(devicePath)
}
