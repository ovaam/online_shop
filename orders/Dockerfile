FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orders .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/orders .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080
CMD ["./orders"]