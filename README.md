# movies

Aplicativo web de listagem de vídeos em Go e Vannila Javascript para aplicação de conceitos base da programação web e arquitetura de software.

## Funcionalidades

- Listagem de filmes;
- Detalhes de um filme;
- Busca de filmes;
- Cadastro / login de usuários;
- Usuário pode adicionar filmes nos favoritos / watchlist.

## Tecnologias

- Go 1.20+ (HTTP server/backend);
- Vanilla Javascript (client/frontend);
- PostgreSQL;
- Docker.

### Client-side

- app.js - Conexão com a UI
- services/API.js - Requisições HTTP para o backend
- components/\* - Web components

### Server-side

- main.go - Configuração do servidor
- providers/\*.go - Repositórios dos Models
- handlers/\*.go - API handlers
- models/\*.go - Models da aplicação
- logger/logger.go - Utilitario para logs da aplicação

## Como rodar

### Pré-requisitos

- Docker
- Docker Compose

### Passo a passo

Antes de iniciar, certifique-se de que o Docker e o Docker Compose estão instalados no seu sistema.

Crie um novo arquivo `.env` na raiz do projeto e copie o conteúdo do arquivo `.env.example` para ele.

1.  **Subir a aplicação**:
    Execute o comando abaixo na raiz do projeto para criar e iniciar os containers (aplicação e banco de dados).

    ```bash
    docker-compose up -d --build
    ```

    A aplicação ficará disponível em `http://localhost:8080`.

2.  **Inicializar o banco de dados**:
    Na primeira execução, é necessário rodar o script de instalação para criar as tabelas e popular o banco de dados.

    ```bash
    docker exec -it movies-app-1 go run database/import/install.go
    ```

3.  **Desenvolvimento (Live Reload)**:
    O ambiente está configurado com `air` para _live reload_. Qualquer alteração salva nos arquivos `.go` reiniciará a aplicação automaticamente dentro do container.

4.  **Parar a aplicação**:
    Para parar e remover os containers:

    ```bash
    docker-compose down
    ```

## Documentação

- Diagrama de Entidade-Relacionamento: [ENTITY_RELATION_DIAGRAM.MD](docs/ENTITY_RELATION_DIAGRAM.MD)
- Diagrama de Arquitetura: [PROJECT_ARCHITECTURE.MD](docs/PROJECT_ARCHITECTURE.MD)
- Guia de Performance Frontend: [FRONTEND_PERFORMANCE_GUIDE.md](docs/FRONTEND_PERFORMANCE_GUIDE.md)

## Models

- Movie
- Genre
- Actor
- User
