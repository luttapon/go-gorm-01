package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("เชื่อมต่อฐานข้อมูลล้มเหลว: %v", err)
	}
	fmt.Println("เชื่อมต่อฐานข้อมูล store_db สำเร็จ")

	seedData(db)

	cat, _ := createCategoryDB(db, "Gaming Gear", "อุปกรณ์เล่นเกม")
	prod, _ := createProductDB(db, "Razer Mouse", 2990.00, 15, cat.ID)
	cust, _ := createCustomerDB(db, "Gamer One", "gamerone@email.com", "089-111-2222")

	order, err := createOrderDB(db, cust.ID, prod.ID, 2)
	if err == nil {
		fmt.Printf("สร้าง Order สำเร็จ: ID=%d, ยอดรวม=%.2f, สถานะ=%s\n", order.ID, order.TotalPrice, order.Status)
	}

	products, _ := getAllProductsDB(db)
	fmt.Printf("จำนวนสินค้าทั้งหมด: %d รายการ\n", len(products))

	priceRange, _ := getProductsByPriceRangeDB(db, 1000, 10000)
	fmt.Printf("สินค้าช่วงราคา 1,000 - 10,000: %d รายการ\n", len(priceRange))

	lowStock, _ := getProductsLowStockDB(db, 5)
	fmt.Printf("สินค้าสต็อกน้อยกว่า 5 ชิ้น: %d รายการ\n", len(lowStock))

	topExpensive, _ := getTopExpensiveProductsDB(db, 3)
	if len(topExpensive) >= 3 {
		fmt.Printf("สินค้าแพงที่สุด 3 อันดับ: 1.%s 2.%s 3.%s\n", topExpensive[0].Name, topExpensive[1].Name, topExpensive[2].Name)
	}

	orders, _ := getOrdersWithDetailsDB(db)
	fmt.Printf("คำสั่งซื้อทั้งหมดในระบบ: %d รายการ\n", len(orders))

	_ = updateProductPriceDB(db, prod.ID, 2790.00)
	_ = updateProductStockDB(db, prod.ID, 5)
	_ = updateOrderStatusDB(db, order.ID, "Completed")
	fmt.Println("อัปเดตราคา, สต็อก และสถานะ Order สำเร็จ")

	_ = softDeleteProductDB(db, prod.ID)
	deleted, _ := getDeletedProductsDB(db)
	fmt.Printf("สินค้าที่ถูก Soft Delete: %d รายการ\n", len(deleted))

	_ = restoreProductDB(db, prod.ID)
	fmt.Println("กู้คืนสินค้า (Restore) สำเร็จ")

	if err := deleteCategoryByIDDB(db, cat.ID); err != nil {
		fmt.Printf("ตรวจสอบความสัมพันธ์: %v\n", err)
	}
	if err := deleteCustomerByIDDB(db, cust.ID); err != nil {
		fmt.Printf("ตรวจสอบความสัมพันธ์: %v\n", err)
	}

	fmt.Println("การทดสอบการทำงานของระบบเสร็จสิ้นเรียบร้อย")
}

func seedData(db *gorm.DB) {
	var count int64
	db.Model(&Product{}).Count(&count)
	if count >= 15 {
		return
	}

	cats := []Category{
		{Name: "Laptop", Description: "แล็ปท็อปและคอมพิวเตอร์พกพา"},
		{Name: "Smartphone", Description: "สมาร์ตโฟนและแท็บเล็ต"},
		{Name: "Accessories", Description: "อุปกรณ์ต่อพ่วง"},
		{Name: "Hardware", Description: "ชิ้นส่วนประกอบคอมพิวเตอร์"},
	}
	for i := range cats {
		db.FirstOrCreate(&cats[i], Category{Name: cats[i].Name})
	}

	prods := []Product{
		{Name: "MacBook Air M3", Price: 39900, Quantity: 10, CategoryID: cats[0].ID},
		{Name: "ASUS Vivobook 15", Price: 22900, Quantity: 8, CategoryID: cats[0].ID},
		{Name: "Dell XPS 13", Price: 45900, Quantity: 4, CategoryID: cats[0].ID},
		{Name: "Lenovo Legion 5", Price: 35900, Quantity: 6, CategoryID: cats[0].ID},
		{Name: "iPhone 16", Price: 32900, Quantity: 15, CategoryID: cats[1].ID},
		{Name: "iPhone 16 Pro", Price: 41900, Quantity: 12, CategoryID: cats[1].ID},
		{Name: "Samsung Galaxy S24", Price: 43900, Quantity: 7, CategoryID: cats[1].ID},
		{Name: "iPad Pro", Price: 31900, Quantity: 9, CategoryID: cats[1].ID},
		{Name: "Logitech MX Master 3S", Price: 3590, Quantity: 20, CategoryID: cats[2].ID},
		{Name: "Mechanical Keyboard", Price: 3890, Quantity: 15, CategoryID: cats[2].ID},
		{Name: "USB-C Hub", Price: 1290, Quantity: 25, CategoryID: cats[2].ID},
		{Name: "Sony WH-1000XM5", Price: 12900, Quantity: 3, CategoryID: cats[2].ID},
		{Name: "Samsung SSD 1TB", Price: 4290, Quantity: 18, CategoryID: cats[3].ID},
		{Name: "RAM DDR5 32GB", Price: 4890, Quantity: 2, CategoryID: cats[3].ID},
		{Name: "RTX 4070", Price: 24500, Quantity: 5, CategoryID: cats[3].ID},
	}
	for i := range prods {
		db.FirstOrCreate(&prods[i], Product{Name: prods[i].Name})
	}

	custs := []Customer{
		{Name: "Somchai", Email: "somchai@email.com", Phone: "081-111-1111"},
		{Name: "Suda", Email: "suda@email.com", Phone: "082-222-2222"},
		{Name: "Wichai", Email: "wichai@email.com", Phone: "083-333-3333"},
		{Name: "Malee", Email: "malee@email.com", Phone: "084-444-4444"},
		{Name: "Anan", Email: "anan@email.com", Phone: "085-555-5555"},
		{Name: "Kannika", Email: "kannika@email.com", Phone: "086-666-6666"},
		{Name: "Prasert", Email: "prasert@email.com", Phone: "087-777-7777"},
		{Name: "Natthapong", Email: "natthapong@email.com", Phone: "088-888-8888"},
	}
	for i := range custs {
		db.FirstOrCreate(&custs[i], Customer{Email: custs[i].Email})
	}

	for i := 0; i < 12; i++ {
		cID := custs[i%len(custs)].ID
		pID := prods[i%len(prods)].ID
		var existing Order
		if err := db.Where("customer_id = ? AND product_id = ?", cID, pID).First(&existing).Error; err != nil {
			db.Create(&Order{
				CustomerID: cID,
				ProductID:  pID,
				Quantity:   1,
				TotalPrice: prods[i%len(prods)].Price,
				Status:     "Completed",
			})
		}
	}
}
