ifeq ($(POSTGRES_SETUP_TEST),)
    POSTGRES_SETUP_TEST := user=postgres password=1111 dbname=TestAvito host=localhost port=5432 sslmode=disable
endif

ifeq ($(POSTGRES_SETUP),)
    POSTGRES_SETUP := user=postgres password=1111 dbname=Avito host=db port=5432 sslmode=disable
endif


INTERNAL_PKG_PATH=$(CURDIR)/internal/storage
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

MOCKGEN_TAG=1.2.0
DOCKER_COMPOSE_FILE := docker-compose.yml

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up


.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down


.PHONY: test-migration-test-up
test-migration-test-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-test-down
test-migration-test-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down


.PHONY: .generate-mockgen-deps
.generate-mockgen-deps:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	@GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@$(MOCKGEN_TAG)
endif

.PHONY: .generate-mockgen
.generate-mockgen:
	PATH="$(LOCAL_BIN):$$PATH" go generate -x -run=mockgen ./...

.PHONY: gofmt
gofmt:
	goimports -l -w $(CURDIR)

.test:
	$(info Running tests...)
	go test ./...

.PHONY: integration-test
integration-test:
	go test -tags=integration -v ./tests

.PHONY: unit-test
unit-test:
	go test -v ./...

.PHONY: unit-test-coverage
unit-test-coverage:
	go test -v ./... -coverprofile=coverage.out

.PHONY: docker-up
docker-up:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: docker-down
docker-down:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
