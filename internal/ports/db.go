package ports

import "order_svc/internal/models"

type Order interface {
	CreateNewOrder(order *models.NewOrder) error
	FetchOrdersByUser(userID string) ([]models.NewOrder, error)
}
