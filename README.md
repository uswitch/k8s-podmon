# k8s-podmon

Watches one or all namespaces for Pods where a container terminates with a non-zero
exit code and if it has the annotation, will notify via a slack channel.

## Usage

    usage: k8s-podmon --slack=SLACK [<flags>]

    Flags:
          --help             Show context-sensitive help (also try --help-long and --help-man).
      -d, --debug            Debug output
          --kubecfg=KUBECFG  Location of kubeconfig, blank for In-Cluster
          --namespace=""     Namespace to follow
          --annotation="com.uswitch.alert/slack"  
                             Annotation to watch for
          --slack=SLACK      Slack webhook

To make a pod "monitored", set the annotation with a value of the slack channel you wish to spam.

For example:

    apiVersion: batch/v1
    kind: Job
    metadata:
    name: boom
    namespace: cloud
    spec:
    template:
    metadata:
      annotations:
        com.uswitch.alert/slack: kubernetes
    spec:
      containers:
      - name: hello
        image: busybox
        args:
        - /bin/sh
        - -c
        - echo Boom; exit 101
      restartPolicy: Never


## Building binary

    CGO_ENABLED=0 go build -o k8s-podmon cmd/*.go
