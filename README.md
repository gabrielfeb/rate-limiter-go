# Rate Limiter em Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Powered-2496ED.svg)](https://www.docker.com)

## Visão Geral

Este projeto consiste na implementação de um Rate Limiter (Limitador de Requisições) em Go, projetado para atuar como um middleware em serviços web. O sistema controla o tráfego de requisições, aplicando limites com base no endereço IP do cliente ou em um token de acesso fornecido, utilizando um banco de dados Redis para persistência de estado.

A arquitetura foi desenvolvida de forma desacoplada, permitindo a fácil substituição da camada de armazenamento (Redis) por outra tecnologia de persistência, graças ao uso do padrão de projeto *Strategy*.

## Principais Funcionalidades

- **Limitação por Endereço IP**: Restringe o número de requisições de um único IP por segundo.
- **Limitação por Token de Acesso**: Aplica limites específicos para tokens fornecidos no cabeçalho `API_KEY`, com precedência sobre os limites de IP.
- **Armazenamento em Redis**: Utiliza o Redis para armazenar e consultar os dados de limitação de forma eficiente e escalável.
- **Configuração Flexível**: Todos os parâmetros, como limites e durações de bloqueio, são configurados via variáveis de ambiente.
- **Arquitetura Extensível**: A lógica do limitador é separada do middleware e da camada de armazenamento, aderindo ao Princípio da Inversão de Dependência (DIP).
- **Ambiente Containerizado**: O projeto é totalmente containerizado com Docker e orquestrado com Docker Compose para facilitar a execução e o teste.

## Arquitetura

A solução é dividida em três camadas principais:

1. **Middleware (`/middleware`)**: Intercepta as requisições HTTP, extrai o identificador (IP ou token) e consulta o limitador.
2. **Lógica do Limiter (`/limiter`)**: Responsável por decidir se a requisição deve ser permitida ou bloqueada, desacoplado da lógica HTTP ou Redis.
3. **Camada de Armazenamento (`/limiter/storage.go`)**: Abstração definida por uma interface (`Storage`) com implementação concreta (`RedisStorage`), permitindo fácil substituição da tecnologia de persistência.

## Pré-requisitos

- Go (versão 1.21 ou superior)
- Docker
- Docker Compose
- Compilador C/C++ (ex: MinGW-w64 no Windows) para testes com *race detector*

## Configuração

Crie um arquivo `.env` na raiz do projeto com base no exemplo abaixo:

```ini
# Configurações do Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# Limites baseados em IP
IP_LIMIT_PER_SECOND=5
IP_BLOCK_DURATION=1m # m=minutos, s=segundos, h=horas

# Limites baseados em Token
TOKEN_LIMIT_PER_SECOND=20
TOKEN_BLOCK_DURATION=5m

# Token padrão para testes
DEFAULT_TEST_TOKEN=my-secret-token-123
```

### Descrição das Variáveis de Ambiente

| Variável               | Descrição                                                        | Valor Padrão             |
|------------------------|------------------------------------------------------------------|--------------------------|
| `REDIS_ADDR`           | Endereço do servidor Redis                                       | `redis:6379`             |
| `REDIS_PASSWORD`       | Senha para autenticação no Redis                                 | (vazio)                  |
| `REDIS_DB`             | Número do banco de dados Redis                                   | `0`                      |
| `IP_LIMIT_PER_SECOND`  | Nº máximo de requisições por segundo por IP                     | `5`                      |
| `IP_BLOCK_DURATION`    | Duração do bloqueio após ultrapassar limite por IP              | `1m`                     |
| `TOKEN_LIMIT_PER_SECOND`| Nº máximo de requisições por segundo por token                 | `20`                     |
| `TOKEN_BLOCK_DURATION` | Duração do bloqueio após ultrapassar limite por token           | `5m`                     |
| `DEFAULT_TEST_TOKEN`   | Token de exemplo para testes manuais                            | `my-secret-token-123`    |

## Como Executar

### Usando Docker (Método Recomendado)

Clone o repositório:

```bash
git clone <url-do-repositorio>
cd rate-limiter-go
```

Crie o arquivo `.env` com base no `.env.example`.

Inicie os serviços:

```bash
docker-compose up --build -d
```

Verifique se os contêineres estão ativos:

```bash
docker-compose ps
```

Verifique os logs da aplicação:

```bash
docker-compose logs -f app
```

Você deve ver a mensagem: `Starting server on port 8080...`

Para encerrar os serviços:

```bash
docker-compose down
```

## Utilização e Teste Manual

Com o servidor rodando em `http://localhost:8080`:

### Requisição Padrão (por IP):

```bash
curl -i http://localhost:8080/
```

### Requisição com Token (cabeçalho `API_KEY`):

```bash
curl -i -H "API_KEY: my-secret-token-123" http://localhost:8080/
```

### Quando o Limite é Excedido

- **HTTP Status:** `429 Too Many Requests`
- **Corpo da Resposta:**  
  ```
  you have reached the maximum number of requests or actions allowed within a certain time frame
  ```

## Executando os Testes Automatizados

### Limpar o cache de testes (opcional)

```bash
go clean -testcache
```

### Rodar todos os testes com cobertura e *race detector*

#### No Linux / macOS:

```bash
CGO_ENABLED=1 go test -v -cover -race ./...
```

#### No Windows (PowerShell):

```powershell
$env:CGO_ENABLED="1"; go test -v -cover -race ./...
```

#### No Windows (CMD):

```cmd
set CGO_ENABLED=1 && go test -v -cover -race ./...
```

### Gerar relatório de cobertura em HTML

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Estrutura do Projeto

```
/rate-limiter-go
|-- .env
|-- Dockerfile
|-- docker-compose.yml
|-- go.mod
|-- go.sum
|-- app/
|   `-- main.go
|-- config/
|   |-- config.go
|   `-- config_test.go
|-- limiter/
|   |-- limiter.go
|   |-- limiter_test.go
|   |-- storage.go
|   `-- storage_mock.go
|-- middleware/
|   |-- ratelimit.go
|   `-- ratelimit_test.go
`-- README.md
```

## Licença

Distribuído sob a licença [MIT](https://opensource.org/licenses/MIT).