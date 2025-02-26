FROM golang:1.21-alpine

WORKDIR /app

# Копируем только необходимые файлы
COPY go.mod go.sum ./
COPY main.go ./
COPY static/ ./static/

# Устанавливаем git и зависимости
RUN apk add --no-cache git && \
    go mod download

# Собираем приложение
RUN go build -o webpush-app

# Создаем пользователя и настраиваем права
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

USER appuser

# Настраиваем сеть
EXPOSE 8080
ENV HOST=0.0.0.0

CMD ["./webpush-app"]