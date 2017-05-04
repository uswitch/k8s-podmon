package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	k8s "github.com/ericchiang/k8s"
	k8sc "github.com/uswitch/k8sc"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/uswitch/k8s-podmon/pkg/podmon"
)

var (
	debug          = kingpin.Flag("debug", "Debug output").Short('d').Bool()
	kubecfg        = kingpin.Flag("kubecfg", "Location of kubeconfig, blank for In-Cluster").String()
	namespace      = kingpin.Flag("namespace", "Namespace to follow").Default(k8s.AllNamespaces).String()
	baseAnnotation = kingpin.Flag("annotation", "Base Annotation to watch for").Default("com.uswitch.alert").String()
	slack          = kingpin.Flag("slack", "Slack webhook").Envar("SLACK").Required().String()
	awsRegion      = kingpin.Flag("aws-region", "AWS Region").Envar("AWS_REGION").Default("eu-west-1").String()
)

func main() {
	kingpin.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Debug logging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	defer cancel()

	// k8s client
	client, err := k8sc.NewClient(*kubecfg)
	if err != nil {
		log.Fatalf("Error starting k8s client: %s", err)
	}

	// Fire off watcher
	alertChan := make(chan podmon.Alert, 10)
	go podmon.Watch(&ctx, client, *namespace, *baseAnnotation, alertChan)

	// Fire off slacker
	slackChan := make(chan podmon.SlackMessage, 5)
	slacker, err := podmon.NewSlackEndpoint(*slack)
	if err != nil {
		log.Fatalf("Error with Slack Webhook: %s", err)
	}
	go slacker.EventLoop(ctx, &wg, slackChan)
	wg.Add(1)

	// Fire off SNS publisher
	snsChan := make(chan podmon.SNSMessage, 5)
	snsEP := podmon.NewSNSEndpoint(awsRegion)
	go snsEP.EventLoop(ctx, &wg, snsChan)
	wg.Add(1)

	// Alert consumer loop
	go func() {
		for {
			a := <-alertChan
			log.Debugf("Got alert: %+v", a)

			for k, v := range a.Annotations {
				if strings.HasPrefix(k, *baseAnnotation) {
					alertType := strings.Split(k, "/")
					if len(alertType) != 2 {
						log.Errorf("Annotation should be of type %s/<something>, I got %s", *baseAnnotation, k)
						continue
					}

					switch alertType[1] {

					case "slack":
						slackChan <- podmon.SlackMessage{
							Text:     fmt.Sprintf("Pod `%s.%s` (container: `%s`) has failed with exit code: `%d`", a.Namespace, a.PodName, a.ContainerName, a.ContainerExitCode),
							Icon:     ":poop:",
							Username: "k8s-PodMon",
							Channel:  v,
						}

					case "sns":
						snsChan <- podmon.SNSMessage{
							Subject:  fmt.Sprintf("Pod %s.%s (container: %s) has failed with exit code: %d", a.Namespace, a.PodName, a.ContainerName, a.ContainerExitCode),
							Message:  fmt.Sprintf("Pod %s.%s (container: %s) has failed with exit code: %d", a.Namespace, a.PodName, a.ContainerName, a.ContainerExitCode),
							TopicARN: v,
						}

					default:
						log.Errorf("Alert type %s not implimented", alertType[1])
						continue
					}
				}
			}
		}
	}()
	wg.Wait()
}
