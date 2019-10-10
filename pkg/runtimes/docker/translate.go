/*
Copyright © 2019 Thorsten Klein <iwilltry42@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package docker

import (
	docker "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	k3d "github.com/rancher/k3d/pkg/types"
)

// TranslateNodeToContainer translates a k3d node specification to a docker container representation
func TranslateNodeToContainer(node *k3d.Node) (*NodeInDocker, error) {

	/* initialize everything that we need */
	containerConfig := docker.Config{}
	hostConfig := docker.HostConfig{}
	networkingConfig := network.NetworkingConfig{}

	/* Name & Image */
	containerConfig.Hostname = node.Name
	containerConfig.Image = node.Image

	/* Command & Arguments */
	containerConfig.Cmd = []string{}

	containerConfig.Cmd = append(containerConfig.Cmd, node.Cmd...)  // contains k3s command and role-specific required flags/args
	containerConfig.Cmd = append(containerConfig.Cmd, node.Args...) // extra flags/args

	/* Environment Variables */
	containerConfig.Env = node.Env

	/* Labels */
	containerConfig.Labels = node.Labels // has to include the role

	/* Auto-Restart */
	if node.Restart {
		hostConfig.RestartPolicy = docker.RestartPolicy{
			Name: "unless-stopped",
		}
	}

	/* Tmpfs Mounts */
	// TODO: do we need this or can the default be a map with empty values already?
	hostConfig.Tmpfs = make(map[string]string)
	for _, mnt := range k3d.DefaultTmpfsMounts {
		hostConfig.Tmpfs[mnt] = ""
	}

	/* They have to run in privileged mode */
	// TODO: can we replace this by a reduced set of capabilities?
	hostConfig.Privileged = true

	/* Volumes */
	// TODO: image volume
	hostConfig.Binds = node.Volumes // TODO: some validation?
	// containerConfig.Volumes = map[string]struct{}{} // TODO: do we need this? We only used binds before

	/* Ports */
	containerConfig.ExposedPorts = nat.PortSet{} // TODO: translate from node.Ports to nat.PortSet
	hostConfig.PortBindings = nat.PortMap{}      // TODO: this and exposedPorts required?

	/* Network */
	networkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
		node.Network: { // TODO: fill
			Aliases: []string{node.Name}, // TODO: fill
		},
	}

	return &NodeInDocker{
		ContainerConfig:  containerConfig,
		HostConfig:       hostConfig,
		NetworkingConfig: networkingConfig,
	}, nil
}
