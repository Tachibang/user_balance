# Step 1: Модульное кэширование
FROM golang:1.22-alpine AS modules
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Step 2: Сборка приложения
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY --from=modules /go/pkg /go/pkg
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./main.go

# Step 3: Финальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/app /app/app
COPY --from=builder /app/migration /app/migration
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY .env /app/.env

EXPOSE 8080
CMD ["/app/app"]
