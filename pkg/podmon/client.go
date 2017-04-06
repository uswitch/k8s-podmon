package podmon

import (
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"

	"github.com/ericchiang/k8s"
	"github.com/ghodss/yaml"
)

// NewClient returns a new k8s client
func NewClient(kubeconfigPath string) (*k8s.Client, error) {
	var err error

	if kubeconfigPath == "" {
		log.Debugln("Creating in-cluster client")
		return k8s.NewInClusterClient()
	}

	log.Debugf("Creating client from kubeconfig")
	var data []byte

	data, err = ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("Read kubeconfig failed: %v", err)
	}

	var config k8s.Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Bad YAML in kubeconfig: %v", err)
	}

	log.Printf("Created client using kubecfg at %s", kubeconfigPath)
	return k8s.NewClient(&config)
}
