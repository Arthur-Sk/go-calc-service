FROM golang:1.17.3-alpine3.14

RUN apk update && apk add protoc
RUN go install github.com/golang/protobuf/protoc-gen-go@latest # necessary to generate code from protocol buffers

COPY ./ /go/src/grpc-service
RUN chmod -R 777 /go/src

WORKDIR /go/src/grpc-service
RUN go mod download
RUN go build -o /go/bin/calc-service calculator/calc_server/server.go
RUN go build -o /go/bin/greet-service greet/greet_server/server.go

WORKDIR /go