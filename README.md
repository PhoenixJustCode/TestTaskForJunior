# 📚 Book API — CRUD-приложение на Go с Docker и Swagger

Полноценное RESTful-приложение для управления книгами, реализованное на Go с использованием PostgreSQL, Gin, Swagger и Docker.

---

## 🚀 Основной функционал

- 🔍 Получение книги по ID
- 📋 Получение всех книг
- ➕ Добавление книги
- ✏️ Обновление книги
- ❌ Удаление книги
- 📑 Swagger UI-документация

---

## 🛠️ Технологии

- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Gin](https://github.com/gin-gonic/gin) — фреймворк для HTTP-сервера
- [Swaggo](https://github.com/swaggo/swag) — генерация Swagger-документации
- [Docker & Docker Compose](https://www.docker.com/)
- [Makefile](https://www.gnu.org/software/make/) — автоматизация команд
- `.env` — переменные окружения

---

## 📁 Структура проекта

```
TestTaskForJun/
│
├── cmd/                 # Точка входа (main.go)
├── pkg/database/        # Подключение к PostgreSQL
├── docs/                # Swagger-файлы (генерируются автоматически)
├── scheme/              # SQL-скрипты для инициализации БД
├── .env                 # Переменные окружения (не коммитить!)
├── .gitignore           # Исключения Git
├── Makefile             # Набор команд для сборки и запуска
├── Dockerfile           # Docker-образ для приложения
├── docker-compose.yml   # Сборка и запуск всех контейнеров
├── Task.pdf             # 📄 Файл с описанием задачи
├── go.mod / go.sum      # Зависимости проекта
└── README.md            # Этот файл
```

---

## 📦 Быстрый запуск с Docker

### ⚙️ 1. Создай `.env` файл

```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_NAME=bookTest_db
DB_PASSWORD=mysecretpassword
```

### 🐳 2. Сборка и запуск контейнеров

> Используй встроенные команды через `make`:

```bash
make build   # Собрать образы
make run     # Запустить контейнеры
```

---

## 🗃️ Миграции базы данных

```bash
make migrate
```

> Выполнит SQL-скрипты из `scheme/` для создания таблицы `books`.

---

## 📑 Swagger-документация

Генерация документации:

```bash
make swag
```

Swagger UI будет доступен после запуска по адресу:  
👉 [`http://localhost:8080/swagger/index.html`](http://localhost:8080/swagger/index.html)

---

## 📬 API Endpoints

| Метод  | Endpoint     | Описание                |
|--------|--------------|-------------------------|
| GET    | /book/:id    | Получить книгу по ID    |
| GET    | /books       | Получить все книги      |
| POST   | /create      | Создать книгу           |
| PUT    | /update      | Обновить книгу          |
| DELETE | /delete/:id  | Удалить книгу по ID     |

---

## 🔎 Пример запроса (POST /create)

```json
{
  "title": "Грокаем алгоритмы",
  "description": "Пошаговое руководство по алгоритмам"
}
```

---

## 📄 Task.pdf

Файл `Task.pdf` содержит подробное описание задачи и требований.  
Он находится в корне проекта.

---

## 🤝 Автор

**Alexander Demchenko**  
- GitHub: [@PhoenixJustCode](https://github.com/PhoenixJustCode)  
- Telegram: [@phoenix_S_E](https://t.me/phoenix_S_E)
