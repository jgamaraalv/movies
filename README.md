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
├── .github/workflows/  # CI/CD com GitHub Actions
│   └── ci-cd.yaml     # Pipeline de CI/CD
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
- **GitHub Actions** - CI/CD automatizado
- **GitHub Container Registry** - Armazenamento de imagens Docker

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
# Banco de Dados
POSTGRES_USER=seu_usuario
POSTGRES_PASSWORD=sua_senha_segura
POSTGRES_DB=movies_db

# Aplicação
JWT_SECRET=seu_secret_jwt_muito_seguro_aqui

# Opcional (produção)
DOCKER_REGISTRY=ghcr.io/seu-usuario
VERSION=latest
APP_PORT=8080
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

---

## Produção

### Arquitetura Docker de Produção

O projeto utiliza um **Dockerfile multi-stage** otimizado para produção:

```
┌─────────────────────────────────────────────────────────────┐
│                    MULTI-STAGE BUILD                        │
├─────────────────────────────────────────────────────────────┤
│  Stage 1: dev             │ Ambiente de desenvolvimento     │
│  Stage 2: frontend-builder│ Build do frontend (Vite)        │
│  Stage 3: backend-builder │ Compilação do Go                │
│  Stage 4: prod            │ Imagem final (~20MB)            │
└─────────────────────────────────────────────────────────────┘
```

### Características de Segurança

| Recurso                 | Descrição                                    |
| ----------------------- | -------------------------------------------- |
| 🔒 Usuário não-root     | Container executa como `appuser` (UID 10001) |
| 📁 Filesystem read-only | Sistema de arquivos em modo somente leitura  |
| 🚫 no-new-privileges    | Impede escalação de privilégios              |
| 🗑️ CAP_DROP ALL         | Remove todas as capabilities Linux           |
| 🌐 Rede isolada         | Serviços em rede interna sem acesso externo  |
| 📊 Resource limits      | Limites de CPU e memória por container       |
| 🩺 Health checks        | Verificação contínua de saúde dos serviços   |
| 📝 Logging estruturado  | Logs com rotação automática                  |

### Deploy Manual com Docker Compose

```bash
# Build e inicialização dos containers de produção
docker compose -f docker-compose.prod.yaml up -d --build

# Verificar status dos containers
docker compose -f docker-compose.prod.yaml ps

# Ver logs em tempo real
docker compose -f docker-compose.prod.yaml logs -f

# Parar serviços
docker compose -f docker-compose.prod.yaml down
```

### Inicialização do Banco de Dados

> ** Automático em Produção**: O banco de dados é inicializado automaticamente na primeira execução!

O `docker-compose.prod.yaml` monta o arquivo `database-dump.sql` no diretório `/docker-entrypoint-initdb.d/` do PostgreSQL. Isso faz com que o script SQL seja executado **automaticamente** quando o volume do banco é criado pela primeira vez.

```yaml
# Configuração no docker-compose.prod.yaml
volumes:
  - ./server/database/import/database-dump.sql:/docker-entrypoint-initdb.d/01-init.sql:ro
```

**Comportamento:**

- **Primeiro deploy**: O banco é criado e populado automaticamente com ~4.800 filmes
- **Deploys subsequentes**: O volume persiste e os dados são mantidos
- **Reset do banco**: Use `docker compose -f docker-compose.prod.yaml down -v` para remover o volume e reinicializar

**Verificar se o banco foi inicializado:**

```bash
# Verificar se as tabelas existem
docker exec movies-postgres psql -U $POSTGRES_USER -d $POSTGRES_DB -c "\dt"

# Contar registros
docker exec movies-postgres psql -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT COUNT(*) FROM movies;"
```

### Variáveis de Ambiente para Produção

Crie um arquivo `.env` com as seguintes variáveis:

```env
# === OBRIGATÓRIAS ===
POSTGRES_USER=movies_prod
POSTGRES_PASSWORD=<senha-forte-aqui>
POSTGRES_DB=movies_production
JWT_SECRET=<secret-jwt-forte-de-256-bits>

# === OPCIONAIS ===
# Registry Docker (para CI/CD)
DOCKER_REGISTRY=ghcr.io/seu-usuario

# Versão da imagem (SHA do commit ou tag semântica)
VERSION=latest

# Porta da aplicação (padrão: 8080)
APP_PORT=8080
```

### Verificar Saúde dos Containers

```bash
# Verificar health check da aplicação
curl http://localhost:8080/health

# Resposta esperada:
# {"status":"healthy"}

# Verificar health de todos os containers
docker compose -f docker-compose.prod.yaml ps
```

---

## CI/CD com GitHub Actions

O projeto inclui um pipeline completo de CI/CD configurado em `.github/workflows/ci-cd.yaml`.

