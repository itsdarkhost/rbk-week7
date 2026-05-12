# rbk-week5

Weather + Users API на Go + PostgreSQL с JWT-аутентификацией.

## Запуск

```bash
docker compose up -d
JWT_SECRET=dev-secret go run ./api
```

По умолчанию сервер слушает `:8080`, PostgreSQL берется из `.env`.
Для запуска API нужен `JWT_SECRET`.

Миграции лежат в `migrations/` и применяются автоматически при старте API

## REST API

Аутентификация:

```http
POST /auth/register
POST /auth/login
```

Register body:

```json
{
  "username": "sultan",
  "email": "sultan@example.com",
  "password": "secret"
}
```

Login body:

```json
{
  "email": "sultan@example.com",
  "password": "secret"
}
```

Все маршруты ниже требуют header:

```http
Authorization: Bearer <access_token>
```

Текущий пользователь:

```http
GET /users/me
```

Города текущего пользователя:

```http
POST   /cities
GET    /cities
DELETE /cities/{city_id}
```

Body для добавления города:

```json
{
  "city": "Almaty"
}
```

Погода и история:

```http
GET /weather
GET /weather/history?city=Almaty&limit=10&offset=0
```

Admin-only маршруты:

```http
GET    /users
GET    /users/{id}
DELETE /users/{id}
```
