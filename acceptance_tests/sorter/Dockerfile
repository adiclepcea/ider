FROM golang:1.7-alpine
COPY . /go/src/ider/acceptance_tests/sorter
RUN cd /go/src/ider && go build acceptance_tests/sorter/sorter.go

ENTRYPOINT ["/go/src/ider/sorter"]
