#!/bin/sh
set -e

echo "🕐 Ждём подключение к PostgreSQL..."
MAX_RETRIES=60
RETRY_COUNT=0

until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "⏳ Postgres is unavailable - sleeping ($RETRY_COUNT/$MAX_RETRIES)"
  sleep 1
  RETRY_COUNT=$((RETRY_COUNT+1))
  if [ "$RETRY_COUNT" -ge "$MAX_RETRIES" ]; then
    echo "❌ Превышено максимальное число попыток подключения к PostgreSQL"
    exit 1
  fi
done

echo "✅ PostgreSQL доступен. Проверка базы данных '$DB_NAME'..."

ATTEMPTS_LEFT=30
until psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q' >/dev/null 2>&1 || [ $ATTEMPTS_LEFT -eq 0 ]; do
  echo "⚠️  База '$DB_NAME' пока не готова. Осталось попыток: $ATTEMPTS_LEFT"
  ATTEMPTS_LEFT=$((ATTEMPTS_LEFT-1))
  sleep 1
done

if [ "$ATTEMPTS_LEFT" -eq 0 ]; then
  echo "❌ База данных '$DB_NAME' не готова. Завершаем выполнение."
  exit 1
fi

echo "🚀 Всё готово. Запускаем приложение: $@"
exec "$@"
