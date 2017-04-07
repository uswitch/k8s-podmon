package main

import (
	"context"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ericchiang/k8s"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/uswitch/k8s-podmon/pkg/podmon"
)

var (
	debug      = kingpin.Flag("debug", "Debug output").Short('d').Bool()
	kubecfg    = kingpin.Flag("kubecfg", "Location of kubeconfig, blank for In-Cluster").String()
	namespace  = kingpin.Flag("namespace", "Namespace to follow").Default(k8s.AllNamespaces).String()
	annotation = kingpin.Flag("annotation", "Annotation to watch for").Default("com.uswitch.alert/slack").String()
	slack      = kingpin.Flag("slack", "Slack webhook").Envar("SLACK").Required().String()
)

func main() {
	kingpin.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Debug logging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// k8s client
	client, err := podmon.NewClient(*kubecfg)
	if err != nil {
		log.Fatalf("Error starting k8s client: %s", err)
	}
	// Slack endpoint
	slacker, err := podmon.NewSlackEndpoint(*slack)
	if err != nil {
		log.Fatalf("Error with Slack Webhook: %s", err)
	}

	ctx := context.Background()
	// Chan to return alerts
	alertChan := make(chan podmon.Alert)

	// Fire off watcher
	go podmon.Watch(&ctx, client, *namespace, *annotation, alertChan)

	// Alert consumer loop
	for {
		a := <-alertChan
		log.Debugf("Got alert: %+v", a)
		// Should have annotation but just in case ...
		if podmon.HasKeyPrefix(&a.Annotations, *annotation) {
			m := podmon.SlackMessage{
				Text:     fmt.Sprintf("Pod `%s.%s` (container: `%s`) has failed with exit code: `%d`", a.Namespace, a.PodName, a.ContainerName, a.ContainerExitCode),
				Icon:     ":poop:",
				Username: "k8s-PodMon",
				Channel:  a.Annotations[*annotation],
			}
			resp, err := slacker.Send(m)
			if err != nil {
				log.Warnf("Slack error: %s", err)
			} else {
				log.Debugf("Got a %d from sending the following to slack: %#v", resp, m)
			}
		}
	}
}
