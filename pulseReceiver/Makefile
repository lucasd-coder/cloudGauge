GOPATH ?= go env GOPATH

.PHONY: generate_mocks run_application test coverage go_tidy

generate_mocks:
	./scripts/mockery.sh

run-application:
	GO_PROFILE=dev GO111MODULE=on go run ./cmd/app/main.go

test:
	GO111MODULE=on go test -coverprofile coverage.out ./...

coverage:
	GO111MODULE=on go test -coverprofile coverage.out ./...
	GO111MODULE=on go tool cover -html=coverage.out

go_tidy:
	go mod tidy -go=1.24.1

.PHONY: docker_dev

docker_dev:
	@docker-compose -f docker-compose.yml up -d

.PHONY: run_application

run_application:
	GO_PROFILE=dev GO111MODULE=on go run -race ./cmd/app/main.go
