package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID			string `json:id`
	Title 		string `json:title`
	Author 		string `json:author`
	Quantity	int	   `json:quantity`
}


var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func findBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil;
		}
	}

	return nil, errors.New("book not found")
}

func getAllBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books);
}

func getBookHandler(c *gin.Context) {
	id := c.Param("id");
	book, err := findBookById(id);

	if err != nil {
		fmt.Println(err);
		c.JSON(http.StatusNotFound, gin.H{"error":"book not found"})
		return 
	}

	c.JSON(http.StatusFound, book);
}


func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id");

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid query parameter"})
		return
	}

	book, err := findBookById(id);

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error:":"book not found"})
		return;
	}

	if book.Quantity <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error:":"book not found"});
		return;
	}

	book.Quantity-=1;
	c.JSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		fmt.Println(err);
		c.JSON(http.StatusNotFound, gin.H{"errorMessage":"book not found"})
		return
	}
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default();
	router.GET("/books", getAllBooks);
	router.POST("/books", createBook);
	router.GET("/books/:id", getBookHandler);
	router.PATCH("/books", checkoutBook);
	router.Run("localhost:8080")
}