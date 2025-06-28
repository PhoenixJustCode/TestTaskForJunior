package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"TestTaskForJun/pkg/database"
	_ "TestTaskForJun/docs" // swag docs
	httpSwagger "github.com/swaggo/http-swagger"
)

var db *database.DB

// Получить книгу по ID
func getBookByID(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		book, err := database.GetBookByID(id, db)
		if err != nil {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	}
}

// Получить все книги
func getAllBooks(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		books, err := database.GetBooks(db)
		if err != nil {
			http.Error(w, "Failed to get books", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	}
}

// Добавить книгу
func createBook(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var book database.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := database.AddBook(book, db); err != nil {
			http.Error(w, "Failed to add book", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Book created")
	}
}

// Обновить книгу
func updateBook(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var book database.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := database.UpdateBook(book, db); err != nil {
			http.Error(w, "Failed to update book", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Book updated")
	}
}

// Удалить книгу
func deleteBook(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := database.DeleteBook(id, db); err != nil {
			http.Error(w, "Failed to delete book", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Book with ID %d deleted\n", id)
	}
}

func main() {
	var err error
	data := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err = database.NewDB(data)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	http.HandleFunc("/book/", getBookByID(db))
	http.HandleFunc("/books", getAllBooks(db))
	http.HandleFunc("/create", createBook(db))
	http.HandleFunc("/update", updateBook(db))
	http.HandleFunc("/delete/", deleteBook(db))

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
