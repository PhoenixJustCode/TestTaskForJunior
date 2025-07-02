package psql

import (
	"TestTaskForJun/pkg/database"
	// envconfig "github.com/kelseyhightower/envconfig"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)



// Получить книгу по ID
func GetBookByID(db *database.DB) http.HandlerFunc {
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
func GetAllBooks(db *database.DB) http.HandlerFunc {
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
func CreateBook(db *database.DB) http.HandlerFunc {
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
		log.Info("book was created")
	}
}

// Обновить книгу
func UpdateBook(db *database.DB) http.HandlerFunc {
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
		log.Info("book was update")
	}
}

// Удалить книгу
func DeleteBook(db *database.DB) http.HandlerFunc {
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
		log.Printf("book was deleted with ID %d", id)
	}
}

