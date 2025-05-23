package main

import (
	"net/http" 
	"log"    

	"github.com/gin-gonic/gin"        
	"github.com/graphql-go/graphql" 
)

// 1. Определение структуры Book
type Book struct {
	Title  string `json:"title"`  
	Author string `json:"author"` 
}

func main() {
	var bookType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Book", 
			Fields: graphql.Fields{
				"title": &graphql.Field{ 
					Type: graphql.String, 
				},
				"author": &graphql.Field{ 
					Type: graphql.String,
				},
			},
		},
	)

	var books = []Book{
		{Title: "The Hitchhiker's Guide to the Galaxy", Author: "Douglas Adams"},
		{Title: "The Lord of the Rings", Author: "J.R.R. Tolkien"},
	}

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query", 
			Fields: graphql.Fields{
				"allBooks": &graphql.Field{ 
					Type: graphql.NewList(bookType), 
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return books, nil
					},
				},
			},
		},
	)

	var schema, err = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	if err != nil {
		log.Fatalf("failed to create schema: %v", err) 
	}

	r := gin.Default()

	r.POST("/graphql", func(c *gin.Context) {
		var body struct {
			Query string `json:"query"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) 
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: body.Query,
		})


		if len(result.Errors) > 0 {
			c.JSON(http.StatusOK, gin.H{"errors": result.Errors}) 
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": result.Data})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err) 
	}
}