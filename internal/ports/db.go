package ports

import "order_svc/internal/models"

type Order interface {
	CreateNewOrder(order *models.NewOrder) error
	FetchOrdersByUser(userID string) ([]models.NewOrder, error)
	DeleteOrderByUserID(userID string) error
	UpdateOrder(orderID, status string) (*models.NewOrder, error)
	FetchAllOrders() ([]models.NewOrder, error)
	FetchUserOrderByStatus(userID, status string) ([]models.NewOrder, error)
	GetDeliveryCode(orderId string) (int32, error)
	UpdateDeliveryTime(orderID string, deliveryTime string) (*models.NewOrder, error)
}
