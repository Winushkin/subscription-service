# 1. Этап сборки
FROM golang:1.26-alpine AS builder

# Устанавливаем swag (если он не установлен в базовом образе)
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости с кэшированием
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# --- ИСПРАВЛЕНИЕ: Сначала копируем весь код ---
COPY . .

# --- ИСПРАВЛЕНИЕ: Теперь генерируем документацию ---
# Убираем `|| true`, чтобы в случае ошибки генерации сборка упала и указала на проблему
RUN swag init -g cmd/main.go

# Собираем бинарник
RUN CGO_ENABLED=0 go build -o myapp ./cmd

# 2. Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/myapp .
COPY --from=builder /app/docs ./docs

# Документируем порт, который использует ваше приложение (например, 8080)
EXPOSE 8080

CMD ["./myapp"]
