package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)


func createOrderDB(db *gorm.DB, customerID uint, productID uint, quantity int) (*Order, error) {
	if quantity <= 0 {
		return nil, errors.New("จำนวนที่สั่งต้องมากกว่า 0 (order quantity must be greater than 0)")
	}

	var createdOrder Order

	err := db.Transaction(func(tx *gorm.DB) error {
		var customer Customer
		if err := tx.First(&customer, customerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("customer ID %d not found", customerID)
			}
			return fmt.Errorf("failed to query customer: %w", err)
		}

		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("product ID %d not found", productID)
			}
			return fmt.Errorf("failed to query product: %w", err)
		}

		if product.Quantity < quantity {
			return fmt.Errorf("insufficient stock for '%s': requested %d, but only %d available", product.Name, quantity, product.Quantity)
		}

		totalPrice := product.Price * float64(quantity)

		newStock := product.Quantity - quantity
		if err := tx.Model(&product).Update("quantity", newStock).Error; err != nil {
			return fmt.Errorf("failed to reduce product stock: %w", err)
		}

		createdOrder = Order{
			CustomerID: customerID,
			ProductID:  productID,
			Quantity:   quantity,
			TotalPrice: totalPrice,
			Status:     "Pending",
		}

		if err := tx.Create(&createdOrder).Error; err != nil {
			return fmt.Errorf("failed to insert order: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := db.Preload("Customer").Preload("Product").First(&createdOrder, createdOrder.ID).Error; err != nil {
		return &createdOrder, nil
	}

	return &createdOrder, nil
}

func getOrdersByCustomerDB(db *gorm.DB, customerID uint) ([]Order, error) {
	var orders []Order
	if err := db.Where("customer_id = ?", customerID).Preload("Customer").Preload("Product").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders by customer: %w", err)
	}
	return orders, nil
}

func getOrdersByStatusDB(db *gorm.DB, status string) ([]Order, error) {
	var orders []Order
	if err := db.Where("status = ?", strings.TrimSpace(status)).Preload("Customer").Preload("Product").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders by status: %w", err)
	}
	return orders, nil
}

func getOrdersByDateRangeDB(db *gorm.DB, startDate time.Time, endDate time.Time) ([]Order, error) {
	var orders []Order
	if err := db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Preload("Customer").Preload("Product").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders by date range: %w", err)
	}
	return orders, nil
}

func getOrdersWithDetailsDB(db *gorm.DB) ([]Order, error) {
	var orders []Order
	if err := db.Preload("Customer").Preload("Product").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders with details: %w", err)
	}
	return orders, nil
}


func updateOrderStatusDB(db *gorm.DB, id uint, status string) error {
	status = strings.TrimSpace(status)
	validStatuses := map[string]bool{
		"Pending":   true,
		"Completed": true,
		"Cancelled": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid order status '%s': must be 'Pending', 'Completed', or 'Cancelled'", status)
	}

	var order Order
	if err := db.First(&order, id).Error; err != nil {
		return fmt.Errorf("order ID %d not found: %w", id, err)
	}

	if err := db.Model(&order).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
