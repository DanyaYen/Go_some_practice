package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Product struct {
	ID    uint `gorm:"primaryKey"`
	Code  string
	Price uint
}

func dbConnect() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	fmt.Println("Database connection established.")
	return db
}

func main() {
	db := dbConnect()

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting underlying sql.DB: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}()

	fmt.Println("\n--- Running AutoMigrate for Product ---")
	err := db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	fmt.Println("AutoMigrate finished successfully.")

	fmt.Println("\n--- Creating Records ---")

	product1 := Product{Code: "L1212", Price: 100}
	fmt.Printf("Creating product: %+v\n", product1)
	result := db.Create(&product1)
	if result.Error != nil {
		log.Fatalf("Failed to create product1: %v", result.Error)
	}
	fmt.Printf("Product created with ID: %d\n", product1.ID)

	product2 := Product{Code: "P1212", Price: 200}
	fmt.Printf("Creating product: %+v\n", product2)
	result = db.Create(&product2)
	if result.Error != nil {
		log.Fatalf("Failed to create product2: %v", result.Error)
	}
	fmt.Printf("Product created with ID: %d\n", product2.ID)

	product3 := Product{Code: "A1212", Price: 300}
	fmt.Printf("Creating product: %+v\n", product3)
	result = db.Create(&product3)
	if result.Error != nil {
		log.Fatalf("Failed to create product3: %v", result.Error)
	}
	fmt.Printf("Product created with ID: %d\n", product3.ID)

	fmt.Println("All records successfully created.")

	fmt.Println("\n--- Reading All Records ---")
	var products []Product

	result = db.Find(&products)
	if result.Error != nil {
		log.Fatalf("Failed to read all products: %v", result.Error)
	}

	fmt.Printf("Found %d products:\n", len(products))
	if len(products) > 0 {
		for _, product := range products {
			fmt.Printf("  %+v\n", product)
		}
	} else {
		fmt.Println("  No products found.")
	}

	fmt.Println("\n--- Reading Record by ID (e.g., ID 2) ---")
	var productByID Product
	targetID := 2

	result = db.First(&productByID, targetID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("Product with ID %d not found.\n", targetID)
		} else {
			log.Fatalf("Failed to read product with ID %d: %v", targetID, result.Error)
		}
	} else {
		fmt.Printf("Found product with ID %d: %+v\n", targetID, productByID)
	}

	fmt.Println("\n--- Updating Record (e.g., Product with ID 1) ---")
	var productToUpdate Product
	updateID := 1

	result = db.First(&productToUpdate, updateID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("Product with ID %d not found for update.\n", updateID)
		} else {
			log.Fatalf("Failed to find product with ID %d for update: %v", updateID, result.Error)
		}
		return
	}

	newPrice := uint(150)
	productToUpdate.Price = newPrice
	productToUpdate.Code = "UPDATED_L1212"

	fmt.Printf("Updating product with ID %d to: %+v\n", updateID, productToUpdate)
	result = db.Save(&productToUpdate)
	if result.Error != nil {
		log.Fatalf("Failed to update product with ID %d: %v", updateID, result.Error)
	}
	fmt.Printf("Product with ID %d successfully updated.\n", updateID)

	fmt.Println("--- Reading Updated Record (ID 1) ---")
	var updatedProduct Product
	result = db.First(&updatedProduct, updateID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("Product with ID %d not found after update.\n", updateID)
		} else {
			log.Fatalf("Failed to read updated product with ID %d: %v", updateID, result.Error)
		}
	} else {
		fmt.Printf("Updated product: %+v\n", updatedProduct)
	}

	fmt.Println("\n--- Deleting Record (e.g., Product with ID 3) ---")
	var productToDelete Product
	deleteID := 3

	result = db.First(&productToDelete, deleteID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("Product with ID %d not found for deletion.\n", deleteID)
		} else {
			log.Fatalf("Failed to find product with ID %d for deletion: %v", deleteID, result.Error)
		}
		return
	}

	fmt.Printf("Deleting product with ID %d...\n", deleteID)
	result = db.Delete(&productToDelete)
	if result.Error != nil {
		log.Fatalf("Failed to delete product with ID %d: %v", deleteID, result.Error)
	}

	fmt.Printf("Successfully deleted %d record(s) with ID %d.\n", result.RowsAffected, deleteID)

	fmt.Println("--- Attempting to Read Deleted Record (ID 3) ---")
	var deletedCheck Product
	result = db.First(&deletedCheck, deleteID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("As expected, product with ID %d is not found after deletion.\n", deleteID)
		} else {
			log.Printf("Warning: Unexpected error when checking for deleted product with ID %d: %v\n", deleteID, result.Error)
		}
	} else {
		fmt.Printf("Warning: Product with ID %d was found after deletion: %+v\n", deleteID, deletedCheck)
	}

	fmt.Println("\n--- Reading All Records After Deletion ---")
	var remainingProducts []Product
	result = db.Find(&remainingProducts)
	if result.Error != nil {
		log.Fatalf("Failed to read remaining products: %v", result.Error)
	}

	fmt.Printf("Found %d remaining products:\n", len(remainingProducts))
	if len(remainingProducts) > 0 {
		for _, product := range remainingProducts {
			fmt.Printf("  %+v\n", product)
		}
	} else {
		fmt.Println("  No products found.")
	}

	fmt.Println("\nProgram finished.")
}
