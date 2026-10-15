package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func createCustomerDB(db *gorm.DB, name string, email string, phone string) (*Customer, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)

	if name == "" {
		return nil, errors.New("customer name cannot be empty")
	}
	if email == "" {
		return nil, errors.New("customer email cannot be empty")
	}

	customer := Customer{
		Name:  name,
		Email: email,
		Phone: phone,
	}

	if err := db.Create(&customer).Error; err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &customer, nil
}

func getCustomerByEmailDB(db *gorm.DB, email string) (*Customer, error) {
	var customer Customer
	if err := db.Where("email = ?", strings.TrimSpace(email)).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("customer with email '%s' not found", email)
		}
		return nil, fmt.Errorf("failed to find customer by email: %w", err)
	}
	return &customer, nil
}


func getAllCustomersDB(db *gorm.DB) ([]Customer, error) {
	var customers []Customer
	if err := db.Find(&customers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch customers: %w", err)
	}
	return customers, nil
}

func getCustomerByIDDB(db *gorm.DB, id uint) (*Customer, error) {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func updateCustomerByIDDB(db *gorm.DB, id uint, name string, email string, phone string) error {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return fmt.Errorf("customer ID %d not found: %w", id, err)
	}

	updates := map[string]interface{}{}
	if strings.TrimSpace(name) != "" {
		updates["name"] = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		updates["email"] = strings.TrimSpace(email)
	}
	if strings.TrimSpace(phone) != "" {
		updates["phone"] = strings.TrimSpace(phone)
	}

	if len(updates) > 0 {
		if err := db.Model(&customer).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update customer: %w", err)
		}
	}
	return nil
}

func deleteCustomerByIDDB(db *gorm.DB, id uint) error {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return fmt.Errorf("customer ID %d not found: %w", id, err)
	}

	var orderCount int64
	if err := db.Model(&Order{}).Where("customer_id = ?", id).Count(&orderCount).Error; err != nil {
		return fmt.Errorf("failed to check customer orders: %w", err)
	}

	if orderCount > 0 {
		return fmt.Errorf("cannot delete customer '%s' (ID: %d): customer has %d existing order(s) in history", customer.Name, id, orderCount)
	}

	if err := db.Delete(&customer).Error; err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}
