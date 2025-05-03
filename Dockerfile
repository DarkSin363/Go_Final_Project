FROM golang:1.24.2 AS builder

RUN apt-get update && apt-get install -y gcc musl-dev

WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server .

FROM alpine:latest

RUN apk update && \
    apk add --no-cache sqlite libc6-compat && \
    mkdir -p /app/data && \
    adduser -D appuser && \
    chown appuser:appuser /app/data

USER appuser
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/web ./web/

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db

VOLUME /data

CMD ["./server"]