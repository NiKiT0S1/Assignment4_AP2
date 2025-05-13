package message

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQClient представляет клиент для работы с RabbitMQ
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// MessagePayload представляет полезную нагрузку сообщения
type MessagePayload struct {
	OrderID   int       `json:"order_id"`
	UserID    int       `json:"user_id"`
	Items     []Item    `json:"items"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// Item представляет элемент заказа
type Item struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// NewRabbitMQClient создаёт новый клиент RabbitMQ
func NewRabbitMQClient(url, queueName string) (*RabbitMQClient, error) {
	// Подключение к RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Создание канала
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Объявление очереди
	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Создание exchange и привязка к очереди
	err = channel.ExchangeDeclare(
		"order_events", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	err = channel.QueueBind(
		queue.Name,      // queue name
		"order.created", // routing key
		"order_events",  // exchange
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

// PublishOrderCreated публикует сообщение о создании заказа
func (c *RabbitMQClient) PublishOrderCreated(payload MessagePayload) error {
	// Логирование начала публикации
	log.Printf("[RabbitMQ Producer] Publishing order created event for OrderID: %d", payload.OrderID)

	// Сериализация полезной нагрузки
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Публикация сообщения
	err = c.channel.Publish(
		"order_events",  // exchange
		"order.created", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	// Логирование успешной публикации
	log.Printf("[RabbitMQ Producer] Successfully published order created event for OrderID: %d", payload.OrderID)
	return nil
}

// ConsumeOrderCreated потребляет сообщения о создании заказа
func (c *RabbitMQClient) ConsumeOrderCreated(handler func(MessagePayload) error) error {
	// Логирование начала потребления
	log.Printf("[RabbitMQ Consumer] Starting to consume order created events")

	// Настройка потребления сообщений
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Обработка сообщений
	go func() {
		for d := range msgs {
			log.Printf("[RabbitMQ Consumer] Received message: %s", d.Body)

			// Десериализация полезной нагрузки
			var payload MessagePayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("[RabbitMQ Consumer] Error parsing message: %v", err)
				d.Nack(false, true) // Ошибка обработки, вернуть в очередь
				continue
			}

			// Обработка сообщения
			log.Printf("[RabbitMQ Consumer] Processing order created event for OrderID: %d", payload.OrderID)
			if err := handler(payload); err != nil {
				log.Printf("[RabbitMQ Consumer] Error handling message: %v", err)
				d.Nack(false, true) // Ошибка обработки, вернуть в очередь
				continue
			}

			// Подтверждение обработки
			d.Ack(false)
			log.Printf("[RabbitMQ Consumer] Successfully processed order created event for OrderID: %d", payload.OrderID)
		}
	}()

	log.Printf("[RabbitMQ Consumer] Waiting for order created events...")
	return nil
}

// Close закрывает соединение с RabbitMQ
func (c *RabbitMQClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
