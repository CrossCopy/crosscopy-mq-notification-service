# syntax=docker/dockerfile:1

FROM golang:1.18.8 AS builder
WORKDIR /go/src/github.com/CrossCopy/crosscopy-mq-notification-service
COPY . ./
RUN go build .

FROM ubuntu:20.04
WORKDIR /root/
RUN apt-get update && apt-get install ca-certificates -y && apt-get clean
COPY --from=builder /go/src/github.com/CrossCopy/crosscopy-mq-notification-service/crosscopy-mq-notification-service ./
CMD ["/root/crosscopy-mq-notification-service"]


# docker build . -t ghcr.io/crosscopy/crosscopy-mq-notification-service:latest
# docker run -it --rm --env-file ./.env ghcr.io/crosscopy/crosscopy-mq-notification-service:latest