FROM golang:1.17.2

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src/github.com/l0calh0st/loki-operator
VOLUME /go/src/github.com/l0calh0st/loki-operator

WORKDIR /go/src/github.com/l0calh0st/loki-operator
