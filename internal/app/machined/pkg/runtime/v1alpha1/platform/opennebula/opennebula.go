// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package opennebula

import (
	"context"
	"fmt"
	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"

	"github.com/talos-systems/go-blockdevice/blockdevice/filesystem"
	"github.com/talos-systems/go-blockdevice/blockdevice/probe"
	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/talos/internal/app/machined/pkg/runtime"
	"github.com/talos-systems/talos/internal/app/machined/pkg/runtime/v1alpha1/platform/errors"
	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/constants"

	"golang.org/x/sys/unix"
)

const (
	mnt            = "/mnt"
	configISOLabel = "CONTEXT"
	configDataPath = "context.sh"
)

// OpenNebula is the concrete type that implements the platform.Platform interface.
type OpenNebula struct{}

// Name implements the platform.Platform interface.
func (n *OpenNebula) Name() string {
	return "OpenNebula"
}

// ConfigurationNetwork implements the network configuration interface.
//nolint:gocyclo
func (n *OpenNebula) ConfigurationNetwork(metadataConfig []byte, confProvider config.Provider) (config.Provider, error) {
	var machineConfig *v1alpha1.Config

	machineConfig, ok := confProvider.(*v1alpha1.Config)
	if !ok {
		return nil, fmt.Errorf("unable to determine machine config type")
	}

	if machineConfig.MachineConfig == nil {
		machineConfig.MachineConfig = &v1alpha1.MachineConfig{}
	}

	if machineConfig.MachineConfig.MachineNetwork == nil {
		machineConfig.MachineConfig.MachineNetwork = &v1alpha1.NetworkConfig{}
	}

	vmConfig := getContextToMap(metadataConfig)

	if machineConfig.MachineConfig.MachineNetwork.NetworkInterfaces == nil {
		networkInterfaces := []*v1alpha1.Device{
			{
				DeviceInterface: "eth0",
				DeviceDHCP:      false,
				DeviceAddresses: []string{vmConfig["ETH0_IP"]},
				DeviceRoutes: []*v1alpha1.Route{
					{
						RouteNetwork: "0.0.0.0/0",
						RouteGateway: vmConfig["ETH0_GATEWAY"],
						RouteMetric:  1024,
					},
				},
			},
		}

		if _, ok := vmConfig["ETH1_IP"]; ok {
			netmask := net.ParseIP(vmConfig["ETH1_MASK"])
			sz, _ := net.IPMask(netmask.To4()).Size()
			networkInterfaces = append(networkInterfaces, &v1alpha1.Device{
				DeviceInterface: "eth1",
				DeviceDHCP:      false,
				DeviceAddresses: []string{vmConfig["ETH1_IP"]},
				DeviceRoutes: []*v1alpha1.Route{
					{
						RouteNetwork: fmt.Sprintf("%s/%d", vmConfig["ETH1_IP"], sz),
						RouteGateway: vmConfig["ETH0_GATEWAY"],
						RouteMetric:  1024,
					},
				},
			})
		}
		machineConfig = &v1alpha1.Config{MachineConfig: &v1alpha1.MachineConfig{
			MachineNetwork: &v1alpha1.NetworkConfig{NetworkInterfaces: networkInterfaces},
		},
		}
	}

	return confProvider, nil
}

// Configuration implements the platform.Platform interface.
//nolint:gocyclo
func (n *OpenNebula) Configuration(context.Context) ([]byte, error) {
	var option *string
	if option = procfs.ProcCmdline().Get(constants.KernelParamConfig).First(); option == nil {
		return nil, fmt.Errorf("%s not found", constants.KernelParamConfig)
	}

	log.Printf("fetching machine config from nebula cdrom mount")
	vmContext, err := n.configFromCD()
	if err != nil {
		return nil, err
	}

	emptyTalosConfig := &v1alpha1.Config{v1alpha1.Version, false, false, nil, nil}

	emptyConfig, err := emptyTalosConfig.Bytes()
	if err != nil {
		return nil, err
	}

	confProvider, err := configloader.NewFromBytes(emptyConfig)
	if err != nil {
		return nil, err
	}

	confProvider, err = n.ConfigurationNetwork(vmContext, confProvider)
	if err != nil {
		return nil, err
	}

	return confProvider.Bytes()
}

// Hostname implements the platform.Platform interface.
func (n *OpenNebula) Hostname(context.Context) (hostname []byte, err error) {
	return nil, nil
}

// Mode implements the platform.Platform interface.
func (n *OpenNebula) Mode() runtime.Mode {
	return runtime.ModeCloud
}

// ExternalIPs implements the runtime.Platform interface.
func (n *OpenNebula) ExternalIPs(context.Context) (addrs []net.IP, err error) {
	return addrs, err
}

// KernelArgs implements the runtime.Platform interface.
func (n *OpenNebula) KernelArgs() procfs.Parameters {
	return []*procfs.Parameter{
		procfs.NewParameter("console").Append("tty0").Append("ttyS0"),
	}
}

//nolint:gocyclo
func (n *OpenNebula) configFromCD() (vmContext []byte, err error) {
	var dev *probe.ProbedBlockDevice

	dev, err = probe.GetDevWithFileSystemLabel(strings.ToLower(configISOLabel))
	if err != nil {
		dev, err = probe.GetDevWithFileSystemLabel(strings.ToUpper(configISOLabel))
		if err != nil {
			return nil, errors.ErrNoConfigSource
		}
	}

	//nolint:errcheck
	defer dev.Close()

	sb, err := filesystem.Probe(dev.Path)
	if err != nil || sb == nil {
		return nil, errors.ErrNoConfigSource
	}

	log.Printf("found config disk (cidata) at %s", dev.Path)

	if err = unix.Mount(dev.Path, mnt, sb.Type(), unix.MS_RDONLY, ""); err != nil {
		return nil, errors.ErrNoConfigSource
	}

	vmContext, err = ioutil.ReadFile(filepath.Join(mnt, configDataPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read: %s", configDataPath)
	}

	if err = unix.Unmount(mnt, 0); err != nil {
		return nil, fmt.Errorf("failed to unmount: %w", err)
	}

	if vmContext == nil || len(vmContext) == 0 {
		return nil, errors.ErrNoConfigSource
	}

	return vmContext, nil
}

func getContextToMap(vmContext []byte) map[string]string {
	entries := strings.Split(string(vmContext), "\n")
	vmConfig := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, "=")
		vmConfig[parts[0]] = parts[1]
	}
	return vmConfig
}
