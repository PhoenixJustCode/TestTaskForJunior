package database

import(
	"database/sql"
)

type Book struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}


type DB struct {
	Conn *sql.DB
}

func (db *DB) Close() {
	db.Conn.Close()
}