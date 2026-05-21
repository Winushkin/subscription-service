# Subscriptions Service API

REST API сервис для агрегации данных об онлайн подписках пользователей. Позволяет управлять подписками, отслеживать расходы и анализировать затраты на подписки за определённые периоды.

---
### Запуск приложения

```bash
# Клонируй репозиторий
git clone https://github.com/Winushkin/subscription-service.git
cd subscription-service

cp .env.example .env

# Запусти с помощью Docker Compose
make buildup

# Приложение будет доступно на http://localhost
```

### Проверка здоровья сервиса

```bash
curl http://localhost/health
# Ответ: {"status":"ok"}
```

---

## API Документация

### Swagger UI

Интерактивная документация доступна на:
```
http://localhost:8080/swagger/index.html
```


### Основные endpoints

#### 1. Создание подписки

```bash
curl -X POST http://localhost/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

**Ответ (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": null,
  "created_at": "2026-05-18T10:30:45.123Z",
  "updated_at": "2026-05-18T10:30:45.123Z"
}
```

#### 2. Получение подписки

```bash
curl http://localhost/api/v1/subscriptions/{id}
```

#### 3. Список подписок

```bash
# Все подписки
curl http://localhost/api/v1/subscriptions

# С фильтром по пользователю
curl "http://localhost/api/v1/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"

# С пагинацией
curl "http://localhost:8080/api/v1/subscriptions?limit=10&offset=0"
```

#### 4. Обновление подписки

```bash
curl -X PUT http://localhost/api/v1/subscriptions/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "price": 500,
    "service_name": "Yandex Plus Premium"
  }'
```

#### 5. Удаление подписки

```bash
curl -X DELETE http://localhost/api/v1/subscriptions/{id}
# Ответ: 204 No Content
```

#### 6. Расчёт стоимости за период

```bash
curl -X POST http://localhost/api/v1/reports/cost \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "service_name": "Yandex Plus",
    "start_date": "06-2025",
    "end_date": "08-2025"
  }'
```

**Ответ:**
```json
{
  "total_cost": 1200,
  "count": 3,
  "currency": "RUB"
}
```

---
### Архитектурные слои

```
HTTP Requests
    ↓
[Gin Router & Middleware]
    ↓
[Handlers] ← Обработка HTTP запросов
    ↓
[Repository Interface] ← Абстракция для хранилища
    ↓
[PostgreSQL Adapter] ← Конкретная реализация
    ↓
[pgxpool.Pool] ← Пул соединений к БД
    ↓
[PostgreSQL Database]
```

---

## Технологический стек

### Backend

| Компонент | Версия | Назначение |
|-----------|--------|-----------|
| **Go** | 1.21+ | Язык программирования |
| **Gin** | v1.9+ | Web фреймворк (HTTP router) |
| **pgx** | v5+ | PostgreSQL драйвер |
| **pgxpool** | v5+ | Пул соединений к БД |
| **Squirrel** | v1+ | Query builder для SQL |
| **zap** | v1+ | Структурированное логирование |
| **swagger** | v1+ | Генерация API документации |

### База данных

| Компонент | Версия |
|-----------|--------|
| **PostgreSQL** | 17-alpine |
| **Миграции** | SQL scripts |

### DevOps

| Компонент | Назначение |
|-----------|-----------|
| **Docker** | Контейнеризация |
| **Docker Compose** | Оркестрация сервисов |

---


### 6. Swagger Documentation

**Зачем:** Автоматическая генерация API документации из кода.

```go
// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body entities.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} entities.Subscription
// @Failure 400 {object} entities.ErrorResponse
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) { ... }
```

**Как генерировать:**
```bash
swag init -g cmd/main.go
```

**Доступно на:**
- Swagger UI: `http://localhost:8080/swagger/index.html`
- JSON схема: `http://localhost:8080/swagger.json`

### 7. Миграции БД

**Зачем:** Версионирование структуры БД и простая инициализация.

```sql
-- migrations/001_create_subscriptions_table.sql
-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date VARCHAR(7) NOT NULL, 
    end_date VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id ON subscriptions(user_id);
CREATE INDEX idx_start_date ON subscriptions(start_date);
CREATE INDEX idx_service_name ON subscriptions(service_name);


-- +goose Down
DROP INDEX IF EXISTS idx_service_name;
DROP INDEX IF EXISTS idx_start_date;
DROP INDEX IF EXISTS idx_user_id;
DROP TABLE IF EXISTS subscriptions;
```

Миграции автоматически выполняются при запуске Docker Compose.

### 8. Context в Go

**Зачем:** Управление таймаутами, отмены операций, передача значений между горутинами.

```go
func (h *Handler) CreateSubscription(c *gin.Context) {
    // Создаём контекст с таймаутом 10 секунд
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()
    
    // Если операция дольше 10 сек - она будет отменена
    sub, err := h.repo.CreateSubscription(ctx, req)
}
```

---
## Тестирование

```bash
go test ./... -v
```
