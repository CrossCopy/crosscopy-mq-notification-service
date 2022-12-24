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

REDIS_HOST=
REDIS_PORT=
REDIS_PASS=

KAFKA_BOOTSTRAP_SERVERS=
KAFKA_SECURITY_PROTOCOL=
KAFKA_SASL_MECHANISMS=
KAFKA_SASL_USERNAME=
KAFKA_SASL_PASSWORD=
```

## Docker

```bash
docker build . -t ghcr.io/crosscopy/crosscopy-mq-notification-service:main
docker run -it --rm --env-file ./.env ghcr.io/crosscopy/crosscopy-mq-notification-service:main

# to deploy
docker run --restart=unless-stopped -d --env-file ./.env ghcr.io/crosscopy/crosscopy-mq-notification-service:main
```

## Development Environment

MacOS cannot be easily used due to dependency issue with Confluent's kafka library.