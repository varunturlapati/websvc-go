FROM golang:1.14 AS builder

WORKDIR $GOPATH/src/github.com/varunturlapati/websvc-go/
COPY . $GOPATH/src/github.com/varunturlapati/websvc-go/

EXPOSE 8081

CMD ["go", "run", "cmd/main.go"]
