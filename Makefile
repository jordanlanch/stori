DB_NAME=stori_test_db
DB_PORT=5432
MIGRATION_DIR=storage/migrations
ROUTE="host=localhost user=postgres password=postgres_password dbname=${DB_NAME} port=${DB_PORT} sslmode=disable"

DB_NAME_TEST=stori_test_db-test
DB_PORT_TEST=5433

define setup_env
    $(eval ENV_FILE := .envtest)
    @echo " - setup env $(ENV_FILE)"
    $(eval include .envtest)
    $(eval export)
endef

## get extra arguments and filter out commands from args
args = $(filter-out $@,$(MAKECMDGOALS))

.PHONY: all test unit_test e2e_test

test:

	echo "Starting test environment"
	$(call setup_env)
	make unit_test
	make e2e_test

unit_test:
	@echo "/////////////////////////////////Starting unit test environment/////////////////////////////////"
	@go test $(shell go list ./... | grep -v /test) -coverprofile coverage.out -covermode count -coverpkg=./internal/core/... && \
	go tool cover -func coverage.out | grep total | awk '{print $3}'

goose_install:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.3

e2e_test: goose_install
	echo "Starting test environment"
	$(call setup_env)
	echo "/////////////////////////////////Deleting fixtures/////////////////////////////////"
	rm -rf ./test/fixtures
	echo "/////////////////////////////////Starting E2E Test/////////////////////////////////"
	docker compose -f docker-compose-test.yml up --build -d
	sh -c 'sleep 5 &&  goose -dir ${MIGRATION_DIR} postgres "host=localhost user=postgres-test password=postgres_password-test dbname=${DB_NAME_TEST} port=${DB_PORT_TEST} sslmode=disable" up'
	cd ./test &&  go test ./... || true
	docker compose -f docker-compose-test.yml down
	echo "/////////////////////////////////Ending E2E Test/////////////////////////////////"


## default that allows accepting extra args
%:
    @:

.PHONY: migration
migration:
	goose -dir ${MIGRATION_DIR} create $(call args,defaultstring) sql
.PHONY: migration
migration-go:
	goose -dir ${MIGRATION_DIR} create $(call args,defaultstring) go

migrate-status:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} status

migrate-up:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} up
migrate-seeds:
	./seeds/goose-custom -dir seeds up -dbstring ${ROUTE}

migrate-down:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} down

migrate-rollback:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} reset

migrate-reset:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} reset

mocks:
	mockery --dir=domain --output=domain/mocks --outpkg=mocks --all
