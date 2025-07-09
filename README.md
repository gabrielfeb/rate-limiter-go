Markdown

# Rate Limiter em Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Powered-2496ED.svg)](https://www.docker.com)

## Visão Geral

Este projeto consiste na implementação de um Rate Limiter (Limitador de Requisições) em Go, projetado para atuar como um middleware em serviços web. O sistema controla o tráfego de requisições, aplicando limites com base no endereço IP do cliente ou em um token de acesso fornecido, utilizando um banco de dados Redis para persistência de estado.

A arquitetura foi desenvolvida de forma desacoplada, permitindo a fácil substituição da camada de armazenamento (Redis) por outra tecnologia de persistência, graças ao uso do padrão de projeto *Strategy*.

## Principais Funcionalidades

-   **Limitação por Endereço IP**: Restringe o número de requisições de um único IP por segundo.
-   **Limitação por Token de Acesso**: Aplica limites específicos para tokens fornecidos no cabeçalho `API_KEY`, com precedência sobre os limites de IP.
-   **Armazenamento em Redis**: Utiliza o Redis para armazenar e consultar os dados de limitação de forma eficiente e escalável.
-   **Configuração Flexível**: Todos os parâmetros, como limites e durações de bloqueio, são configurados via variáveis de ambiente.
-   **Arquitetura Extensível**: A lógica do limitador é separada do middleware e da camada de armazenamento, aderindo ao Princípio da Inversão de Dependência (DIP).
-   **Ambiente Containerizado**: O projeto é totalmente containerizado com Docker e orquestrado com Docker Compose para facilitar a execução e o teste.

## Arquitetura

A solução é dividida em três camadas principais:

1.  **Middleware (`/middleware`)**: O ponto de entrada que intercepta as requisições HTTP. É responsável por extrair o identificador (IP ou token) e consultar a lógica do limitador.
2.  **Lógica do Limiter (`/limiter`)**: O núcleo do sistema. Recebe um identificador e decide se a requisição deve ser permitida ou bloqueada, sem conhecer detalhes de HTTP ou Redis.
3.  **Camada de Armazenamento (`/limiter/storage.go`)**: Uma abstração definida por uma interface (`Storage`). A implementação concreta (`RedisStorage`) contém a lógica específica para interagir com o Redis. Essa abstração permite trocar o Redis por outro banco de dados sem alterar o resto da aplicação.

## Pré-requisitos

Para executar este projeto, você precisará ter as seguintes ferramentas instaladas:

-   Go (versão 1.21 ou superior)
-   Docker
-   Docker Compose
-   Um compilador C/C++ (como o MinGW-w64 no Windows) para executar os testes com o *race detector*.

## Configuração

A aplicação é configurada através de um arquivo `.env` na raiz do projeto. Crie este arquivo a partir do exemplo abaixo:

**.env.example**
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
A tabela a seguir descreve cada variável de ambiente:

| Variável | Descrição | Valor Padrão |
| :--- | :--- | :--- |
| `REDIS_ADDR` | Endereço do servidor Redis (usado pelo contêiner Go para se conectar). | `redis:6379` |
| `REDIS_PASSWORD` | Senha para autenticação no Redis. | (vazio) |
| `REDIS_DB` | Número do banco de dados Redis a ser utilizado. | `0` |
| `IP_LIMIT_PER_SECOND` | Número máximo de requisições por segundo para um único IP. | `5` |
| `IP_BLOCK_DURATION` | Duração do bloqueio para um IP que excedeu o limite. | `1m` (1 minuto) |
| `TOKEN_LIMIT_PER_SECOND`| Número máximo de requisições por segundo para um token de acesso. | `20` |
| `TOKEN_BLOCK_DURATION` | Duração do bloqueio para um token que excedeu o limite. | `5m` (5 minutos) |
| `DEFAULT_TEST_TOKEN` | Um token de exemplo para ser usado nos testes manuais. | `my-secret-token-123`|

### Como Executar

```markdown
## Como Executar

### Usando Docker (Método Recomendado)

Este é o método mais simples e garante que todo o ambiente (aplicação + Redis) funcione corretamente.

1.  **Clone o Repositório**
    ```sh
    git clone <url-do-repositorio>
    cd rate-limiter-go
    ```

2.  **Crie o arquivo `.env`**
    Copie o conteúdo do `.env.example` (mostrado acima) para um novo arquivo chamado `.env`.

3.  **Inicie os Serviços**
    Execute o seguinte comando para construir a imagem e iniciar os contêineres em segundo plano:
    ```sh
    docker-compose up --build -d
    ```

4.  **Verifique se os Contêineres estão Rodando**
    ```sh
    docker-compose ps
    ```
    Você deve ver os serviços `app` e `redis` com o status `Up` ou `running`.

5.  **Verifique os Logs da Aplicação**
    Para confirmar a conexão com o Redis e o início do servidor, verifique os logs:
    ```sh
    docker-compose logs -f app
    ```
    Você deve ver a mensagem "Starting server on port 8080...".

### Encerrando o Ambiente

Para parar e remover todos os contêineres e redes criadas, use:
```sh
docker-compose down

Com certeza. Aqui estão todas as outras seções do arquivo README.md, prontas para você copiar, colar e montar seu documento completo.

Título e Badges
Markdown

# Rate Limiter em Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Powered-2496ED.svg)](https://www.docker.com)
Visão Geral
Markdown

## Visão Geral

Este projeto consiste na implementação de um Rate Limiter (Limitador de Requisições) em Go, projetado para atuar como um middleware em serviços web. O sistema controla o tráfego de requisições, aplicando limites com base no endereço IP do cliente ou em um token de acesso fornecido, utilizando um banco de dados Redis para persistência de estado.