### Pipeline Overview

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Test      │───▶│   Build     │───▶│   Scan      │───▶│   Deploy    │
│  Backend    │    │   Docker    │    │   Trivy     │    │  (manual)   │
│  Frontend   │    │   Image     │    │   SARIF     │    │             │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### Funcionalidades do Pipeline

| Etapa             | Descrição                                               |
| ----------------- | ------------------------------------------------------- |
| **Test Backend**  | Testes Go, linting, verificação de formatação           |
| **Test Frontend** | Build de verificação do Vite                            |
| **Build & Push**  | Build multi-stage e push para GitHub Container Registry |
| **Security Scan** | Scan de vulnerabilidades com Trivy                      |
| **Deploy**        | Deploy automático (configurável)                        |

### Triggers do Pipeline

- **Push para `main`**: Executa pipeline completo com deploy
- **Pull Request**: Executa apenas testes (sem deploy)
- **Manual**: Permite execução via GitHub UI

### Configurar GitHub Secrets

Para o pipeline funcionar, configure os seguintes secrets no GitHub:

| Secret           | Descrição                                 |
| ---------------- | ----------------------------------------- |
| `GITHUB_TOKEN`   | Automático (não precisa configurar)       |
| `SERVER_HOST`    | IP/hostname do servidor (para deploy SSH) |
| `SERVER_USER`    | Usuário SSH do servidor                   |
| `SERVER_SSH_KEY` | Chave SSH privada para deploy             |

### Deploy Automático

O pipeline está configurado com múltiplas opções de deploy:

#### Opção 1: Deploy via SSH para VPS/VM

Descomente a seção no workflow e configure os secrets:

```yaml
- name: 🚀 Deploy to server
  uses: appleboy/ssh-action@v1.0.3
  with:
    host: ${{ secrets.SERVER_HOST }}
    username: ${{ secrets.SERVER_USER }}
    key: ${{ secrets.SERVER_SSH_KEY }}
    script: |
      cd /opt/movies
      docker compose -f docker-compose.prod.yaml pull
      docker compose -f docker-compose.prod.yaml up -d
```

#### Opção 2: Deploy para Plataformas PaaS

O pipeline pode ser adaptado para:

- **Fly.io**: `flyctl deploy`
- **Railway**: API de deploy
- **Render**: Webhook de deploy
- **DigitalOcean App Platform**: API de deploy

### Executar Pipeline Manualmente

1. Vá para **Actions** no repositório GitHub
2. Selecione **CI/CD Pipeline**
3. Clique em **Run workflow**
4. Escolha o ambiente de deploy

---

## Scripts NPM

Na pasta `web/`:

| Script            | Descrição                                  |
| ----------------- | ------------------------------------------ |
| `npm run dev`     | Build em watch mode para desenvolvimento   |
| `npm run build`   | Build de produção (otimizado e minificado) |
| `npm run preview` | Preview do build de produção localmente    |

## API Endpoints

### Health Check

- `GET /health` - Verifica saúde da aplicação e conexão com banco

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

## Comandos Úteis

### Desenvolvimento

```bash
# Subir ambiente de desenvolvimento
docker compose up -d --build

# Ver logs em tempo real
docker compose logs -f app

# Executar comando dentro do container
docker exec -it movies-app-1 sh

# Rebuild apenas o backend
docker compose up -d --build app
```

### Produção

```bash
# Build de produção
docker compose -f docker-compose.prod.yaml build

# Deploy com nova versão
VERSION=v1.0.0 docker compose -f docker-compose.prod.yaml up -d

# Verificar recursos dos containers
docker stats

# Backup do banco de dados
docker exec movies-postgres pg_dump -U $POSTGRES_USER $POSTGRES_DB > backup.sql
```

### Manutenção

```bash
# Limpar imagens não utilizadas
docker image prune -a

# Limpar volumes órfãos
docker volume prune

# Ver uso de disco
docker system df

# Logs do sistema
docker compose -f docker-compose.prod.yaml logs --tail=100
```

## Parar a Aplicação

Para parar e remover os containers:

```bash
# Desenvolvimento
docker compose down

# Produção
docker compose -f docker-compose.prod.yaml down
```

Para remover também os volumes (dados do banco):

```bash
docker compose down -v
```

## Estrutura de Dados

### Principais Entidades

- **Movie** - Informações dos filmes (título, sinopse, elenco, etc.)
- **User** - Usuários do sistema
- **Actor** - Atores/atrizes
- **Genre** - Gêneros cinematográficos
- **UserMovie** - Relação entre usuários e filmes (favoritos/watchlist)

## Segurança

### Aplicação

- Senhas são hasheadas com bcrypt
- Autenticação via JWT (JSON Web Tokens)
- Validação de dados no backend (Value Objects)
- Sanitização de inputs

### Containers (Produção)

- Usuário não-root em todos os containers
- Filesystem read-only
- Capabilities Linux removidas
- Limites de recursos (CPU/memória)
- Rede isolada entre serviços
- Health checks ativos
- Logging com rotação automática

## Licença

MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.
