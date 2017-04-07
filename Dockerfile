FROM scratch
LABEL maintainer Tom Taylor <tom.taylor@uswitch.com>

ADD k8s-podmon /
ENTRYPOINT ["/k8s-podmon"]
