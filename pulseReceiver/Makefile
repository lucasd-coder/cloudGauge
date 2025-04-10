GOPATH ?= go env GOPATH

.PHONY: generate-mocks run-application test coverage go-tidy

generate-mocks:
	./scripts/mockery.sh

run-application:
	GO_PROFILE=dev GO111MODULE=on go run ./cmd/app/main.go

test:
	GO111MODULE=on go test -coverprofile coverage.out ./...

coverage:
	GO111MODULE=on go test -coverprofile coverage.out ./...
	GO111MODULE=on go tool cover -html=coverage.out

go-tidy:
	go mod tidy -go=1.24.1
