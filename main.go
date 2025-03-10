package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Vova4o/webpushnotification/webpushforsite"
)

// Структура для хранения подписчиков в памяти
type SubscriptionStore struct {
	sync.RWMutex
	subscribers []webpushforsite.Subscription
}

func main() {
	// Инициализируем хранилище подписчиков
	store := &SubscriptionStore{
		subscribers: make([]webpushforsite.Subscription, 0),
	}

	client, err := webpushforsite.NewClient("https://example.com")
	if err != nil {
		log.Fatal("Ошибка создания клиента:", err)
	}

	// API endpoints
	http.HandleFunc("/api/push-key", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"publicKey": client.GetPublicKey(),
		})
	})

	http.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		var subscription webpushforsite.Subscription
		if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Сохраняем подписку в памяти
		store.Lock()
		store.subscribers = append(store.subscribers, subscription)
		store.Unlock()

		log.Printf("Новая подписка получена. Всего подписчиков: %d\n", len(store.subscribers))
		w.WriteHeader(http.StatusOK)
	})

	// Добавим эндпоинт для отправки уведомлений всем подписчикам
	http.HandleFunc("/api/send-notification", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var message webpushforsite.Message
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			log.Printf("Ошибка декодирования сообщения: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Получено сообщение для отправки: %+v", message)

		store.RLock()
		failedCount := 0
		log.Printf("Всего подписчиков: %d", len(store.subscribers))

		for i, sub := range store.subscribers {
			log.Printf("Отправка уведомления подписчику %d: %+v", i, sub)

			// Проверка валидности подписки
			if sub.Endpoint == "" {
				log.Printf("Ошибка: пустой endpoint для подписчика %d", i)
				failedCount++
				continue
			}

			// Логирование данных для отправки
			jsonData, _ := json.Marshal(message)
			log.Printf("Отправляемые данные: %s", string(jsonData))

			// Отправка с детальным логированием
			if err := client.SendNotification(&sub, &message); err != nil {
				log.Printf("Ошибка отправки уведомления: %v", err)
				failedCount++
			} else {
				log.Printf("Уведомление успешно отправлено подписчику %d", i)
			}
		}
		store.RUnlock()

		result := map[string]interface{}{
			"total":     len(store.subscribers),
			"failed":    failedCount,
			"succeeded": len(store.subscribers) - failedCount,
		}
		log.Printf("Результат отправки: %+v", result)
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/test-notification", func(w http.ResponseWriter, r *http.Request) {
		message := webpushforsite.Message{
			Title: "Тестовое уведомление",
			Body:  "Это тестовое push-уведомление " + time.Now().Format("15:04:05"),
			URL:   "http://localhost:8080",
		}

		store.RLock()
		defer store.RUnlock()

		if len(store.subscribers) == 0 {
			log.Println("Нет подписчиков!")
			http.Error(w, "Нет активных подписчиков", http.StatusBadRequest)
			return
		}

		// Пробуем отправить первому подписчику
		sub := store.subscribers[0]
		log.Printf("Отправка тестового уведомления подписчику: %+v", sub)

		err := client.SendNotification(&sub, &message)
		if err != nil {
			log.Printf("Ошибка отправки: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Тестовое уведомление отправлено успешно")
		w.Write([]byte("OK"))
	})

	// Раздача статических файлов
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
