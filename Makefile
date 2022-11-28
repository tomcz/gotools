BASE_DIR ?= $(shell git rev-parse --show-toplevel 2>/dev/null)

.PHONY: precommit
precommit: format build

.PHONY: build
build: lint test

.PHONY: format
format:
	${BASE_DIR}/scripts/format.sh

.PHONY: lint
lint:
	${BASE_DIR}/scripts/lint.sh

.PHONY: test
test:
ifdef DB_HOST
	go test -cover -tags=integration ./...
else
	go test -cover ./...
endif

.PHONY: tidy
tidy:
	go mod tidy -compat=1.19

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
		-t mysql:5.7                    \
		./scripts/wait-for-mysql.sh
	@docker run --rm                    \
		-e DB_HOST="database"           \
		-e DB_DATABASE="gotools_test"   \
		-e DB_USER="root"               \
		-e DB_PASSWORD="sekret"         \
		--network gotools_local         \
		-v "${BASE_DIR}:/code"          \
		-w /code                        \
		-t golang:1.19                  \
		make test
