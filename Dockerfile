FROM golang:1.25-alpine as builder

WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build

FROM restic/restic:0.18.1
RUN apk add --no-cache rclone bash curl docker-cli python3 py3-pip dcron
COPY --from=builder /app/autorestic /usr/bin/autorestic

ENV CRON_SCHEDULE="10 2 * * *"
ENV AUTORESTIC_CONFIG="/data/.autorestic.yml"

RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo 'echo "SHELL=/bin/bash" > /etc/crontabs/root' >> /entrypoint.sh && \
    echo 'echo "PATH=/usr/local/bin:/usr/bin:/bin" >> /etc/crontabs/root' >> /entrypoint.sh && \
    echo 'echo "$CRON_SCHEDULE autorestic backup -c $AUTORESTIC_CONFIG -a >> /var/log/autorestic/backup.log 2>&1" >> /etc/crontabs/root' >> /entrypoint.sh && \
    echo 'echo "Starting cron with schedule: $CRON_SCHEDULE"' >> /entrypoint.sh && \
    echo 'crond -f -l 2' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

ENTRYPOINT []
CMD ["/entrypoint.sh"]
