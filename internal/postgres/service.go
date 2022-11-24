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

func (postgresDB *PostgresDB) DeleteOrderByUserID(userID string) error {
	if err := postgresDB.DB.Where("user_id = ?", userID).Delete(&models.CartPackageItem{}).Error; err != nil {
		return err
	}
	if err := postgresDB.DB.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		return err
	}
	if err := postgresDB.DB.Where("user_id = ?", userID).Delete(&models.NewOrder{}).Error; err != nil {
		return err
	}
	return nil
}

func (postgresDB *PostgresDB) UpdateOrder(orderID, status string) (*models.NewOrder, error) {
	order := &models.NewOrder{}
	err := postgresDB.DB.Model(order).Where("order_id = ?", orderID).Update("status", status).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
