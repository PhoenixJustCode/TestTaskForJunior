package database

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func GetBookByID(id int, db *DB) (Book, error) {
	query := "SELECT id, title, description FROM books WHERE id = $1"
	row := db.Conn.QueryRow(query, id)

	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.WithFields(log.Fields{
				"handler": "Get book by id book",
				"problem": "book not found",
			}).Error(err)
			return book, err
		}
		return book, err
	}

	return book, nil
}

func GetBooks(db *DB) ([]Book, error) {
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
			log.WithFields(log.Fields{
				"handler": "Get book",
				"problem": "with scan book from DB",
			}).Error(err)
			continue
		}
		books = append(books, b)
	}
	return books, nil
}

func AddBook(book Book, db *DB) error {
	query := "INSERT INTO books (title, description) VALUES ($1, $2)"
	_, err := db.Conn.Exec(query, book.Title, book.Description)
	log.WithFields(log.Fields{
		"handler": "Add book",
		"problem": "with exec sql",
	}).Error(err)
	return err
}

func UpdateBook(book Book, db *DB) error {
	query := "UPDATE books SET title = $1, description = $2 WHERE id = $3"
	_, err := db.Conn.Exec(query, book.Title, book.Description, book.ID)
	log.WithFields(log.Fields{
		"handler": "Update book",
		"problem": "with exec sql",
	}).Error(err)
	return err
}

func DeleteBook(id int, db *DB) error {
	query := "DELETE FROM books WHERE id = $1"
	_, err := db.Conn.Exec(query, id)
	log.WithFields(log.Fields{
		"handler": "Delete book",
		"problem": "with exec sql",
	}).Error(err)
	return err
}
