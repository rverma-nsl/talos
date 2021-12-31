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
ETH0_DNS='8.8.8.8 1.1.1.1 1.0.0.1'
ETH0_GATEWAY='164.52.208.1'
ETH0_IP='164.52.212.8'
ETH0_MAC='02:00:a4:34:d4:08'
ETH0_MASK='255.255.240.0'
ETH0_NETWORK='164.52.208.0'
ETH1_DNS='8.8.8.8 8.8.4.4'
ETH1_GATEWAY='10.2.18.1'
ETH1_IP='10.2.18.9'
ETH1_MAC='02:00:0a:02:12:09'
ETH1_MASK='255.255.254.0'
ETH1_NETWORK='10.2.18.0'
SSH_PUBLIC_KEY='ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDONRZFKq/IKdIbQYiaAHooP1M1o7UKxLfEGYrYYVIEhorwymmtJVwxK2qHbESaWrjDSdxRuAI5J6x+9Ruvxbha0H7XmQdbHLzkXKJcyeW44hAgMxRRMv9sMLC6CWp7XJOtQlzR/96a9KFSd7+ijQPGJxLCFKFlJgu+Sf5OlBYpP/3z9BW1fpC+PLoxcxxHeaDbvCiRcpTjd+wsBKGbm635heCewZ3Hf4lHdQO5+5WuLNud4mjUHw1uYEyACD0TrbgjIFk/tX7OkUO7t7M0k92ugFL34f7nsFozXNDnD5xFbxsWKhDJbNRbGvClxkHs68ZoMHjR003afg3wC9QSa7Rrsy0KpkJsxjbYa0KQ37aA3ifHZHqTtmEDLhpQP7G17p/Ts+l/874fpEJQei7BQP9LEugpYAQUjo4DcugE7+WAFZy9F6hJgrLougcQEbK3DIkpVsKt6Ue8b3N9P720= rverma@rverma.local'
	`)

	n := &opennebula.OpenNebula{}

	defaultMachineConfig := &v1alpha1.Config{}

	machineConfig := &v1alpha1.Config{
		MachineConfig: &v1alpha1.MachineConfig{
			MachineNetwork: &v1alpha1.NetworkConfig{
				NetworkHostname: "localhost",
				NameServers:     []string{"8.8.8.8", "1.1.1.1", "1.0.0.1"},
				NetworkInterfaces: []*v1alpha1.Device{
					{
						DeviceInterface: "eth0",
						DeviceDHCP:      false,
						DeviceAddresses: []string{"164.52.212.8"},
						DeviceCIDR:      "164.52.208.0/20",
						DeviceRoutes: []*v1alpha1.Route{
							{
								RouteNetwork: "0.0.0.0/0",
								RouteGateway: "164.52.208.1",
							},
						},
					},
					{
						DeviceInterface: "eth1",
						DeviceDHCP:      false,
						DeviceAddresses: []string{"10.2.18.9"},
						DeviceCIDR:      "10.2.18.0/23",
						DeviceRoutes: []*v1alpha1.Route{
							{
								RouteNetwork: "10.2.18.0/23",
								RouteGateway: "10.2.18.1",
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
						DeviceDHCP:      false,
						DeviceAddresses: []string{"10.2.18.6"},
						DeviceCIDR:      "10.2.18.0/23",
						DeviceRoutes: []*v1alpha1.Route{
							{
								RouteNetwork: "0.0.0.0/0",
								RouteGateway: "10.2.18.1",
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
