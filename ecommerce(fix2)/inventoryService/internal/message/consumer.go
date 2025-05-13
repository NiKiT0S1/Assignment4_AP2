package message

import (
	"database/sql"
	"fmt"
	"inventoryService/internal/domain"
	"log"
)

// 6) Получение и обработка сообщения в Inventory Service

// MessageConsumer обрабатывает сообщения из очереди
type MessageConsumer struct {
	productUsecase domain.ProductUsecase
	rabbitClient   *RabbitMQClient
}

// NewMessageConsumer создает новый обработчик сообщений
func NewMessageConsumer(uc domain.ProductUsecase, rabbitClient *RabbitMQClient) *MessageConsumer {
	return &MessageConsumer{
		productUsecase: uc,
		rabbitClient:   rabbitClient,
	}
}

// Start запускает потребителя сообщений
func (c *MessageConsumer) Start() error {
	return c.rabbitClient.ConsumeOrderCreated(c.handleOrderCreated)
}

// handleOrderCreated обрабатывает событие создания заказа
func (c *MessageConsumer) handleOrderCreated(payload MessagePayload) error {
	log.Printf("[Inventory Consumer] Processing order %d with %d items", payload.OrderID, len(payload.Items))

	// Validate order data
	if len(payload.Items) == 0 {
		log.Printf("[Inventory Consumer] Order %d has no items, skipping", payload.OrderID)
		return nil
	}

	// Обработка каждого товара в заказе
	for _, item := range payload.Items {
		// Validate product ID
		if item.ProductID <= 0 {
			log.Printf("[Inventory Consumer] Invalid product ID %d in order %d", item.ProductID, payload.OrderID)
			return fmt.Errorf("invalid product ID %d", item.ProductID)
		}

		// Validate quantity
		if item.Quantity <= 0 {
			log.Printf("[Inventory Consumer] Invalid quantity %d for product %d in order %d",
				item.Quantity, item.ProductID, payload.OrderID)
			return fmt.Errorf("invalid quantity %d for product %d", item.Quantity, item.ProductID)
		}

		// Получаем текущий товар
		product, err := c.productUsecase.GetByID(item.ProductID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("[Inventory Consumer] Product %d not found", item.ProductID)
				return fmt.Errorf("product %d not found", item.ProductID)
			}
			log.Printf("[Inventory Consumer] Error getting product %d: %v", item.ProductID, err)
			return fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}

		// Проверяем достаточность запаса
		if product.Stock < item.Quantity {
			log.Printf("[Inventory Consumer] Insufficient stock for product %d: requested %d, available %d",
				item.ProductID, item.Quantity, product.Stock)
			return fmt.Errorf("insufficient stock for product %d: requested %d, available %d",
				item.ProductID, item.Quantity, product.Stock)
		}

		// Обновляем запас
		product.Stock -= item.Quantity
		if err := c.productUsecase.Update(product); err != nil {
			log.Printf("[Inventory Consumer] Error updating product %d stock: %v", item.ProductID, err)
			return fmt.Errorf("failed to update product %d stock: %w", item.ProductID, err)
		}

		log.Printf("[Inventory Consumer] Updated stock for product %d: new stock %d",
			item.ProductID, product.Stock)
	}

	log.Printf("[Inventory Consumer] Successfully processed order %d", payload.OrderID)
	return nil
}
