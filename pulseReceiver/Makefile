GOPATH ?= go env GOPATH

# Diretório onde os mocks serão gerados
MOCK_DIR := ./internal/mocks

# Lista de arquivos-fonte de onde as interfaces são extraídas
MOCK_SOURCES := \
	./internal/domain/usageaggregation/interfaces.go \
	./internal/shared/shared.go \
	./internal/domain/usageaggregationhistory/usegeaggregationhistory_interfaces.go
	

DOCKER_COMPOSE = docker compose

.PHONY: generate_mocks run_application test coverage go_tidy

generate_mocks:
	@command -v mockgen >/dev/null 2>&1 || { \
		echo "mockgen não encontrado. Instalando..."; \
		go install go.uber.org/mock/mockgen@latest; \
	}
	@echo "Gerando mocks..."
	@for src in $(MOCK_SOURCES); do \
		base=$$(basename $$src .go); \
		dir=$$(dirname $$src); \
		pkg=$$(basename $(MOCK_DIR)); \
		dest=$(MOCK_DIR)/$${base}_mock.go; \
		mockgen -source=$$src -destination=$$dest -package=$$pkg || exit 1; \
		echo "✓ Mock gerado: $$dest"; \
	done

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
	sleep 25
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

docker_down:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

docker_restart:
	$(MAKE) docker_down
	$(MAKE) docker_dev