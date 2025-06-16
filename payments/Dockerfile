FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app
COPY . .

# Установка зависимостей
RUN go mod download

# Сборка
RUN CGO_ENABLED=0 GOOS=linux go build -o payments .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/payments .
CMD ["./payments"]