.PHONY: build
build: clean test

.PHONY: clean
clean:
	rm -rf target

.PHONY: test
test:
	go test -cover ./...
