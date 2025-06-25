# 📚 Book API — CRUD приложение на Go

RESTful CRUD-приложение для управления сущностью **Book**, реализованное на **Go**, с использованием **PostgreSQL**, **Gin** и **Swagger** для документации API.

---

## 🚀 Возможности

- 🔍 Получение книги по ID
- 📋 Получение всех книг
- ➕ Добавление новой книги
- ✏️ Обновление книги
- ❌ Удаление книги
- 📑 Swagger UI документация

---

## 🛠️ Технологии

- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Swaggo](https://github.com/swaggo/gin-swagger) — Swagger документация
- [godotenv](https://github.com/joho/godotenv) — загрузка .env переменных

---

## 📂 Структура проекта

```
TestTaskForJun/
│
├── cmd/ # Точка входа приложения (main.go)
│ └── main.go
│
├── pkg/
│ └── database/
│     └── database.go  # Логика взаимодействия с базой данных (PostgreSQL)
│
├── docs/ # Swagger-документация (генерируется автоматически)
│ ├── docs.go
│ ├── swagger.json
│ └── swagger.yaml
│
├── scheme/ # SQL-скрипты для инициализации БД
│ └── init.up.sql
│
├── .env # Переменные окружения (не коммитить!)
├── .gitignore # Исключения Git
├── go.mod # Модуль Go и зависимости
├── go.sum # Контрольные суммы зависимостей
└── README.md # Инструкция по запуску проекта
```

---

## 📥 Установка и запуск

### 🔧 Клонирование проекта

```bash
git clone https://github.com/yourusername/TestTaskForJun.git
cd TestTaskForJun
```

### ⚙️ Создание `.env`

Создай `.env` файл в корне проекта:

```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_NAME=bookTest_db
DB_PASSWORD=yourpassword
```

### 📦 Установка зависимостей

```bash
go mod tidy
```

### 🐘 Настройка базы данных PostgreSQL

Создай таблицу `books`: # SQL code in scheme/init.up.sql

```sql
CREATE TABLE
    books (
        id serial not null unique,
        title varchar(255) not null,
        description varchar(255) not null
    );
```

### 📚 Генерация Swagger-документации

```bash
swag init -g cmd/main.go
```

---

## ▶️ Запуск

```bash
go run ./cmd/main.go
```

---

## 📬 API Endpoints

| Метод  | Endpoint      | Описание             |
| ------ | ------------- | -------------------- |
| GET    | `/book/:id`   | Получить книгу по ID |
| GET    | `/books`      | Получить все книги   |
| POST   | `/create`     | Создать новую книгу  |
| PUT    | `/update`     | Обновить книгу       |
| DELETE | `/delete/:id` | Удалить книгу по ID  |

---

## 🔎 Swagger UI

Swagger документация доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

---

## 📌 Пример запроса (POST /create)

```json
{
  "title": "Грокаем алгоритмы",
  "description": "Пошаговое руководство по алгоритмам"
}
```

---

## 🤝 Автор

**Alexander Demchenko**  
GitHub: [@PhoenixJustCode](https://github.com/PhoenixJustCode)  
Telegram: [@phoenix_S_E](https://t.me/phoenix_S_E)
