# crosscopy-mq-notification-service

A notification service (golang) for sending emails and potentially more updates according to data and events received in message queue.

## Environment Variables

Make a `.env` with the following variables.

```
EMAIL_ADDRESS=noreply@crosscopy.io
EMAIL_PASSWORD=
MAIL_SERVER=
KAFKA_MODE=cloud
KAFKA_GROUP_ID=crosscopy-notification


KAFKA_BOOTSTRAP_SERVERS=
KAFKA_SECURITY_PROTOCOL=
KAFKA_SASL_MECHANISMS=
KAFKA_SASL_USERNAME=
KAFKA_SASL_PASSWORD=
```

## Docker

```bash
docker build . -t ghcr.io/crosscopy/crosscopy-mq-notification-service:latest
docker run -it --rm --env-file ./.env ghcr.io/crosscopy/crosscopy-mq-notification-service:latest
```