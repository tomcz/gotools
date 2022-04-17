BASE_DIR ?= $(shell git rev-parse --show-toplevel 2>/dev/null)

.PHONY: precommit
precommit: generate format build

.PHONY: build
build: lint test

.PHONY: generate
generate:
ifeq (, $(shell which genny))
	go install github.com/cheekybits/genny@latest
endif
	go generate ./...

.PHONY: format
format:
ifeq (, $(shell which goimports))
	go install golang.org/x/tools/cmd/goimports@latest
endif
	goimports -w -local github.com/tomcz/gotools .

.PHONY: lint
lint:
ifeq (, $(shell which staticcheck))
	go install honnef.co/go/tools/cmd/staticcheck@latest
endif
	staticcheck ./...

.PHONY: test
test:
ifdef DB_HOST
	go test -cover -tags=integration ./...
else
	go test -cover ./...
endif

COMPOSE_CMD := docker-compose -p gotools -f scripts/docker-compose.yml

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
		-t golang:1.18                  \
		make test
