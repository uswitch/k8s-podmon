package k8sc

import (
	"fmt"
	"github.com/ericchiang/k8s"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// Creates a Kubernetes API client. If no path is specified the client is assumed to be in-cluster
func NewClient(kubeconfigPath string) (*k8s.Client, error) {
	if kubeconfigPath == "" {
		return k8s.NewInClusterClient()
	}

	var data []byte

	data, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("kubeconfig read failed: %v", err)
	}

	var config k8s.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failure parsing config: %v", err)
	}

	return k8s.NewClient(&config)
}
