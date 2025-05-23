package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string `json:"title" gorm:"not null"`
	Content string `json:"content" gorm:"not null"`
}

type Handler struct {
	DB *gorm.DB
}

func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Post{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}

func (h *Handler) getAllPosts(c echo.Context) error {
	var posts []Post
	result := h.DB.Find(&posts)
	if result.Error != nil {
		c.Logger().Errorf("Database error fetching all posts: %v", result.Error)
		return c.String(http.StatusInternalServerError, "Failed to fetch posts")
	}
	return c.JSON(http.StatusOK, posts)
}

func (h *Handler) getPostByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post ID format")
	}

	var post Post
	result := h.DB.First(&post, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "Post not found")
		}
		c.Logger().Errorf("Database error fetching post %d: %v", id, result.Error)
		return c.String(http.StatusInternalServerError, "Failed to fetch post")
	}

	return c.JSON(http.StatusOK, post)
}

func (h *Handler) createPost(c echo.Context) error {
	post := new(Post)
	if err := c.Bind(post); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON body")
	}

	if post.Title == "" || post.Content == "" {
		return c.String(http.StatusBadRequest, "Title and Content cannot be empty")
	}

	result := h.DB.Create(post)
	if result.Error != nil {
		c.Logger().Errorf("Database error creating post: %v", result.Error)
		return c.String(http.StatusInternalServerError, "Failed to create post")
	}

	return c.JSON(http.StatusCreated, post)
}

func (h *Handler) updatePost(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post ID format")
	}

	var post Post
	result := h.DB.First(&post, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "Post not found")
		}
		c.Logger().Errorf("Database error finding post %d for update: %v", id, result.Error)
		return c.String(http.StatusInternalServerError, "Failed to find post for update")
	}

	if err := c.Bind(&post); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON body")
	}

	if post.Title == "" || post.Content == "" {
		return c.String(http.StatusBadRequest, "Title and Content cannot be empty")
	}

	result = h.DB.Save(&post)
	if result.Error != nil {
		c.Logger().Errorf("Database error updating post %d: %v", id, result.Error)
		return c.String(http.StatusInternalServerError, "Failed to update post")
	}

	return c.JSON(http.StatusOK, post)
}

func (h *Handler) deletePost(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post ID format")
	}

	result := h.DB.Delete(&Post{}, id)

	if result.Error != nil {
		c.Logger().Errorf("Database error deleting post %d: %v", id, result.Error)
		return c.String(http.StatusInternalServerError, "Failed to delete post")
	}

	if result.RowsAffected == 0 {
		return c.String(http.StatusNotFound, "Post not found")
	}

	return c.NoContent(http.StatusNoContent)
}

func setupRoutes(e *echo.Echo, h *Handler) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/posts", h.getAllPosts)
	e.GET("/posts/:id", h.getPostByID)
	e.POST("/posts", h.createPost)
	e.PUT("/posts/:id", h.updatePost)
	e.DELETE("/posts/:id", h.deletePost)
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	handler := &Handler{DB: db}

	e := echo.New()

	setupRoutes(e, handler)

	e.Logger.Fatal(e.Start(":8080"))
}