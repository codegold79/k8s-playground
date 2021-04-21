package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type RemoteClusterInfo struct {
	Host           string `yaml:"host"`
	ServiceAccount struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Token     string `yaml:"token"`
	}
}

// parseYAML retrieves remote cluster hosts, service account namespaces and tokens.
func parseYAML(yamlFile string) (RemoteClusterInfo, error) {
	source, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return RemoteClusterInfo{}, err
	}
	
	var config RemoteClusterInfo
	if err := yaml.Unmarshal(source, &config); err != nil {
		return RemoteClusterInfo{}, err
	}

	return config, nil
}
