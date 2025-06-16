FROM golang:1.24-alpine AS builder

# Устанавливаем git и зависимости
RUN apk add --no-cache git

WORKDIR /app

# Копируем только файлы модулей сначала для кэширования
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем остальные файлы
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o orders .

# Финальный образ
FROM alpine:latest
WORKDIR /app

# Копируем бинарный файл
COPY --from=builder /app/orders .

CMD ["./orders"]