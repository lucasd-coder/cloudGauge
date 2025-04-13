GOPATH ?= go env GOPATH

DOCKER_COMPOSE = docker compose

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
	$(DOCKER_COMPOSE) up -d
	sleep 20
	$(DOCKER_COMPOSE) exec kafka bash /create-topics.sh

.PHONY: run_application

run_application:
	GO_PROFILE=dev GO111MODULE=on go run -race ./cmd/app/main.go

.PHONY: wire-gen

wire-gen:
	@command -v wire >/dev/null 2>&1 || { \
		echo "wire não encontrado. Instalando..."; \
		go install github.com/google/wire/cmd/wire@latest; \
	}
	wire ./internal/inject
