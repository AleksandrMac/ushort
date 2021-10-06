GIT_COMMIT=$(shell git rev-list -1 HEAD)
GITHASH_COMMIT=$(shell git log --format="%h" -n 1)

.PHONY: test check
test:
	go test -race -coverprofile=coverage_controller.out -timeout 30s github.com/AleksandrMac/ushort/pkg/controller
	go test -race -coverprofile=coverage_utils.out -timeout 30s github.com/AleksandrMac/ushort/pkg/utils

check:
	golangci-lint run
build:
	go build -ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.GitHashCommit=$(GITHASH_COMMIT)" cmd/ushort/main.go
