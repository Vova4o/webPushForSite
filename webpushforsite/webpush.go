package webpushforsite

import (
	"encoding/json"

	"github.com/SherClockHolmes/webpush-go"
)

// Subscription содержит данные подписки браузера
type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     Keys   `json:"keys"`
}

// Keys содержит ключи для шифрования сообщений
type Keys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

// Message представляет структуру push-уведомления
type Message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon,omitempty"`
	URL   string `json:"url,omitempty"`
}

// Client для отправки web push уведомлений
type Client struct {
	vapidPublicKey  string
	vapidPrivateKey string
	subject         string // обычно это URL вашего сайта
}

// NewClient создает новый клиент для отправки уведомлений
func NewClient(subject string) (*Client, error) {
	// Генерируем VAPID ключи
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return nil, err
	}

	return &Client{
		vapidPublicKey:  publicKey,
		vapidPrivateKey: privateKey,
		subject:         subject,
	}, nil
}

// SendNotification отправляет push-уведомление подписчику
func (c *Client) SendNotification(sub *Subscription, msg *Message) error {
	// Преобразуем нашу подписку в формат библиотеки
	s := webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.Keys.P256dh,
			Auth:   sub.Keys.Auth,
		},
	}

	// Преобразуем сообщение в JSON
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Отправляем уведомление
	_, err = webpush.SendNotification(
		payload,
		&s,
		&webpush.Options{
			VAPIDPublicKey:  c.vapidPublicKey,
			VAPIDPrivateKey: c.vapidPrivateKey,
			TTL:             30,
			Subscriber:      c.subject,
		},
	)

	return err
}

// GetPublicKey возвращает публичный ключ в формате base64url
func (c *Client) GetPublicKey() string {
	return c.vapidPublicKey
}
