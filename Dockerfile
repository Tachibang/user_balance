# Step 1: Модульное кэширование
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache cmake
RUN go mod download && \
    make build

# Step 2: Финальный образ
FROM alpine:latest
WORKDIR /app

# differnce between COPY and ADD
COPY --from=builder /bin/app /app/app
COPY --from=builder /app/migration /app/migration
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY .env /app/.env

EXPOSE 8080
CMD ["/app/app"]
