# Movies - Aplicativo Web de Listagem de Filmes

Aplicativo web full-stack para listagem e gerenciamento de filmes, desenvolvido com **Go** e **Vanilla JavaScript**. Projeto com foco na aplicação de conceitos fundamentais de programação web e arquitetura de software, seguindo os princípios de **Clean Architecture** e **Domain-Driven Design (DDD)**.

## Funcionalidades

- **Listagem de Filmes**

  - Top 10 filmes mais populares
  - Filmes aleatórios para descoberta
  - Busca avançada com filtros por gênero e ordenação
  - Detalhes completos de cada filme (sinopse, elenco, trailer)

- **Sistema de Autenticação**

  - Cadastro de novos usuários
  - Login seguro com JWT
  - Gerenciamento de conta

- **Coleções Pessoais**
  - Adicionar filmes aos favoritos
  - Criar lista de desejos (watchlist)
  - Visualizar coleções pessoais

## Arquitetura

O projeto segue os princípios de **Clean Architecture** e **DDD**, organizando o código em camadas bem definidas:

- **Domain Layer**: Entidades, Value Objects e interfaces de repositório
- **Application Layer**: Casos de uso (use cases) que orquestram a lógica de negócio
- **Interface Layer**: Handlers HTTP que processam requisições
- **Infrastructure Layer**: Implementações concretas (PostgreSQL, logger, JWT)

### Estrutura do Projeto

```
movies/
├── server/              # Backend Go (Clean Architecture)
│   ├── cmd/api/        # Entry point da aplicação
│   ├── internal/        # Código interno
│   │   ├── domain/     # Camada de domínio
│   │   ├── usecase/    # Casos de uso
│   │   ├── handler/    # Handlers HTTP
│   │   └── infrastructure/  # Implementações
│   ├── models/         # DTOs
│   ├── pkg/           # Pacotes reutilizáveis
│   └── database/      # Scripts de banco de dados
│
├── web/                # Frontend (código fonte)
│   ├── src/
│   │   ├── components/  # Web Components
│   │   ├── services/    # Serviços (API, Router, Store)
│   │   ├── app.js       # Entry point
│   │   └── styles.css   # Estilos
│   ├── index.html
│   └── package.json
│
└── public/             # Build/dist (gerado automaticamente)
```

Para mais detalhes sobre a arquitetura, consulte a [documentação completa](docs/PROJECT_ARCHITECTURE.MD).

## Tecnologias

### Backend

- **Go 1.24+** - Linguagem de programação
- **PostgreSQL** - Banco de dados relacional
- **JWT** - Autenticação e autorização
- **Air** - Hot reload em desenvolvimento

### Frontend

- **Vanilla JavaScript** - Sem frameworks, JavaScript puro
- **ES Modules** - Módulos ES6 nativos
- **Web Components** - Componentes reutilizáveis
- **Vite 5.4+** - Build tool e otimizações

### DevOps

- **Docker** - Containerização
- **Docker Compose** - Orquestração de containers

## Pré-requisitos

### Para Desenvolvimento Local

- **Docker** 20.10+
- **Docker Compose** 2.0+
- **Node.js** 20+ e **npm** (opcional, para desenvolvimento do frontend localmente)

### Para Produção

- **Docker** 20.10+
- **Docker Compose** 2.0+

## Instalação e Configuração

### 1. Clone o repositório

```bash
git clone <repository-url>
cd movies
```

### 2. Configure as variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto baseado no `.env.example`:

```bash
cp .env.example .env
```

Edite o `.env` com suas configurações:

```env
POSTGRES_USER=seu_usuario
POSTGRES_PASSWORD=sua_senha
POSTGRES_DB=movies_db
DATABASE_URL=postgres://seu_usuario:sua_senha@postgres:5432/movies_db?sslmode=disable
JWT_SECRET=seu_secret_jwt_aqui
```

## Desenvolvimento

### Opção 1: Tudo no Docker (Recomendado)

Esta é a forma mais simples e recomendada para desenvolvimento:

