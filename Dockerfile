FROM golang:1.17.3-alpine3.14

RUN apk update && apk add protoc
RUN go get github.com/golang/protobuf/protoc-gen-go

WORKDIR /go

COPY ./ /go/src/grpc-service
RUN chmod -R 777 /go/src
RUN go build -o bin/grpc-service src/grpc-service/main.go

CMD ["./bin/grpc-service"]