package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func createProductDB(db *gorm.DB, name string, price float64, quantity int, categoryID uint) (*Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("ราคาต้องไม่ติดลบ (price cannot be negative)")
	}
	if quantity < 0 {
		return nil, errors.New("สต็อกต้องไม่ติดลบ (quantity cannot be negative)")
	}

	var category Category
	if err := db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category ID %d does not exist", categoryID)
		}
		return nil, fmt.Errorf("failed to check category: %w", err)
	}

	product := Product{
		Name:       name,
		Price:      price,
		Quantity:   quantity,
		CategoryID: categoryID,
	}

	if err := db.Create(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &product, nil
}


func getAllProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch all products: %w", err)
	}
	return products, nil
}


func getProductByIDDB(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}
	return &product, nil
}


func getProductByNameDB(db *gorm.DB, name string) ([]Product, error) {
	var products []Product
	searchTerm := "%" + strings.TrimSpace(name) + "%"
	if err := db.Preload("Category").Where("name LIKE ?", searchTerm).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find products by name pattern: %w", err)
	}
	return products, nil
}


func getProductsByCategoryDB(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find products by category: %w", err)
	}
	return products, nil
}


func getProductsByPriceRangeDB(db *gorm.DB, minPrice float64, maxPrice float64) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find products by price range: %w", err)
	}
	return products, nil
}


func getProductsLowStockDB(db *gorm.DB, threshold int) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("quantity < ?", threshold).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find low stock products: %w", err)
	}
	return products, nil
}


func getProductsByCategoryAndPriceDB(db *gorm.DB, categoryID uint, minPrice float64, maxPrice float64) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("category_id = ? AND price BETWEEN ? AND ?", categoryID, minPrice, maxPrice).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find products by category and price range: %w", err)
	}
	return products, nil
}


func getProductsSortedByPriceDB(db *gorm.DB, order string) ([]Product, error) {
	orderUpper := strings.ToUpper(strings.TrimSpace(order))
	if orderUpper != "ASC" && orderUpper != "DESC" {
		orderUpper = "ASC"
	}
	var products []Product
	if err := db.Preload("Category").Order("price " + orderUpper).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find sorted products: %w", err)
	}
	return products, nil
}


func getTopExpensiveProductsDB(db *gorm.DB, limit int) ([]Product, error) {
	if limit <= 0 {
		limit = 5
	}
	var products []Product
	if err := db.Preload("Category").Order("price DESC").Limit(limit).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find top expensive products: %w", err)
	}
	return products, nil
}


func updateProductByIDDB(db *gorm.DB, id uint, name string, price float64, quantity int, categoryID uint) error {
	if price < 0 {
		return errors.New("ราคาต้องไม่ติดลบ (price cannot be negative)")
	}
	if quantity < 0 {
		return errors.New("สต็อกต้องไม่ติดลบ (quantity cannot be negative)")
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	if categoryID != 0 {
		var cat Category
		if err := db.First(&cat, categoryID).Error; err != nil {
			return fmt.Errorf("category ID %d does not exist: %w", categoryID, err)
		}
		product.CategoryID = categoryID
	}

	if strings.TrimSpace(name) != "" {
		product.Name = name
	}
	product.Price = price
	product.Quantity = quantity

	if err := db.Save(&product).Error; err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func updateProductPriceDB(db *gorm.DB, id uint, newPrice float64) error {
	if newPrice < 0 {
		return errors.New("ราคาต้องไม่ติดลบ (price cannot be negative)")
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	if err := db.Model(&product).Update("price", newPrice).Error; err != nil {
		return fmt.Errorf("failed to update product price: %w", err)
	}
	return nil
}

func updateProductStockDB(db *gorm.DB, id uint, changeQty int) error {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	newStock := product.Quantity + changeQty
	if newStock < 0 {
		return fmt.Errorf("สต็อกต้องไม่ติดลบ: สต็อกปัจจุบัน %d ไม่เพียงพอต่อการลด %d (resulting stock %d < 0)", product.Quantity, -changeQty, newStock)
	}

	if err := db.Model(&product).Update("quantity", newStock).Error; err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}
	return nil
}

func softDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	if err := db.Delete(&product).Error; err != nil {
		return fmt.Errorf("failed to soft delete product: %w", err)
	}
	return nil
}

func hardDeleteProductDB(db *gorm.DB, id uint) error {
	var orderCount int64
	if err := db.Model(&Order{}).Where("product_id = ?", id).Count(&orderCount).Error; err != nil {
		return fmt.Errorf("failed to check order history: %w", err)
	}

	if orderCount > 0 {
		return fmt.Errorf("cannot hard delete product ID %d: product is referenced in %d order(s) history (only soft delete is permitted)", id, orderCount)
	}

	result := db.Unscoped().Delete(&Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to hard delete product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product ID %d not found for hard delete", id)
	}
	return nil
}

func getDeletedProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Unscoped().Where("deleted_at IS NOT NULL").Preload("Category").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch soft-deleted products: %w", err)
	}
	return products, nil
}

func restoreProductDB(db *gorm.DB, id uint) error {
	result := db.Unscoped().Model(&Product{}).Where("id = ? AND deleted_at IS NOT NULL", id).Update("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("failed to restore product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no soft-deleted product found with ID %d to restore", id)
	}
	return nil
}
