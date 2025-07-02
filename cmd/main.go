package main

import (
	_ "TestTaskForJun/docs" // swag docs
	"TestTaskForJun/pkg/database"
	"fmt"
	envconfig "github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"TestTaskForJun/internal/psql"
)

var db *database.DB

type envData struct {
	Host     string
	Port     int
	User     string
	Name     string
	Password string
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	var cfg envData
	err := envconfig.Process("DB", &cfg)
	if err != nil {
		log.Fatal("Ошибка обработки переменных окружения:", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password,
	)

	db, err = database.NewDB(dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	http.HandleFunc("/book/", psql.GetBookByID(db))
	http.HandleFunc("/books", psql.GetAllBooks(db))
	http.HandleFunc("/create", psql.CreateBook(db))
	http.HandleFunc("/update", psql.UpdateBook(db))
	http.HandleFunc("/delete/", psql.DeleteBook(db))

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Info("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
