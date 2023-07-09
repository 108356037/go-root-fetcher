FROM golang:1.20 AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download 

COPY . .
RUN go build -ldflags '-w -extldflags "-static"' -o main

FROM alpine:3.18 AS base
RUN apk update
#RUN apt-get install -y vim curl
WORKDIR /app
COPY --from=builder /build/main /app/main
CMD ["/app/main"]