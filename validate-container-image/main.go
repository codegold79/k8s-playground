package main

import (
	// These two crypto imports are needed to pass validation. For more information read
	// https://github.com/opencontainers/go-digest#usage
	_ "crypto/sha256"
	_ "crypto/sha512"
	"fmt"

	dockerref "github.com/docker/distribution/reference"
)

func main() {
	images := []string{
		"velero",
		"velero/velero:v10.10.1",
		"velero:latest",
		"vel#ro:latest",      // 3
		"velero:lat?st",      // 4
		"vel#ro::lat?st",     // 5
		"vel#ro..lat?st",     // 6
		"velero._.lat3234st", // 7
		"docker.hub/velero:latest",
		"http://wwww.harbor.registry.com/velero:v1.7", // 9
		"wwww.harbor.registry.com/velero:v1.7",
		"https://wwww.harbor.registry.com/velero:v1.7", // 11
		"projects.registry.vmware.com/tanzu_migrator/velero:v0.0.7c",
		" projects.registry.vmware.com/tanzu_migrator/velero:v0.0.7c ",                                                         // 13
		"projects.registry.vmware.com/akamaipoc/nginx@sha256:34f3f875e745861ff8a37552ed7eb4b673544d2c56c7cc58f9a9bec5b4b3530e", // 14
		"foo@sha256:xyz", // 15
	}

	for i, img := range images {
		if err := validateImage(img); err != nil {
			fmt.Println(i, img, ":", err)
		}
	}
}

func validateImage(image string) error {
	// Kubernetes uses docker distribution reference to help with image parsing.
	// https://github.com/kubernetes/kubernetes/blob/c1e69551be1a72f0f8db6778f20658199d3a686d/pkg/kubelet/dockershim/libdocker/helpers.go
	// For my purpose, if an image can be parsed without error, the image is valid.
	_, err := dockerref.ParseNormalizedNamed(image)
	return err
}
