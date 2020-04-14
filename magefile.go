// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

var (
	DockerImg = "varunturlapati/blog"
	DockerTag = "1.0"
	DockerFilePath = "docker/Dockerfile"
	BuildPath = "."
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Docker() error {
	fmt.Println("Building docker image...")
	return sh.Run("docker", "build", "-t", fmt.Sprintf("%s:%s", DockerImg, DockerTag), "-f", DockerFilePath, BuildPath)
}
