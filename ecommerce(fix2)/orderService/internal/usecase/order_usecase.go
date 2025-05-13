package usecase

import (
	"log"
	"orderService/internal/domain"
	"orderService/internal/message"
)

type orderUsecase struct {
	repo     domain.OrderRepository
	producer *message.MessageProducer
}

// 4) Сохранение заказа в БД и отправка события
func NewOrderUsecase(r domain.OrderRepository, p *message.MessageProducer) domain.OrderUsecase {
	return &orderUsecase{r, p}
}

func (uc *orderUsecase) Create(o *domain.Order) error {
	o.Status = "pending"

	// 4.1) Сохраняем заказ в БД
	if err := uc.repo.Create(o); err != nil {
		return err
	}

	// 4.2) Публикуем событие создания заказа
	log.Printf("[OrderUsecase] Publishing order created event for OrderID: %d", o.ID)
	if err := uc.producer.PublishOrderCreated(o); err != nil {
		log.Printf("[OrderUsecase] Failed to publish order created event: %v", err)
		// Продолжаем выполнение даже при ошибке публикации
	}

	return nil
}

func (uc *orderUsecase) GetByID(id int) (*domain.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *orderUsecase) UpdateStatus(id int, status string) error {
	return uc.repo.UpdateStatus(id, status)
}

func (uc *orderUsecase) ListByUser(userID int) ([]domain.Order, error) {
	return uc.repo.ListByUser(userID)
}
