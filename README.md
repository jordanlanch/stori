# 🔍 Stori Test

This service processes transaction files, calculates summary information, and sends the summary via email. It also implements rate limiting to control the number of requests.

## 📋 Endpoints

### Process Transactions
```bash
curl --location 'http://localhost:8080/process-transactions' \
--header 'Content-Type: application/json'
```

## 💻 Requirements
- **Port**: 8080 - REST
- **Tools**:
  - `make`
  - `docker` version 20.10.21
  - `docker-compose` version 1.29.2

## 🚀 Run the App

1. **Copy environment variables file**:
   ```sh
   cp .env.example .env
   ```

2. **Run the application**:
   ```sh
   ./run-dev.sh
   ```

3. **Run database migrations** (in a new shell):
   ```sh
   make migrate-up
   ```

## 🏛️ Architecture Layers

- **Router**: Defines API routes and associates them with controllers.
- **Controller**: Handles HTTP requests and responses.
- **Usecase**: Contains business logic and application rules.
- **Repository**: Manages data persistence and retrieval.
- **Domain**: Defines core business entities and logic.

## 🧪 Run Unit Tests and Integration Tests

### Unit Tests
```sh
make unit_test
```

### End-to-End (E2E) Tests
```sh
make e2e_test
```

### All Tests (Unit and E2E)
```sh
make test
```

## 📜 Environment Variables

Ensure you have the following variables set in your `.env` file:

```ini
APP_ENV=test
SERVER_ADDRESS=:8080
CONTEXT_TIMEOUT=5
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
EMAIL_FROM=test-email@example.com
EMAIL_TO=test-recipient@example.com
EMAIL_PASSWORD=test-email-password
SMTP_HOST=test.smtp.example.com
SMTP_PORT=587
CSV_FILE_PATH=/app/test/transactions.csv
FAKE_EMAIL=true
RATE_LIMIT=1000
REDIS_TIMEOUT_SEC=5
CACHE_DURATION_SEC=600
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres_password
DB_NAME=stori_test_db
DB_PORT=5432
```

## 🛠️ Makefile Commands

### Migrations

#### Create a new migration
```sh
make migration name=<migration_name>
```

#### Create a new Go migration
```sh
make migration-go name=<migration_name>
```

#### Check migration status
```sh
make migrate-status
```

#### Apply all up migrations
```sh
make migrate-up
```

#### Apply seed data
```sh
make migrate-seeds
```

#### Rollback the last migration
```sh
make migrate-down
```

#### Reset all migrations
```sh
make migrate-reset
```

### Mocks

#### Generate mocks
```sh
make mocks
```

## 📂 Project Structure

- `cmd/` - Main application entry points.
- `internal/`
  - `config/` - Configuration management.
  - `core/` - Business logic and domain models.
  - `infrastructure/` - External services and frameworks.
  - `interface/` - API layer and interactions.
- `test/` - Test utilities and end-to-end tests.

## 📝 Notes

- Ensure Redis and PostgreSQL are running before starting the application.
- Use `docker-compose` to easily manage service dependencies.
```
