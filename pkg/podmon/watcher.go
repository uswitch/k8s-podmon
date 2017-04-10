package podmon

import (
	"context"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ericchiang/k8s"
)

// Alert stuct is the message that will be passed when a failure is detected
type Alert struct {
	PodName           string
	Namespace         string
	Annotations       map[string]string
	Labels            map[string]string
	ContainerName     string
	ContainerExitCode int32
}

// HasKeyPrefix checks for a key starting with k
func HasKeyPrefix(m *map[string]string, k string) bool {
	for kk := range *m {
		if strings.HasPrefix(kk, k) {
			return true
		}
	}
	return false
}

// Watch a namespace for events
func Watch(ctx *context.Context, c *k8s.Client, namespace, annotation string, alertChan chan Alert) {
	if namespace == k8s.AllNamespaces {
		log.Infoln("Setting up watch on all namespaces")
	} else {
		log.Infof("Setting up watch on: %s", namespace)
	}

	wp, err := c.CoreV1().WatchPods(*ctx, namespace)
	if err != nil {
		log.Fatalf("Error getting watch: %s", err)
	}

	for {
		evt, pod, err := wp.Next()
		if err != nil {
			log.Warnf("Error getting event: %s", err)
			continue
		}

		if *evt.Type == "MODIFIED" && HasKeyPrefix(&pod.Metadata.Annotations, annotation) {
			for _, c := range pod.Status.ContainerStatuses {
				if c.State.Terminated != nil {
					if *c.State.Terminated.ExitCode > 0 {
						a := Alert{
							PodName:           *pod.Metadata.Name,
							Namespace:         *pod.Metadata.Namespace,
							Annotations:       pod.Metadata.Annotations,
							Labels:            pod.Metadata.Labels,
							ContainerName:     *c.Name,
							ContainerExitCode: *c.State.Terminated.ExitCode,
						}
						alertChan <- a
						log.Debugln("Event sent to notifier.")
					}
				}
			}
		}
	}
}
