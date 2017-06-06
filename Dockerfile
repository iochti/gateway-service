FROM golang:1.8
MAINTAINER Luc CHMIELOWSKI <luc.chmielowski@gmail.com>

# Envs
ENV GO15VENDOREXPERIMENT=1

EXPOSE 3000

RUN mkdir -p /go/src/github.com/iochti/gateway-service
WORKDIR /go/src/github.com/iochti/gateway-service
COPY . /go/src/github.com/iochti/gateway-service

RUN go get github.com/tools/godep
RUN godep restore
RUN go install ./...

RUN rm -rf /go/src/github.com/iochti/gateway-service

CMD ["gateway-service"]
