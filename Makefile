.PHONY: fmt lint lint-fix test buf-lint buf

fmt:
	golangci-lint fmt ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test ./...

buf-lint:
	buf lint

buf:
	buf lint & buf generate