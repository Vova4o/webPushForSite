FROM golang:1.21-alpine

WORKDIR /app

# Устанавливаем git для загрузки зависимостей
RUN apk add --no-cache git

# Копируем файлы проекта
COPY . .

# Загружаем зависимости и собираем приложение
RUN go mod download && \
    go build -o main

# Создаем пользователя для запуска приложения
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

USER appuser

# Важно: прослушиваем все интерфейсы (0.0.0.0)
ENV HOST=0.0.0.0
EXPOSE 8080

# Запускаем приложение на 0.0.0.0
CMD ["./main"]