A arquitetura foi desenvolvida de forma desacoplada, permitindo a fácil substituição da camada de armazenamento (Redis) por outra tecnologia de persistência, graças ao uso do padrão de projeto *Strategy*.
Principais Funcionalidades
Markdown

## Principais Funcionalidades

-   **Limitação por Endereço IP**: Restringe o número de requisições de um único IP por segundo.
-   **Limitação por Token de Acesso**: Aplica limites específicos para tokens fornecidos no cabeçalho `API_KEY`, com precedência sobre os limites de IP.
-   **Armazenamento em Redis**: Utiliza o Redis para armazenar e consultar os dados de limitação de forma eficiente e escalável.
-   **Configuração Flexível**: Todos os parâmetros, como limites e durações de bloqueio, são configurados via variáveis de ambiente.
-   **Arquitetura Extensível**: A lógica do limitador é separada do middleware e da camada de armazenamento, aderindo ao Princípio da Inversão de Dependência (DIP).
-   **Ambiente Containerizado**: O projeto é totalmente containerizado com Docker e orquestrado com Docker Compose para facilitar a execução e o teste.
Arquitetura
Markdown

## Arquitetura

A solução é dividida em três camadas principais:

1.  **Middleware (`/middleware`)**: O ponto de entrada que intercepta as requisições HTTP. É responsável por extrair o identificador (IP ou token) e consultar a lógica do limitador.
2.  **Lógica do Limiter (`/limiter`)**: O núcleo do sistema. Recebe um identificador e decide se a requisição deve ser permitida ou bloqueada, sem conhecer detalhes de HTTP ou Redis.
3.  **Camada de Armazenamento (`/limiter/storage.go`)**: Uma abstração definida por uma interface (`Storage`). A implementação concreta (`RedisStorage`) contém a lógica específica para interagir com o Redis. Essa abstração permite trocar o Redis por outro banco de dados sem alterar o resto da aplicação.
Pré-requisitos
Markdown

## Pré-requisitos

Para executar este projeto, você precisará ter as seguintes ferramentas instaladas:

-   Go (versão 1.21 ou superior)
-   Docker
-   Docker Compose
-   Um compilador C/C++ (como o MinGW-w64 no Windows) para executar os testes com o *race detector*.
Configuração (Exemplo de .env)
Markdown

## Configuração

A aplicação é configurada através de um arquivo `.env` na raiz do projeto. Crie este arquivo a partir do exemplo abaixo:

**.env.example**
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
(A Tabela de Configuração foi fornecida na resposta anterior).


---
### Como Executar

```markdown
## Como Executar

### Usando Docker (Método Recomendado)

Este é o método mais simples e garante que todo o ambiente (aplicação + Redis) funcione corretamente.

1.  **Clone o Repositório**
    ```sh
    git clone <url-do-repositorio>
    cd rate-limiter-go
    ```

2.  **Crie o arquivo `.env`**
    Copie o conteúdo do `.env.example` (mostrado acima) para um novo arquivo chamado `.env`.

3.  **Inicie os Serviços**
    Execute o seguinte comando para construir a imagem e iniciar os contêineres em segundo plano:
    ```sh
    docker-compose up --build -d
    ```

4.  **Verifique se os Contêineres estão Rodando**
    ```sh
    docker-compose ps
    ```
    Você deve ver os serviços `app` e `redis` com o status `Up` ou `running`.

5.  **Verifique os Logs da Aplicação**
    Para confirmar a conexão com o Redis e o início do servidor, verifique os logs:
    ```sh
    docker-compose logs -f app
    ```
    Você deve ver a mensagem "Starting server on port 8080...".

### Encerrando o Ambiente

Para parar e remover todos os contêineres e redes criadas, use:
```sh
docker-compose down

---
### Utilização e Teste Manual

```markdown
## Utilização e Teste Manual

Com a aplicação rodando, você pode usar uma ferramenta como o `curl` para testar o rate limiting. A aplicação estará disponível em `http://localhost:8080`.

**Requisição Padrão (limitada por IP):**
```sh
curl -i http://localhost:8080/


Requisição com Token de Acesso:
O token deve ser fornecido no cabeçalho API_KEY.

Bash

curl -i -H "API_KEY: my-secret-token-123" http://localhost:8080/
Resposta de Limite Excedido:
Quando um limite é ultrapassado, o servidor responderá com:

Código HTTP: 429 Too Many Requests

Corpo da Resposta: you have reached the maximum number of requests or actions allowed within a certain time frame

---
### Executando os Testes Automatizados

```markdown
## Executando os Testes Automatizados

Para rodar a suíte de testes de unidade e componente localmente:

1.  **Limpe o Cache de Testes (Opcional, mas recomendado)**
    ```sh
    go clean -testcache
    ```
2.  **Execute a Suíte de Testes Completa**
    Este comando roda todos os testes, calcula a cobertura de código e ativa o *race detector*.
    ```sh
    # No Windows (PowerShell)
    $env:CGO_ENABLED="1"; go test -v -cover -race ./...

    # No Windows (CMD)
    set CGO_ENABLED=1 && go test -v -cover -race ./...

    # No Linux / macOS
    CGO_ENABLED=1 go test -v -cover -race ./...
    ```

3.  **Gerar Relatório de Cobertura Visual**
    Para analisar a cobertura de forma detalhada em HTML:
    ```sh
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out
    ```

## Estrutura do Projeto

```plaintext
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

---
### Licença

```markdown
## Licença

Distribuído sob a licença MIT.