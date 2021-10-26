GIT_COMMIT=$(shell git rev-list -1 HEAD)
GITHASH_COMMIT=$(shell git log --format="%h" -n 1)

.PHONY: test check
test:
	go test -race -coverprofile=coverage.out -timeout 30s ./...
check:
	golangci-lint run
build:
	go build -ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.GitHashCommit=$(GITHASH_COMMIT)" go cmd/ushort/main.go
