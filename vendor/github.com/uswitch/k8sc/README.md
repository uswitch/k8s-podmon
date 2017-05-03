# Kubernetes Client

In-cluster and out-of-cluster client creation for [https://github.com/ericchiang/k8s](https://github.com/ericchiang/k8s).

`k8sc.NewClient` makes it easier to build programs that will connect running within a Kubernetes cluster and using `~/.kube/config` from a development machine. Calling `NewClient` with a path will load the configuration from that file.

```go
// this will connect using deserialized config
c, err := k8sc.NewClient("/home/me/.kube/config")

// this will connect using the in-cluster client
c, err := k8sc.NewClient("")
```

## License

Licensed under Apache License v2.0