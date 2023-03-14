BASE_DIR ?= $(shell git rev-parse --show-toplevel 2>/dev/null)

.PHONY: precommit
precommit: format lint test

.PHONY: format
format:
	goimports -w -local github.com/tomcz/gotools .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
ifdef DB_HOST
	go test -race -cover -tags=integration ./...
else
	go test -race -cover ./...
endif

.PHONY: tidy
tidy:
	go mod tidy -compat=1.20

COMPOSE_CMD := docker compose -p gotools -f scripts/docker-compose.yml

.PHONY: docker
docker:
	@${COMPOSE_CMD} up -d
	@bash -c "trap '${COMPOSE_CMD} down' EXIT; $(MAKE) docker-run"

.PHONY: docker-run
docker-run:
	@docker run --rm                    \
		-e WAIT_DB_HOST="database"      \
		-e WAIT_DB_USER="root"          \
		-e WAIT_DB_PASSWORD="sekret"    \
		--network gotools_local         \
		-v "${BASE_DIR}:/code"          \
		-w /code                        \
		-t mysql:8.0                    \
		./scripts/wait-for-mysql.sh
	@docker run --rm                    \
		-e DB_HOST="database"           \
		-e DB_USER="root"               \
		-e DB_PASSWORD="sekret"         \
		--network gotools_local         \
		-v "${BASE_DIR}:/code"          \
		-w /code                        \
		-t golang:1.20                  \
		make test
