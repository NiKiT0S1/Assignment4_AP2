package message

import (
	"orderService/internal/domain"
	"time"
)

// 5) Публикация сообщения в RabbitMQ

// MessageProducer отправляет сообщения в очередь
type MessageProducer struct {
	rabbitClient *RabbitMQClient
}

// NewMessageProducer создаёт нового продюсера сообщений
func NewMessageProducer(rabbitClient *RabbitMQClient) *MessageProducer {
	return &MessageProducer{
		rabbitClient: rabbitClient,
	}
}

// PublishOrderCreated публикует событие о создании заказа
func (p *MessageProducer) PublishOrderCreated(order *domain.Order) error {
	// Подготовка данных товаров
	var items []Item
	for _, orderItem := range order.Items {
		items = append(items, Item{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		})
	}

	// Создание полезной нагрузки сообщения
	payload := MessagePayload{
		OrderID:   order.ID,
		UserID:    order.UserID,
		Items:     items,
		Status:    order.Status,
		Timestamp: time.Now(),
	}

	// Публикация сообщения в клиенте RabbitMQ
	return p.rabbitClient.PublishOrderCreated(payload)
}
