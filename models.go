package main

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Category struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

type Product struct {
	gorm.Model
	Name       string   `gorm:"type:varchar(150);not null" json:"name"`
	Price      float64  `gorm:"type:decimal(12,2);not null" json:"price"`
	Quantity   int      `gorm:"not null;default:0" json:"quantity"`
	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Orders     []Order  `gorm:"foreignKey:ProductID" json:"orders,omitempty"`
}

type Customer struct {
	gorm.Model
	Name   string  `gorm:"type:varchar(100);not null" json:"name"`
	Email  string  `gorm:"type:varchar(150);unique;not null" json:"email"`
	Phone  string  `gorm:"type:varchar(30)" json:"phone"`
	Orders []Order `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
}

type Order struct {
	gorm.Model
	CustomerID uint     `gorm:"not null" json:"customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ProductID  uint     `gorm:"not null" json:"product_id"`
	Product    Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	TotalPrice float64  `gorm:"type:decimal(12,2);not null" json:"total_price"`
	Status     string   `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"` // Pending, Completed, Cancelled
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func initDB() (*gorm.DB, error) {
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "store_db")

	var auth string
	if password != "" {
		auth = fmt.Sprintf("%s:%s", user, password)
	} else {
		auth = user
	}

	dsnWithoutDB := fmt.Sprintf("%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", auth, host, port)
	serverDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
	}

	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName)
	if err := serverDB.Exec(createDBSQL).Error; err != nil {
		return nil, fmt.Errorf("failed to create database %s: %w", dbName, err)
	}

	sqlDB, err := serverDB.DB()
	if err == nil {
		_ = sqlDB.Close()
	}

	dsnWithDB := fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", auth, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsnWithDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", dbName, err)
	}

	sqlDB, err = db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	err = db.AutoMigrate(&Category{}, &Product{}, &Customer{}, &Order{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database tables: %w", err)
	}

	return db, nil
}
