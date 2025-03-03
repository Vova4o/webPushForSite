# Web Push Notification Service

Этот проект предоставляет сервис для отправки web push уведомлений. Он включает серверную часть на Go и клиентскую часть для взаимодействия с браузером.

## Структура проекта

- `main.go`: Основной файл сервера, который обрабатывает подписки и отправку уведомлений.
- `webpushforsite/webpush.go`: Логика для работы с web push уведомлениями.
- `static/`: Статические файлы, включая `service-worker.js` и `test.html`.
- `example/`: Примеры использования.

## Установка

1. Клонируйте репозиторий:

    ```bash
    git clone https://github.com/Vova4o/webPushForSite
    cd webPushForSite
    ```

2. Установите зависимости:

    ```bash
    go mod tidy
    ```

3. Соберите проект:

    ```bash
    go build -o webpush-app
    ```

## Запуск

1. Запустите сервер:

    ```bash
    ./webpush-app
    ```

2. Откройте браузер и перейдите на `http://localhost:8080`.

## Использование

### Подписка на уведомления

1. Откройте `http://localhost:8080/test.html`.
2. Нажмите кнопку "Создать подписку".
3. Разрешите уведомления в браузере.

### Отправка уведомлений

1. Откройте `http://localhost:8080/test.html`.
2. Нажмите кнопку "Тест уведомления через API".

### Прямое тестирование уведомлений

1. Откройте `http://localhost:8080/test.html`.
2. Нажмите кнопку "Прямой тест уведомлений".

## Примеры

### Пример использования в `example/`

В папке `example/` вы найдете пример использования сервиса для отправки уведомлений.

```go
package main

import (
    "log"
    "net/http"
    "github.com/Vova4o/webpushnotification/webpushforsite"
)

func main() {
    client, err := webpushforsite.NewClient("https://example.com")
    if err != nil {
        log.Fatal("Ошибка создания клиента:", err)
    }

    http.HandleFunc("/api/push-key", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{
            "publicKey": client.GetPublicKey(),
        })
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}