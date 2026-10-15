package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func createCategoryDB(db *gorm.DB, name string, description string) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category := Category{
		Name:        name,
		Description: description,
	}

	if err := db.Create(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

func deleteCategoryByIDDB(db *gorm.DB, id uint) error {
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("category ID %d not found", id)
		}
		return fmt.Errorf("failed to find category: %w", err)
	}

	var productCount int64
	if err := db.Model(&Product{}).Where("category_id = ?", id).Count(&productCount).Error; err != nil {
		return fmt.Errorf("failed to check associated products: %w", err)
	}

	if productCount > 0 {
		return fmt.Errorf("cannot delete category '%s' (ID: %d): there are %d product(s) associated with this category", category.Name, id, productCount)
	}

	if err := db.Delete(&category).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func getAllCategoriesDB(db *gorm.DB) ([]Category, error) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return categories, nil
}

func getCategoryByIDDB(db *gorm.DB, id uint) (*Category, error) {
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
