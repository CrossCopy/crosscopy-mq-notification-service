deploy:
	docker run --restart=unless-stopped -d --env-file ./.env ghcr.io/crosscopy/crosscopy-mq-notification-service:main