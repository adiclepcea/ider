FROM golang:1.7-alpine
COPY . /go/src/ider/acceptance_tests/client
RUN cd /go/src/ider && go build acceptance_tests/client/client.go

ENTRYPOINT ["/go/src/ider/client"]
