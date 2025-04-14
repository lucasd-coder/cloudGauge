# CloudGauge

[![codecov](https://codecov.io/gh/lucasd-coder/cloudGauge/graph/badge.svg?token=DD8Z1G8C4K)](https://codecov.io/gh/lucasd-coder/cloudGauge)
[![Documentation](https://img.shields.io/badge/Documentation-TODO-blue)](https://github.com/lucasd-coder/cloudGauge/tree/master/docs) 

## Tecnologias

Este projeto utiliza as seguintes tecnologias:

* [Golang](https://go.dev/) v1.24.1
* [PostgreSQL](https://www.postgresql.org/)
* [Kafka](https://kafka.apache.org/)
* [Google Wire](https://github.com/google/wire)

## Diagramas do Projeto

Os diagramas do projeto podem ser encontrados na pasta `docs`.

## Documentação

A documentação detalhada do projeto está disponível. [docs](https://github.com/lucasd-coder/cloudGauge/tree/master/docs)

## Como Rodar o Projeto

### Pré-requisitos

Certifique-se de ter o [Docker](https://www.docker.com/) e o [Make](https://www.gnu.org/software/make/) instalados em sua máquina.

### Subindo a Aplicação

Utilize o `Makefile` para facilitar a execução dos comandos Docker:

* **Para subir a aplicação e reiniciar os containers (caso já existam):**
    ```bash
    make docker_restart
    ```

* **Para subir a aplicação em containers Docker (para desenvolvimento):**
    ```bash
    make docker_dev
    ```

* **Para executar a aplicação diretamente (sem Docker):**
    ```bash
    make run_application
    ```

### Executando os Testes

Para rodar todos os testes da aplicação, execute o seguinte comando:

```bash
make test
```
### 🧪 Gerar Mocks

Para adicionar novos mocks, inclua o caminho das interfaces na variável MOCK_SOURCES dentro do Makefile:

```bash
make generate_mocks
```

### 🔗 Gerar Injeção de Dependências

```bash
make wire-gen
```

### 📎 Observações

* **Certifique-se de que o Docker esteja instalado e rodando para utilizar os comandos de ambiente.**

### 📁 Estrutura de Pastas

```bash

├── docs
│  ├── arquitetura.drawio
│  ├── arquitetura.png
│  ├── diagrama_classe.drawio
│  ├── diagrama_classe.png
│  ├── diagrama_simples.drawio
│  └── diagrama_simples.png
├── LICENSE
├── pulseReceiver
│  ├── cmd
│  │  └── app
│  │     └── main.go
│  ├── config
│  │  ├── config-dev.yml
│  │  ├── config.go
│  │  └── config.yml
│  ├── coverage.out
│  ├── docker-compose.yml
│  ├── go.mod
│  ├── go.sum
│  ├── infra
│  │  ├── docker
│  │  │  └── Dockerfile
│  │  ├── http
│  │  │  ├── health.http
│  │  │  └── usage_aggregation.http
│  │  ├── kafka
│  │  │  └── create-topics.sh
│  │  └── postgres
│  │     └── initdb.sh
│  ├── internal
│  │  ├── app
│  │  │  └── app.go
│  │  ├── controller
│  │  │  ├── controller.go
│  │  │  ├── controller_test.go
│  │  │  └── usageaggregation.go
│  │  ├── domain
│  │  │  ├── usageaggregation
│  │  │  │  ├── interfaces.go
│  │  │  │  ├── repository
│  │  │  │  │  ├── repository.go
│  │  │  │  │  └── repository_test.go
│  │  │  │  ├── service
│  │  │  │  │  ├── service.go
│  │  │  │  │  ├── usageaggregation_service.go
│  │  │  │  │  └── usageaggregation_test.go
│  │  │  │  ├── usageaggregation.go
│  │  │  │  └── usageaggregation_test.go
│  │  │  └── usageaggregationhistory
│  │  │     ├── repository
│  │  │     │  ├── repository.go
│  │  │     │  └── repository_test.go
│  │  │     ├── usageaggregationhistory.go
│  │  │     └── usegeaggregationhistory_interfaces.go
│  │  ├── inject
│  │  │  ├── wire.go
│  │  │  └── wire_gen.go
│  │  ├── mocks
│  │  │  ├── interfaces_mock.go
│  │  │  ├── shared_mock.go
│  │  │  └── usegeaggregationhistory_interfaces_mock.go
│  │  ├── processor
│  │  │  ├── processor.go
│  │  │  └── processor_usageaggregation.go
│  │  ├── provider
│  │  │  ├── jobscheduled
│  │  │  │  └── jobscheduler.go
│  │  │  ├── kafka
│  │  │  │  ├── publisher.go
│  │  │  │  └── subscription.go
│  │  │  ├── logger
│  │  │  │  ├── interface.go
│  │  │  │  └── slog.go
│  │  │  ├── middleware
│  │  │  │  └── logger.go
│  │  │  ├── migrations
│  │  │  │  └── migrations.go
│  │  │  ├── postgres
│  │  │  │  └── postgres.go
│  │  │  └── validator
│  │  │     ├── validator.go
│  │  │     └── validator_test.go
│  │  ├── scheduledtaskrunner
│  │  │  ├── scheduledtaskrunner.go
│  │  │  └── startdispatchscheduler.go
│  │  ├── server
│  │  │  └── server.go
│  │  ├── shared
│  │  │  ├── errors
│  │  │  │  └── errors.go
│  │  │  └── shared.go
│  │  └── subscription
│  │     ├── pulseReceiver.go
│  │     └── subscription.go
│  └── Makefile
└── README.md
```
