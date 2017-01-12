FROM golang:latest

MAINTAINER Jonathan Reyna "jreyna@juniper.net"

ENV REFRESHED_AT=2017-01-11 \
PREFIX=LSPTOOL

ADD . $GOPATH/src/github.com/WOWLABS/LSPTool
RUN go install -v github.com/WOWLABS/LSPTool/Server

WORKDIR $GOPATH/src/github.com/WOWLABS/LSPTool/Server

EXPOSE 9002

ENTRYPOINT ["/go/bin/Server"]
