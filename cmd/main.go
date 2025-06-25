// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"TestTaskForJun/pkg/database"
	_ "TestTaskForJun/docs" // swag docs
	httpSwagger "github.com/swaggo/http-swagger"
)

var db *database.DB

// @title Book API
// @version 1.0
// @description Это Тестовое задание на тему CRUD-приложение на Go.
// @host localhost:8080
// @BasePath /

// @Summary Получить книгу по ID
// @Produce json
// @Param id path int true "ID книги"
// @Success 200 {object} database.Book
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Book not found"
// @Router /book/{id} [get]
func getBookByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/book/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	book, err := db.GetBookByID(id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// @Summary Получить все книги
// @Produce json
// @Success 200 {array} database.Book
// @Failure 500 {string} string "Failed to get books"
// @Router /books [get]
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	books, err := db.GetBooks()
	if err != nil {
		http.Error(w, "Failed to get books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// @Summary Добавить книгу
// @Accept json
// @Produce json
// @Param book body database.Book true "Книга"
// @Success 201 {string} string "Book created"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to add book"
// @Router /create [post]
func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var book database.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.AddBook(book); err != nil {
		http.Error(w, "Failed to add book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Book created")
}

// @Summary Обновить книгу
// @Accept json
// @Produce json
// @Param book body database.Book true "Книга"
// @Success 200 {string} string "Book updated"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to update book"
// @Router /update [put]
func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var book database.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.UpdateBook(book); err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Book updated")
}

// @Summary Удалить книгу
// @Produce json
// @Param id path int true "ID книги"
// @Success 200 {string} string "Book deleted"
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Failed to delete book"
// @Router /delete/{id} [delete]
func deleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := db.DeleteBook(id); err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Book with ID %d deleted\n", id)
}

func main() {
	var err error
	db, err = database.NewDB("host=127.0.0.1 port=5432 user=postgres dbname=bookTest_db sslmode=disable")
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	http.HandleFunc("/book/", getBookByID)
	http.HandleFunc("/books", getAllBooks)
	http.HandleFunc("/create", createBook)
	http.HandleFunc("/update", updateBook)
	http.HandleFunc("/delete/", deleteBook)

	http.Handle("/swagger/", httpSwagger.WrapHandler)
	
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}