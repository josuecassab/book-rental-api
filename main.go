package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type book struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Genre    string `json:"genre"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{Id: "1", Title: "The art of seduction", Genre: "non-fiction", Author: "Robert greene", Quantity: 5},
	{Id: "2", Title: "Sapiens", Genre: "non-fiction", Author: "Yuval Noah Harari", Quantity: 6},
	{Id: "3", Title: "El infinito en un junco", Genre: "non-fiction", Author: "Irene Vallejo", Quantity: 3},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {

	var newBook map[string]interface{}
	byte_body, _ := io.ReadAll(c.Request.Body)

	err := json.Unmarshal(byte_body, &newBook)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not decode body"})
		return
	}
	fmt.Println(string(byte_body))
	// c.BindJSON(&newBook)
	fmt.Println(newBook)

	quantityInt, _ := strconv.Atoi(newBook["quantity"].(string))
	newBookStruct := book{
		Id:       fmt.Sprintf("%d", len(books)+1),
		Title:    newBook["title"].(string),
		Genre:    newBook["genre"].(string),
		Author:   newBook["author"].(string),
		Quantity: quantityInt,
	}

	books = append(books, newBookStruct)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func fetchBookById(id string) (*book, error) {
	for i, b := range books {
		if b.Id == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("not id found")
}

func getBookById(c *gin.Context) {
	id := c.Param("id")

	book, err := fetchBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {

	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad query parameter"})
		return
	}

	book, err := fetchBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book id not found"})
		return
	}

	if book.Quantity < 1 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available to rent"})
		return
	}
	fmt.Println(book)
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

func returnBook(c *gin.Context) {

	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad query parameter"})
		return
	}

	book, err := fetchBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book id not found"})
		return
	}

	book.Quantity += 1

	c.IndentedJSON(http.StatusOK, book)
}

func renderHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "welcome to book library",
		"books": books})
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.GET("/", renderHTML)
	r.GET("/books", getBooks)
	r.GET("/books/:id", getBookById)
	r.POST("/books", createBook)
	r.PUT("/checkout/", checkoutBook)
	r.PUT("/return/", returnBook)
	r.Run()
}
