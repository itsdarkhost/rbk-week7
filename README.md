# rbk-week3

Weather + Users API на Go + PostgreSQL.

## Запуск

```bash
docker compose up -d
go run ./cmd/api
```

По умолчанию сервер слушает `:8080`, PostgreSQL берется из `.env`.

Миграции лежат в `migrations/` и применяются автоматически при старте API

## REST API

Пользователи:

```http
POST   /users
GET    /users
GET    /users/{id}
PUT    /users/{id}
DELETE /users/{id}
```

Тело для создания/обновления пользователя:

```json
{
  "username": "sultan"
}
```

Города пользователя:

```http
POST   /users/{id}/cities
GET    /users/{id}/cities
DELETE /users/{id}/cities/{city_id}
```

Body для добавления города:

```json
{
  "city": "Almaty"
}
```

Погода и история:

```http
GET /users/{id}/weather
GET /users/{id}/weather/history?city=Almaty&limit=10&offset=0
```