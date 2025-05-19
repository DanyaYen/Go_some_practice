package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Completed bool   `json:"completed" db:"completed"`
}

func setupRouter(db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/todos", func(c *gin.Context) {
        db := c.MustGet("db").(*sqlx.DB)

		var todos []Todo
		selectSQL := "SELECT id, title, completed FROM todos"

		err := db.Select(&todos, selectSQL)
		if err != nil {
			log.Printf("Error fetching todos: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
			return
		}

		c.JSON(http.StatusOK, todos)
	})

	r.POST("/todos", func(c *gin.Context) {
        db := c.MustGet("db").(*sqlx.DB)

		var input Todo
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		insertSQL := "INSERT INTO todos (title, completed) VALUES (?, ?)"

		result, err := db.Exec(insertSQL, input.Title, input.Completed)
		if err != nil {
			log.Printf("Error inserting todo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last insert ID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get created todo ID"})
			return
		}

		createdTodo := Todo{ID: int(id), Title: input.Title, Completed: input.Completed}

		c.JSON(http.StatusCreated, createdTodo)
	})

	r.GET("/todos/:id", func(c *gin.Context) {
        db := c.MustGet("db").(*sqlx.DB)

		idSTR := c.Param("id")
		id, err := strconv.Atoi(idSTR)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo ID format"})
			return
		}

		var todo Todo
		selectOneSQL := "SELECT id, title, completed FROM todos WHERE id = ?"

		err = db.Get(&todo, selectOneSQL, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			log.Printf("Error fetching todo with ID %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todo"})
			return
		}

		c.JSON(http.StatusOK, todo)
	})

	r.PUT("/todos/:id", func(c *gin.Context) {
        db := c.MustGet("db").(*sqlx.DB)

		idSTR := c.Param("id")
		id, err := strconv.Atoi(idSTR)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo ID format"})
			return
		}

        var existingTodo Todo
        err = db.Get(&existingTodo, "SELECT id FROM todos WHERE id = ?", id)
        if err != nil {
            if err == sql.ErrNoRows {
                c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
                return
            }
            log.Printf("Error checking existence of todo with ID %d: %v", id, err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check todo existence"})
            return
        }

		var input Todo
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		updateSQL := "UPDATE todos SET title = ?, completed = ? WHERE id = ?"

		result, err := db.Exec(updateSQL, input.Title, input.Completed, id)
		if err != nil {
			log.Printf("Error updating todo with ID %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
			return
		}

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            log.Printf("Warning: Could not get rows affected after update for ID %d: %v", id, err)
        } else if rowsAffected == 0 {
             c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found after check"})
             return
        }


		var updatedTodo Todo
		err = db.Get(&updatedTodo, "SELECT id, title, completed FROM todos WHERE id = ?", id)
		if err != nil {
			log.Printf("Error fetching updated todo with ID %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated todo"})
			return
		}

		c.JSON(http.StatusOK, updatedTodo)
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
        db := c.MustGet("db").(*sqlx.DB)

		idSTR := c.Param("id")
		id, err := strconv.Atoi(idSTR)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo ID format"})
			return
		}

        var existingID int
        err = db.Get(&existingID, "SELECT id FROM todos WHERE id = ?", id)
        if err != nil {
            if err == sql.ErrNoRows {
                c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
                return
            }
            log.Printf("Error checking existence of todo with ID %d before delete: %v", id, err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check todo existence before delete"})
            return
        }

		deleteSQL := "DELETE FROM todos WHERE id = ?"

		result, err := db.Exec(deleteSQL, id)
		if err != nil {
			log.Printf("Error deleting todo with ID %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
			return
		}

        rowsAffected, err := result.RowsAffected()
         if err != nil {
            log.Printf("Warning: Could not get rows affected after delete for ID %d: %v", id, err)
        } else if rowsAffected == 0 {
             log.Printf("Warning: Delete query for ID %d affected 0 rows", id)
        }

		c.Status(http.StatusNoContent)
	})

	return r
}

func main() {
	dbPath := "./todos.db"

    fmt.Printf("Connecting to database at path: %s\n", dbPath)
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

	CreateTableSQL := `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        completed INTEGER NOT NULL
    );`

	fmt.Println("Creating 'todos' table (if not exists)...")
	_, err = db.Exec(CreateTableSQL)
	if err != nil {
		log.Fatalf("Error creating 'todos' table: %v", err)
	}
	fmt.Println("'todos' table created or already exists.")


	r := setupRouter(db)

	fmt.Println("Starting Gin server on :8080...")
	err = r.Run()
    if err != nil {
        log.Fatalf("Error running Gin server: %v", err)
    }
}