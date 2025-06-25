package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

type Book struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DB struct {
	Conn *sql.DB
}

func NewDB(dataSource string) (*DB, error) {
	conn, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return &DB{Conn: conn}, nil
}

func (db *DB) Close() {
	db.Conn.Close()
}

func (db *DB) GetBookByID(id int) (Book, error) {
	query := "SELECT id, title, description FROM books WHERE id = $1"
	row := db.Conn.QueryRow(query, id)

	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return book, fmt.Errorf("book not found")
		}
		return book, err
	}

	return book, nil
}

func (db *DB) GetBooks() ([]Book, error) {
	rows, err := db.Conn.Query("SELECT id, title, description FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Description)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		books = append(books, b)
	}
	return books, nil
}

func (db *DB) AddBook(book Book) error {
	query := "INSERT INTO books (title, description) VALUES ($1, $2)"
	_, err := db.Conn.Exec(query, book.Title, book.Description)
	return err
}

func (db *DB) UpdateBook(book Book) error {
	query := "UPDATE books SET title = $1, description = $2 WHERE id = $3"
	_, err := db.Conn.Exec(query, book.Title, book.Description, book.ID)
	return err
}

func (db *DB) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id = $1"
	_, err := db.Conn.Exec(query, id)
	return err
}
