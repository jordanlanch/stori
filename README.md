# 🔍 stori_TEST
This service control de notification and limit rate by MessageType

# Enpoint
```
  curl --location 'http://localhost:8080/process-transactions' \
--header 'Content-Type: application/json'
```


# 💻 Requirements
  - Port: [8080] - REST
  - make
  - docker version 20.10.21.
  - docker-compose version 1.29.2.

copy .env.example to .env
# 🚀 Run the app
```sh
$ ./run-dev.sh
```
# Run Migratios after Run App in new Shell

```
  make migrate-up
```

## Architecture Layers of the project

- Router
- Controller
- Usecase
- Repository
- Domain


# Run UNIT Tests (unit test and integration test)

```
  make unit_test
```

# Run E2E Tests

```
  make e2e_test
```

# Run Unit Test and E2E Tests

```
  make test
```
