
-- SQL Export: store_db
-- Assignment: Go + GORM + MySQL
-- Database: store_db

CREATE DATABASE IF NOT EXISTS `store_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `store_db`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `description` text,
  PRIMARY KEY (`id`),
  KEY `idx_categories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `categories` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `description`) VALUES
(1, NOW(), NOW(), NULL, 'Laptop', 'แล็ปท็อปและคอมพิวเตอร์พกพาสำหรับการทำงานและเล่นเกม'),
(2, NOW(), NOW(), NULL, 'Smartphone', 'สมาร์ตโฟนและแท็บเล็ตระดับพรีเมียม'),
(3, NOW(), NOW(), NULL, 'Accessories', 'อุปกรณ์ต่อพ่วง เมาส์ คีย์บอร์ด และฮับ'),
(4, NOW(), NOW(), NULL, 'Hardware', 'ชิ้นส่วนประกอบคอมพิวเตอร์ การ์ดจอ SSD RAM');

DROP TABLE IF EXISTS `products`;
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(150) NOT NULL,
  `price` decimal(12,2) NOT NULL,
  `quantity` bigint NOT NULL DEFAULT '0',
  `category_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`),
  KEY `fk_categories_products` (`category_id`),
  CONSTRAINT `fk_categories_products` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `products` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `price`, `quantity`, `category_id`) VALUES
(1, NOW(), NOW(), NULL, 'MacBook Air M3', 39900.00, 10, 1),
(2, NOW(), NOW(), NULL, 'ASUS Vivobook 15', 22900.00, 8, 1),
(3, NOW(), NOW(), NULL, 'Dell XPS 13', 45900.00, 4, 1),
(4, NOW(), NOW(), NULL, 'Lenovo Legion 5', 35900.00, 6, 1),
(5, NOW(), NOW(), NULL, 'iPhone 16', 32900.00, 15, 2),
(6, NOW(), NOW(), NULL, 'iPhone 16 Pro', 41900.00, 12, 2),
(7, NOW(), NOW(), NULL, 'Samsung Galaxy S24 Ultra', 43900.00, 7, 2),
(8, NOW(), NOW(), NULL, 'iPad Pro 11-inch', 31900.00, 9, 2),
(9, NOW(), NOW(), NULL, 'Logitech Mouse MX Master 3S', 3590.00, 20, 3),
(10, NOW(), NOW(), NULL, 'Mechanical Keyboard Keychron K2', 3890.00, 15, 3),
(11, NOW(), NOW(), NULL, 'USB-C Hub 7-in-1', 1290.00, 25, 3),
(12, NOW(), NOW(), NULL, 'Sony WH-1000XM5 Headphones', 12900.00, 3, 3),
(13, NOW(), NOW(), NULL, 'Samsung 990 Pro SSD 1TB', 4290.00, 18, 4),
(14, NOW(), NOW(), NULL, 'Corsair DDR5 32GB RAM', 4890.00, 2, 4),
(15, NOW(), NOW(), NULL, 'NVIDIA GeForce RTX 4070', 24500.00, 5, 4);

DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(150) NOT NULL,
  `phone` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_customers_email` (`email`),
  KEY `idx_customers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `customers` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `email`, `phone`) VALUES
(1, NOW(), NOW(), NULL, 'Somchai Jaidee', 'somchai@email.com', '081-111-1111'),
(2, NOW(), NOW(), NULL, 'Suda Rakdee', 'suda@email.com', '082-222-2222'),
(3, NOW(), NOW(), NULL, 'Wichai Sukjai', 'wichai@email.com', '083-333-3333'),
(4, NOW(), NOW(), NULL, 'Malee Charoen', 'malee@email.com', '084-444-4444'),
(5, NOW(), NOW(), NULL, 'Anan Panyarat', 'anan@email.com', '085-555-5555'),
(6, NOW(), NOW(), NULL, 'Kannika Srisuk', 'kannika@email.com', '086-666-6666'),
(7, NOW(), NOW(), NULL, 'Prasert Meechok', 'prasert@email.com', '087-777-7777'),
(8, NOW(), NOW(), NULL, 'Natthapong Wongsuwan', 'natthapong@email.com', '088-888-8888');

DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `customer_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `quantity` bigint NOT NULL,
  `total_price` decimal(12,2) NOT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'Pending',
  PRIMARY KEY (`id`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `fk_orders_customer` (`customer_id`),
  KEY `fk_orders_product` (`product_id`),
  CONSTRAINT `fk_orders_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_orders_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


INSERT INTO `orders` (`id`, `created_at`, `updated_at`, `deleted_at`, `customer_id`, `product_id`, `quantity`, `total_price`, `status`) VALUES
(1, NOW(), NOW(), NULL, 1, 1, 1, 39900.00, 'Completed'),
(2, NOW(), NOW(), NULL, 1, 9, 2, 7180.00, 'Completed'),
(3, NOW(), NOW(), NULL, 2, 5, 1, 32900.00, 'Completed'),
(4, NOW(), NOW(), NULL, 2, 11, 3, 3870.00, 'Pending'),
(5, NOW(), NOW(), NULL, 3, 2, 1, 22900.00, 'Pending'),
(6, NOW(), NOW(), NULL, 3, 10, 1, 3890.00, 'Completed'),
(7, NOW(), NOW(), NULL, 4, 7, 1, 43900.00, 'Pending'),
(8, NOW(), NOW(), NULL, 5, 13, 2, 8580.00, 'Completed'),
(9, NOW(), NOW(), NULL, 6, 8, 1, 31900.00, 'Completed'),
(10, NOW(), NOW(), NULL, 6, 14, 1, 4890.00, 'Pending'),
(11, NOW(), NOW(), NULL, 7, 15, 1, 24500.00, 'Cancelled'),
(12, NOW(), NOW(), NULL, 8, 3, 1, 45900.00, 'Completed');

SET FOREIGN_KEY_CHECKS = 1;
