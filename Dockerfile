FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o go-backend-template ./cmd/server/main.go
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o go-backend-template-seed ./cmd/seeder/main.go

FROM migrate/migrate AS migrator

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache tzdata \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /app/logs /app/storage/models \
    && chown -R app:app /app

ENV TZ=Asia/Jakarta

COPY --from=builder /app/go-backend-template .
COPY --from=builder /app/go-backend-template-seed .
COPY --from=migrator /usr/local/bin/migrate /usr/local/bin/migrate

COPY ./scripts/docker-entrypoint.sh /app/entrypoint.sh

USER app

EXPOSE 8000

ENTRYPOINT ["/bin/sh", "/app/entrypoint.sh"]
CMD ["server"]
