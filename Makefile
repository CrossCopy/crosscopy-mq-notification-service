deploy:
	docker run --restart=unless-stopped -d --env-file ./.env --name crosscopy-mq-notification-service ghcr.io/crosscopy/crosscopy-mq-notification-service:main

destroy:
	docker stop crosscopy-mq-notification-service && docker rm crosscopy-mq-notification-service