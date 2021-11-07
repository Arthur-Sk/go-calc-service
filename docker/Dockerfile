FROM golang:1.17.3-alpine3.14

COPY ./src /go/src/grpc-service
RUN go build -o bin/grpc-service src/grpc-service/main.go

CMD ["./bin/grpc-service"]