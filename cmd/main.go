package main

import (
	_ "TestTaskForJun/docs" // swag docs
	"TestTaskForJun/internal/database/psql"
	"TestTaskForJun/internal/database"
	"fmt"
	envconfig "github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"TestTaskForJun/internal/auth"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var db *database.DB

type envData struct {
	Host     string
	Port     int
	User     string
	Name_DB     string
	Password string
	JWTSecret string `envconfig:"JWT_SECRET"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	godotenv.Load()
}

func main() {
	var cfg envData
	err := envconfig.Process("DB", &cfg)
	if err != nil {
		log.Fatal("Ошибка обработки переменных окружения:", err)
	}

	auth.InitJWT(cfg.JWTSecret)
	
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Name_DB, cfg.Password,
	)
	fmt.Println(dsn)

	db, err = database.NewDB(dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()


	r := mux.NewRouter()
	r.Handle("/swagger/", httpSwagger.WrapHandler)
	
	r.HandleFunc("/register", auth.RegisterHandler(db.Conn)).Methods("POST")
	r.HandleFunc("/login", auth.LoginHandler(db.Conn)).Methods("POST")
	r.HandleFunc("/refresh", auth.RefreshHandler(db.Conn)).Methods("POST")

	// защищённые роуты
	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.JWTMiddleware)
	api.HandleFunc("/books", psql.GetAllBooks(db)).Methods("GET")
	api.HandleFunc("/book/{id}", psql.GetBookByID(db)).Methods("GET")
	api.HandleFunc("/create", psql.CreateBook(db)).Methods("POST")
	api.HandleFunc("/update/{id}", psql.UpdateBook(db)).Methods("PUT")
	api.HandleFunc("/delete/{id}", psql.DeleteBook(db)).Methods("DELETE")

	http.Handle("/", r)
	log.Info("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
