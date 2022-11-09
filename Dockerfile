# syntax=docker/dockerfile:1

FROM golang:1.18.8 AS builder
WORKDIR /go/src/github.com/CrossCopy/crosscopy-mq-notification-service
COPY . ./
RUN go build .

FROM ubuntu:20.04
ENV EMAIL_ADDRESS ""
ENV EMAIL_PASSWORD ""
ENV MAIL_SERVER mail.privateemail.com:587
ENV KAFKA_MODE local
ENV KAFKA_GROUP_ID crosscopy-notification
ENV KAFKA_BOOTSTRAP_SERVERS localhost:9093
ENV KAFKA_SECURITY_PROTOCOL ""
ENV KAFKA_SASL_MECHANISMS ""
ENV KAFKA_SASL_USERNAME ""
ENV KAFKA_SASL_PASSWORD ""
WORKDIR /root/
COPY --from=builder /go/src/github.com/CrossCopy/crosscopy-mq-notification-service/crosscopy-mq-notification-service ./
CMD ["/root/crosscopy-mq-notification-service"]


# docker build . -t crosscopy-mq-notification-service:latest
# docker run -it --rm crosscopy-mq-notification-service:latest