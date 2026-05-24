# rbk-week5

Weather + Users API на Go + PostgreSQL с JWT-аутентификацией, разделенный на два сервиса:

- `api-service` - REST API, бизнес-логика, PostgreSQL, миграции, repository/service layers, DI, JWT, structured logging.
- `gateway-service` - HTTP gateway для внешнего weather API. Не работает с БД, использует `http.Client` с timeout и общается с API service по внутреннему адресу Docker Compose.

## Структура

```text
.
├── docker-compose.yml
├── api-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── cmd/api/
│   ├── internal/
│   ├── migrations/
│   └── tests/
└── gateway-service/
    ├── Dockerfile
    ├── go.mod
    ├── cmd/gateway/
    ├── internal/
    └── tests/
```

## Запуск

Создайте `.env` по примеру:

```bash
cp .env.example .env
```

Запуск всех сервисов:

```bash
docker compose up --build
```

После запуска доступны:

```http
GET http://localhost:8080/health
GET http://localhost:8081/health
GET http://localhost:8081/api-health
```

`api-service` вызывает `gateway-service` для получения погоды через `WEATHER_API_URL=http://gateway-service:8081`.
`gateway-service` ходит во внешний `https://wttr.in` и может проверить API service через `API_SERVICE_URL=http://api-service:8080`.

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

## Gateway API

```http
GET /health
GET /api-health
GET /weather?city=Almaty
```

## Тесты

```bash
(cd api-service && go test ./...)
(cd gateway-service && go test ./...)
```

CI настроен в `.github/workflows/ci.yml`: запускает тесты обоих сервисов и проверяет сборку Docker-образов.
