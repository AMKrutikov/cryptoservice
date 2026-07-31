FROM golang:1.26-alpine AS builder
WORKDIR /app

# 1. Копируем файлы зависимостей и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# 2. Копируем весь исходный код проекта (включая вашу готовую папку docs/)
COPY . .

# 3. Сразу компилируем Go-приложение
# CGO_ENABLED=0 (Отключение зависимости от C-библиотек)
# GOOS=linux (Целевая операционная система)
RUN CGO_ENABLED=0 GOOS=linux go build -o cryptoservice ./cmd/app/main.go

# ЭТАП 2: Финальный легковесный образ
FROM alpine:3.20

WORKDIR /app

# 1. Устанавливаем SSL-сертификаты для HTTPS-запросов к CoinGecko
# Они необходимы для работы ActualizeRates (в легковесном alpine:3.20 их нет)
RUN apk --no-cache add ca-certificates

# 2. Копируем скомпилированный бинарник
COPY --from=builder /app/cryptoservice .

# 3. Копируем папку с файлом конфигурации
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./cryptoservice"]
