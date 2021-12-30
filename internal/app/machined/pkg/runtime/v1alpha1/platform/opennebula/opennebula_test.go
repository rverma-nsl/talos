// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package opennebula_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/talos-systems/talos/internal/app/machined/pkg/runtime/v1alpha1/platform/opennebula"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
)

type ConfigSuite struct {
	suite.Suite
}

func (suite *ConfigSuite) TestNetworkConfig1() {
	cfg := []byte(`
ETH0_GATEWAY='164.52.208.1'
ETH0_IP='164.52.218.55'
ETH0_MAC='02:00:a4:34:da:37'
ETH0_MASK='255.255.240.0'
ETH0_NETWORK='164.52.208.0'
ETH0_VLAN_ID='461'
ETH1_IP='172.16.125.83'
ETH1_MAC='02:00:ac:10:7d:53'
ETH1_MASK='255.255.224.0'
	`)

	n := &opennebula.OpenNebula{}

	defaultMachineConfig := &v1alpha1.Config{}

	machineConfig := &v1alpha1.Config{
		MachineConfig: &v1alpha1.MachineConfig{
			MachineNetwork: &v1alpha1.NetworkConfig{
				NetworkHostname: "localhost",
				NameServers:     []string{"8.8.8.8", "8.8.4.4"},
				NetworkInterfaces: []*v1alpha1.Device{
					{
						DeviceInterface: "eth0",
						DeviceDHCP:      false,
						DeviceAddresses: []string{"164.52.218.55"},
						DeviceRoutes: []*v1alpha1.Route{
							{
								RouteNetwork: "/0",
								RouteGateway: "164.52.208.1",
							},
						},
					},
					{
						DeviceInterface: "eth1",
						DeviceDHCP:      false,
						DeviceAddresses: []string{"172.16.125.83"},
						DeviceRoutes: []*v1alpha1.Route{
							{
								RouteNetwork: "10.2.18.0/23",
								RouteGateway: "164.52.208.1",
							},
						},
					},
				},
			},
		},
	}

	result, err := n.ConfigurationNetwork(cfg, defaultMachineConfig)
	suite.Require().NoError(err)
	suite.Assert().Equal(machineConfig, result)
}
func (suite *ConfigSuite) TestNetworkConfig2() {
	cfg := []byte(`
ETH0_DNS='8.8.8.8 8.8.4.4'
ETH0_GATEWAY='10.2.18.1'
ETH0_IP='10.2.18.6'
ETH0_MAC='02:00:0a:02:12:06'
ETH0_MASK='255.255.254.0'
ETH0_NETWORK='10.2.18.0'
	`)

	n := &opennebula.OpenNebula{}

	defaultMachineConfig := &v1alpha1.Config{}

	machineConfig := &v1alpha1.Config{
		MachineConfig: &v1alpha1.MachineConfig{
			MachineNetwork: &v1alpha1.NetworkConfig{
				NetworkHostname: "localhost",
				NameServers:     []string{"8.8.8.8", "8.8.4.4"},
				NetworkInterfaces: []*v1alpha1.Device{
					{
						DeviceInterface: "eth0",
						DeviceMTU:       1450,
						DeviceDHCP:      false,
						DeviceAddresses: []string{"2000:0:100::/56"},
						DeviceRoutes: []*v1alpha1.Route{
							{
								RouteNetwork: "::/0",
								RouteGateway: "2000:0:100:2fff:ff:ff:ff:ff",
								RouteMetric:  1024,
							},
							{
								RouteNetwork: "2000:0:100:2f00::/58",
								RouteGateway: "2000:0:100:2fff:ff:ff:ff:f0",
								RouteMetric:  1024,
							},
						},
					},
					{
						DeviceInterface: "eth1",
						DeviceMTU:       9000,
						DeviceDHCP:      true,
						DeviceAddresses: []string{"2000:0:ff00::1/56"},
						DeviceRoutes: []*v1alpha1.Route{
							{
								RouteNetwork: "::/0",
								RouteGateway: "2000:0:ff00::",
								RouteMetric:  1024,
							},
						},
					},
				},
			},
		},
	}

	result, err := n.ConfigurationNetwork(cfg, defaultMachineConfig)

	suite.Require().NoError(err)
	suite.Assert().Equal(machineConfig, result)

}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
