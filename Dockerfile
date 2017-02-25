FROM golang:1.7-alpine
COPY . /go/src/ider
RUN cd /go/src/ider && go build acceptance_tests/serverIderIntegration.go

ENTRYPOINT ["/go/src/ider/serverIderIntegration","/data"]
