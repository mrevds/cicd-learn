# CICD Learn - Book API

Простое REST API для управления книгами с использованием Go, Gin, GORM и PostgreSQL.

## Функциональность

- **POST /books** - Создание новой книги
- **GET /books** - Получение списка всех книг
- **GET /ping** - Проверка работоспособности
- **GET /newping** - Дополнительная проверка

## Модель Book

```json
{
  "id": 1,
  "title": "Название книги",
  "author": "Автор книги", 
  "isbn": "978-3-16-148410-0",
  "created_at": "2025-09-20T10:00:00Z",
  "updated_at": "2025-09-20T10:00:00Z"
}
```

## Конфигурация базы данных

Приложение автоматически определяет хост PostgreSQL:
- Сначала пытается подключиться к **postgres** (для Docker окружения)
- Если не удается, подключается к **localhost** (для локального запуска)

Параметры подключения:
- **Port:** 5432
- **User:** postgres  
- **Password:** postgres
- **Database:** cicd_learn

## Запуск с Docker Compose

```bash
# Запуск PostgreSQL и приложения
docker-compose up -d

# Остановка
docker-compose down
```

## Локальный запуск

```bash
# Установка зависимостей
go mod tidy

# Компиляция
go build -o app .

# Запуск (убедитесь, что PostgreSQL запущен)
./app
```

## Примеры использования API

### Создание книги

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Война и мир",
    "author": "Лев Толстой",
    "isbn": "978-5-17-084161-1"
  }'
```

### Получение всех книг

```bash
curl http://localhost:8080/books
```

### Проверка работоспособности

```bash
curl http://localhost:8080/ping
curl http://localhost:8080/newping
```

## Технический стек

- **Go** - язык программирования
- **Gin** - веб-фреймворк
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - база данных
- **Docker** - контейнеризация