package main

import (
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
)

type Item struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Quantity int    `db:"quantity"`
}

func main() {
	dbPath := "./test.db"

	fmt.Printf("Connecting to SQLite database at path: %s\n", dbPath)
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}()

	fmt.Println("Database connection established.")

	fmt.Println("\n--- Creating 'items' table (if not exists) ---")
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		quantity INTEGER NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating 'items' table: %v", err)
	}
	fmt.Println("'items' table created or already exists.")

	fmt.Println("\n--- Inserting 3 records into the table ---")

	insertSQL := "INSERT INTO items (name, quantity) VALUES (?, ?)"

	fmt.Println("Inserting record: 'item1', 10")
	_, err = db.Exec(insertSQL, "item1", 10)
	if err != nil {
		log.Fatalf("Error inserting record 'item1': %v", err)
	}

	fmt.Println("Inserting record: 'item2', 20")
	_, err = db.Exec(insertSQL, "item2", 20)
	if err != nil {
		log.Fatalf("Error inserting record 'item2': %v", err)
	}

	fmt.Println("Inserting record: 'item3', 30")
	_, err = db.Exec(insertSQL, "item3", 30)
	if err != nil {
		log.Fatalf("Error inserting record 'item3': %v", err)
	}
	fmt.Println("All records successfully inserted.")

	fmt.Println("\n--- Reading all records from the table ---")
	selectALLSQL := "SELECT id, name, quantity FROM items"

	var items []Item

	err = db.Select(&items, selectALLSQL)
	if err != nil {
		log.Fatalf("Error reading all records: %v", err)
	}

	if len(items) == 0 {
		fmt.Println("Table is empty.")
	} else {
		fmt.Println("Found records:")
		for _, item := range items {
			fmt.Printf("ID: %d, Name: %s, Quantity: %d\n", item.ID, item.Name, item.Quantity)
		}
	}


	fmt.Println("\n--- Updating record with ID = 1 ---")
	updateSQL := "UPDATE items SET quantity = ? WHERE id = ?"
	itemIDToUpdate := 1
	newQuantity := 15

	result, err := db.Exec(updateSQL, newQuantity, itemIDToUpdate)
	if err != nil {
		log.Fatalf("Error updating record with ID=%d: %v", itemIDToUpdate, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Warning: Could not get rows affected after UPDATE: %v", err)
	} else {
		fmt.Printf("Record with ID=%d successfully updated. Rows affected: %d\n", itemIDToUpdate, rowsAffected)
	}

	fmt.Println("--- Reading updated record with ID = 1 ---")
	selectOneSQL := "SELECT id, name, quantity FROM items WHERE id = ?"
	var updatedItem Item

	err = db.Get(&updatedItem, selectOneSQL, itemIDToUpdate)
	if err != nil {
		log.Fatalf("Error reading updated record with ID=%d: %v", itemIDToUpdate, err)
	}
	fmt.Printf("Updated record: ID: %d, Name: %s, Quantity: %d\n", updatedItem.ID, updatedItem.Name, updatedItem.Quantity)


	fmt.Println("\n--- Deleting record with ID = 2 ---")
	deleteSQL := "DELETE FROM items WHERE id = ?"
	itemIDToDelete := 2

	result, err = db.Exec(deleteSQL, itemIDToDelete)
	if err != nil {
		log.Fatalf("Error deleting record with ID=%d: %v", itemIDToDelete, err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		log.Printf("Warning: Could not get rows affected after DELETE: %v", err)
	} else {
		fmt.Printf("Record with ID=%d successfully deleted. Rows affected: %d\n", itemIDToDelete, rowsAffected)
	}

	fmt.Println("\n--- Counting remaining records ---")
	selectCountSQL := "SELECT COUNT(*) FROM items"
	var count int

	err = db.Get(&count, selectCountSQL)
	if err != nil {
		log.Fatalf("Error counting records: %v", err)
	}
	fmt.Printf("Total records in table after deletion: %d\n", count)

	fmt.Println("\n--- Performing cleanup (dropping table and VACUUM) ---")

	// _, err = db.Exec("DROP TABLE IF EXISTS items")
	// if err != nil {
	// 	log.Printf("Warning: Error dropping 'items' table: %v", err)
	// } else {
	// 	fmt.Println("'items' table dropped.")
	// }

    // _, err = db.Exec("VACUUM")
    // if err != nil {
    //     log.Printf("Warning: Error executing VACUUM: %v", err)
    // } else {
    //     fmt.Println("VACUUM executed.")
    // }

    // _, err = db.Exec("PRAGMA optimize")
    // if err != nil {
    //     log.Printf("Warning: Error executing PRAGMA optimize: %v", err)
    // } else {
    //     fmt.Println("PRAGMA optimize executed.")
    // }

	fmt.Println("\nProgram finished successfully.")
}