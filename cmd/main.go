package main

import (
	"context"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/uswitch/k8s-podmon/pkg/podmon"
)

var (
	debug     = kingpin.Flag("debug", "Debug output").Short('d').Bool()
	namespace = kingpin.Flag("namespace", "Namespace to follow").Default("cloud").String()
)

func main() {
	kingpin.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Debug logging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	c, err := podmon.NewClient("/home/tom/.kube/config")
	if err != nil {
		log.Fatalf("Error starting k8s client: %s", err)
	}

	ctx := context.Background()

	podmon.Watch(&ctx, c, *namespace)
}
