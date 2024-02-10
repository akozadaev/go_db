package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	//gorm.Model
	ID    uint `gorm:"primarykey"`
	Name  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Name: "Мороженка", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                       // find product with integer primary key
	db.First(&product, "name = ?", "Мороженка") // find product with code Мороженка

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 201)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 202, Name: "Мороженка"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 203, "Name": "Мороженка1"})

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&Product{}).Where("id = ?", 1).Limit(10).Order("id desc").Find(&[]Product{})
	})

	fmt.Println(sql)

	var name any
	rows, err := db.Model(&Product{}).Where("name = ?", "Мороженка1").Select("name").Rows()
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&name)

		fmt.Println(name)
	}
	// Delete - delete product
	db.Delete(&product, 1)
	db.Delete(&product, "name = ?", "Мороженка")

}
