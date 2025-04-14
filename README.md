# CloudGauge

[![codecov](https://codecov.io/gh/lucasd-coder/cloudGauge/graph/badge.svg?token=9MYI8R7IYZ)](https://codecov.io/gh/lucasd-coder/cloudGauge)
[![Documentation](https://img.shields.io/badge/Documentation-TODO-blue)](https://github.com/lucasd-coder/cloudGauge/docs) ## Tecnologias

Este projeto utiliza as seguintes tecnologias:

* [Golang](https://go.dev/) v1.24.1
* [PostgreSQL](https://www.postgresql.org/)
* [Kafka](https://kafka.apache.org/)
* [Google Wire](https://github.com/google/wire)

## Diagramas do Projeto

Os diagramas do projeto podem ser encontrados na pasta `docs`.

## Documentação

A documentação detalhada do projeto está disponível [**TODO: Adicionar o link da documentação aqui**].

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

