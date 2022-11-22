## build
FROM golang:1.18-alpine AS build-env

ADD . /go/src/pontus

WORKDIR /go/src/pontus

RUN export GOPROXY=https://goproxy.cn && go build -v -mod=vendor

## run
FROM alpine:3.9

LABEL maintainer="colynn.liu <colynn.liu@gmail.com>"

ADD config /pontus/config

ADD assets /pontus/assets

# ADD api/swagger-spec /pontus/api/swagger-spec

RUN mkdir -p /pontus && mkdir -p /pontus/logs && touch /pontus/logs/pontus.log

WORKDIR /pontus

COPY --from=build-env /go/src/pontus/pontus /pontus

ENV PATH $PATH:/pontus

EXPOSE 8000
CMD ["./pontus"]