package postgres

import (
	"log"
	"order_svc/internal/models"
)

// CreateNewOrder creates a new order
func (postgresDB *PostgresDB) CreateNewOrder(order *models.NewOrder) error {
	postgresDB.Init()
	err := postgresDB.DB.Create(order).Error
	if err != nil {
		log.Println("Error creating new order:", err)
		return err
	}
	return nil
}

func (postgresDB *PostgresDB) FetchOrdersByUser(userID string) ([]models.NewOrder, error) {
	var orders []models.NewOrder
	if err := postgresDB.DB.Where("user_id = ?", userID).Preload("CartItems").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (postgresDB *PostgresDB) CartOrderExist(cartID string) bool {
	order := &models.NewOrder{}

	result := postgresDB.DB.First(order, "cart_id = ?", cartID)

	return result.RowsAffected > 0
}

func (postgresDB *PostgresDB) OrderExist(orderID string) bool {
	order := &models.NewOrder{}

	result := postgresDB.DB.First(order, "id = ?", orderID)

	return result.RowsAffected > 0
}

func (postgresDB *PostgresDB) UpdateOrder(orderID, status string) (*models.NewOrder, error) {
	order := &models.NewOrder{}
	err := postgresDB.DB.Model(order).Where("id = ?", orderID).Update("status", status).Error
	return order, err
}

func (postgresDB *PostgresDB) FetchCartItemByUserID(userID string) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	if err := postgresDB.DB.Where("userId = ?", userID).Preload("CartPackageItem").Find(&cartItems).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}