```bash
# Subir todos os serviços
docker-compose up -d --build
```

Isso iniciará:

- **Backend Go** na porta `8080` com hot reload (Air)
- **Frontend Vite** em watch mode, buildando automaticamente para `public/`
- **PostgreSQL** na porta `5432`

A aplicação estará disponível em `http://localhost:8080`.

#### Inicializar o banco de dados

Na primeira execução, é necessário popular o banco de dados:

```bash
docker exec movies-app-1 go run ./database/import/install.go
```

#### Desenvolvimento

- **Backend**: Alterações em arquivos `.go` reiniciam automaticamente (Air)
- **Frontend**: Alterações em `web/` são buildadas automaticamente para `public/` (Vite watch mode)

### Opção 2: Desenvolvimento Híbrido

Para desenvolvimento do frontend localmente (sem Docker):

#### 1. Instalar dependências do frontend

```bash
cd web
npm install
```

#### 2. Rodar build em watch mode

```bash
npm run dev
```

Isso observará mudanças em `web/` e buildará automaticamente para `public/`.

#### 3. Subir apenas backend e banco via Docker

```bash
docker-compose up postgres app -d
```

#### 4. Inicializar o banco de dados

```bash
docker exec movies-app-1 go run ./database/import/install.go
```

## Produção

### Build de Produção

Para gerar os arquivos otimizados do frontend:

```bash
cd web
npm install
npm run build
```

Isso gerará arquivos minificados e otimizados em `public/`.

### Deploy com Docker

O projeto inclui um `Dockerfile` multi-stage otimizado para produção:

```bash
docker-compose -f docker-compose.prod.yaml up -d --build
```

## Scripts NPM

Na pasta `web/`:

| Script            | Descrição                                  |
| ----------------- | ------------------------------------------ |
| `npm run dev`     | Build em watch mode para desenvolvimento   |
| `npm run build`   | Build de produção (otimizado e minificado) |
| `npm run preview` | Preview do build de produção localmente    |

## API Endpoints

### Autenticação

- `POST /api/account/register/` - Registrar novo usuário
- `POST /api/account/authenticate/` - Autenticar usuário (login)

### Filmes

- `GET /api/movies/top` - Listar top 10 filmes mais populares
- `GET /api/movies/random` - Listar filmes aleatórios
- `GET /api/movies/search?q={query}&order={order}&genre={genre}` - Buscar filmes
- `GET /api/movies/{id}` - Obter detalhes de um filme
- `GET /api/genres` - Listar todos os gêneros

### Coleções (Requer autenticação)

- `GET /api/account/favorites/` - Listar filmes favoritos
- `GET /api/account/watchlist/` - Listar watchlist
- `POST /api/account/save-to-collection/` - Adicionar filme à coleção

**Autenticação**: Endpoints protegidos requerem header `Authorization: Bearer {token}`

## Testes

_Seção para testes quando implementados_

## Documentação Adicional

- [Arquitetura do Projeto](docs/PROJECT_ARCHITECTURE.MD) - Detalhes sobre Clean Architecture e DDD
- [Diagrama de Entidade-Relacionamento](docs/ENTITY_RELATION_DIAGRAM.MD) - Estrutura do banco de dados
- [Guia de Performance Frontend](docs/FRONTEND_PERFORMANCE_GUIDE.md) - Otimizações e boas práticas

## Parar a Aplicação

Para parar e remover os containers:

```bash
docker-compose down
```

Para remover também os volumes (dados do banco):

```bash
docker-compose down -v
```

## Estrutura de Dados

### Principais Entidades

- **Movie** - Informações dos filmes (título, sinopse, elenco, etc.)
- **User** - Usuários do sistema
- **Actor** - Atores/atrizes
- **Genre** - Gêneros cinematográficos
- **UserMovie** - Relação entre usuários e filmes (favoritos/watchlist)

## Segurança

- Senhas são hasheadas com bcrypt
- Autenticação via JWT (JSON Web Tokens)
- Validação de dados no backend (Value Objects)
- Sanitização de inputs
