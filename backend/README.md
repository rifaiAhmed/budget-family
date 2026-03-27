# Family Budget Control Backend

Golang (1.22+), Gin, PostgreSQL, GORM, JWT, zap, viper, golang-migrate.

## Quick Start

### 1) Run dependencies

```bash
docker compose up -d postgres redis
```

### 2) Run migrations

Using the compose migration tool:

```bash
docker compose --profile tools run --rm migrate up
```

### 3) Run the API

Locally:

```bash
go run ./cmd
```

Or with Docker:

```bash
docker compose up --build app
```

Server default: `:8080`

## Environment Variables

- `SERVER_MODE` (debug|release)
- `SERVER_ADDRESS` (e.g. `:8080`)
- `DATABASE_DSN` (Postgres DSN)
- `AUTH_JWT_SECRET`
- `AUTH_ACCESS_TTL` (e.g. `15m`)
- `AUTH_REFRESH_TTL` (e.g. `720h`)
- `AUTH_PASSWORD_COST` (bcrypt cost, e.g. `12`)
- `AUTH_ISSUER` (e.g. `budget-family`)
- `REDIS_ENABLED` (true|false)
- `REDIS_ADDR` (e.g. `localhost:6379`)

## API Overview

### Auth

- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/me`

### Family

- `POST /families`
- `GET /families`
- `POST /families/invite`

### Wallets

- `GET /wallets?family_id=...&page=1&limit=20`
- `POST /wallets`
- `PUT /wallets/:id`
- `DELETE /wallets/:id`

### Categories

- `GET /categories?family_id=...&type=expense`
- `POST /categories`

### Transactions

- `POST /transactions`
- `GET /transactions?family_id=...&from=YYYY-MM-DD&to=YYYY-MM-DD&wallet_id=...&category_id=...&type=income`
- `GET /transactions/summary?family_id=...&from=YYYY-MM-DD&to=YYYY-MM-DD`
- `GET /transactions/:id`
- `PUT /transactions/:id`
- `DELETE /transactions/:id`

### Budgets

- `POST /budgets`
- `GET /budgets?family_id=...&month=...&year=...`
- `GET /budgets/usage?family_id=...&month=...&year=...`

### Goals

- `POST /goals`
- `GET /goals?family_id=...`
- `PUT /goals/:id`

### Bills

- `POST /bills`
- `GET /bills?family_id=...`

## Example Requests

### Register

```bash
curl -s -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"John","email":"john@example.com","phone":"","password":"password123"}'
```

### Login

```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"john@example.com","password":"password123"}'
```

Save `access_token` and use it:

```bash
export TOKEN=YOUR_ACCESS_TOKEN
```

### Create Family

```bash
curl -s -X POST http://localhost:8080/families \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"My Family"}'
```

### Create Wallet

```bash
curl -s -X POST http://localhost:8080/wallets \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"family_id":"FAMILY_UUID","name":"Cash","type":"cash","balance":"100000"}'
```

### Create Expense Transaction

```bash
curl -s -X POST http://localhost:8080/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"family_id":"FAMILY_UUID","wallet_id":"WALLET_UUID","category_id":"CATEGORY_UUID","amount":"50000","type":"expense","note":"groceries","transaction_date":"2026-03-01"}'
```

### Budget Usage

```bash
curl -s "http://localhost:8080/budgets/usage?family_id=FAMILY_UUID&month=3&year=2026" \
  -H "Authorization: Bearer $TOKEN"
```
