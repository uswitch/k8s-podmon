package podmon

import (
	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/ericchiang/k8s"
	"github.com/ericchiang/k8s/api/v1"
)

func IsJob(p v1.Pod) bool {
	for k := range p.Metadata.Labels {
		if k == "job-name" {
			return true
		}
	}
	return false
}

// Watch a namespace for events
func Watch(ctx *context.Context, c *k8s.Client, namespace string) {
	log.Infof("Setting up watch on: %s", namespace)
	wp, err := c.CoreV1().WatchPods(*ctx, namespace)
	if err != nil {
		log.Fatalf("Error getting watch: %s", err)
	}

	for {
		evt, pod, err := wp.Next()
		if err != nil {
			log.Fatalf("Error getting event: %s", err)
		}
		if *evt.Type == "MODIFIED" {
			// l := pod.Metadata.Labels
			// a := pod.Metadata.Annotations

			log.Infof("Event: %s Pod:%s\n%+v", *pod.Metadata.Name, *pod.Status.Phase, pod.Status.ContainerStatuses)
		}
	}
}
