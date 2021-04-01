# syntax = docker/dockerfile:1.0.0-experimental
# Build Container
FROM golang:1.15.2 as builder

ENV GO111MODULE on

WORKDIR /go/src/bitbucket.org/latonaio

COPY go.mod .

#RUN go get k8s.io/client-go@v0.20.4

RUN go mod download


COPY . .

RUN go build -o microservice-monitor .


# Runtime Container
FROM alpine:3.12

RUN apk add --no-cache libc6-compat tzdata

COPY --from=builder /go/src/bitbucket.org/latonaio/microservice-monitor .

CMD ["./microservice-monitor"]

