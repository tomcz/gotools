BASE_DIR ?= $(shell git rev-parse --show-toplevel 2>/dev/null)

.PHONY: build
build: lint test

.PHONY: lint
lint:
 ifeq (, $(shell which staticcheck))
	go install honnef.co/go/tools/cmd/staticcheck@latest
 endif
	staticcheck ./...

.PHONY: test
test:
	go test -cover ./...

.PHONY: docker-test-run
docker-test-run:
	@docker run --rm                    \
		-e WAIT_DB_HOST="database"      \
		-e WAIT_DB_USER="root"          \
		-e WAIT_DB_PASSWORD="sekret"    \
		--network gotools_local         \
		-v "${BASE_DIR}:/code"          \
		-w /code                        \
		-t mysql:5.7                    \
		./wait-for-mysql.sh
	@docker run --rm                    \
		-e DB_HOST="database"           \
		-e DB_DATABASE="gotools_test"   \
		-e DB_USER="root"               \
		-e DB_PASSWORD="sekret"         \
		--network gotools_local         \
		-v "${BASE_DIR}:/code"          \
		-w /code                        \
		-t golang:1.16                  \
		make lint test

.PHONY: docker-test
docker-test:
	@docker-compose up -d
	@bash -c "trap 'docker-compose down' EXIT; $(MAKE) docker-test-run"
