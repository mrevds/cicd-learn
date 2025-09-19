FROM golang:1.20-alpine AS builder

# Устанавливаем git и другие зависимости
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Запускаем тесты
RUN go test -v ./...

# Собираем статически (важно для Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Финальный образ
FROM alpine:latest

# Добавляем сертификаты и timezone данные
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]