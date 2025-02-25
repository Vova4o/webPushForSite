package webpushforsite

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
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
	privateKey *ecdh.PrivateKey
	publicKey  []byte
	subject    string // обычно это URL вашего сайта
}

// NewClient создает новый клиент для отправки уведомлений
func NewClient(subject string) (*Client, error) {
	curve := ecdh.P256()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey().Bytes()

	return &Client{
		privateKey: privateKey,
		publicKey:  publicKey,
		subject:    subject,
	}, nil
}

// SendNotification отправляет push-уведомление подписчику
func (c *Client) SendNotification(subscription *Subscription, message *Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", subscription.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	// Добавляем необходимые заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TTL", "180") // время жизни уведомления в секундах

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetPublicKey возвращает публичный ключ в формате base64url
func (c *Client) GetPublicKey() string {
	return base64.URLEncoding.EncodeToString(c.publicKey)
}
