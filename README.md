# Movies – Movie Listing Web Application

A full-stack web application for listing and managing movies, built with **Go** and **Vanilla JavaScript**. The project focuses on applying fundamental web programming and software architecture concepts, following the principles of **Clean Architecture** and **Domain-Driven Design (DDD)**.

## Features

* **Movie Listing**

  * Top 10 most popular movies
  * Random movies for discovery
  * Advanced search with genre filters and sorting
  * Full movie details (synopsis, cast, trailer)

* **Authentication System**

  * User registration
  * Secure login with JWT
  * Account management

* **Personal Collections**

  * Add movies to favorites
  * Create a watchlist
  * View personal collections

* **AI-Powered Recommendation System**

  * Neural Collaborative Filtering model (NCF/NeuMF) trained with TensorFlow
  * 128-dimensional embeddings stored via pgvector in PostgreSQL
  * Real-time hybrid recommendation algorithm combining:
    * **Genre affinity** (35%) — prioritizes genres the user likes most
    * **Embedding similarity** (25%) — latent features from the NCF model
    * **Collaborative filtering** (20%) — movies liked by similar users
    * **Movie quality** (12% score + 8% popularity) — tiebreaker
  * Recommendations automatically recomputed on each user interaction
  * Cold-start support for new users via genre-based fallback

## Architecture

The project follows **Clean Architecture** and **DDD** principles, organizing the code into well-defined layers:
* **Domain Layer**: Entities, Value Objects, and repository interfaces
* **Application Layer**: Use cases that orchestrate business logic
* **Interface Layer**: HTTP handlers that process requests
* **Infrastructure Layer**: Concrete implementations (PostgreSQL, logger, JWT)

### Project Structure

```
movies/
├── server/              # Go backend (Clean Architecture)
│   ├── cmd/api/         # Application entry point
│   ├── internal/        # Internal code
│   │   ├── domain/      # Domain layer
│   │   ├── usecase/     # Use cases
│   │   ├── handler/     # HTTP handlers
│   │   └── infrastructure/  # Implementations
│   ├── models/          # DTOs
│   ├── pkg/             # Reusable packages
│   └── database/        # Database scripts
│
├── web/                 # Frontend (source code)
│   ├── src/
│   │   ├── components/  # Web Components
│   │   ├── services/    # Services (API, Router, Store)
│   │   ├── app.js       # Entry point
│   │   └── styles.css   # Styles
│   ├── index.html
│   └── package.json
│
├── recommender/         # ML recommendation system
│   ├── models/          # NCF/NeuMF model definition
│   ├── data/            # Training data and embeddings
│   ├── train.py         # Training pipeline
│   └── generate_embeddings.py  # Embedding extraction
│
├── .github/workflows/   # CI/CD with GitHub Actions
│   └── ci-cd.yaml       # CI/CD pipeline
│
└── public/              # Build/dist (auto-generated)
```

For more details about the architecture, see the [full documentation](docs/PROJECT_ARCHITECTURE.MD).

## Technologies

### Backend

* **Go 1.24+** – Programming language
* **PostgreSQL** – Relational database
* **JWT** – Authentication and authorization
* **Air** – Hot reload in development

### Frontend

* **Vanilla JavaScript** – No frameworks, pure JavaScript
* **ES Modules** – Native ES6 modules
* **Web Components** – Reusable components
* **Vite 5.4+** – Build tool and optimizations

### Machine Learning

* **TensorFlow / Keras** – NCF model training
* **Python 3.11+** – Training pipeline and embedding extraction
* **pgvector** – PostgreSQL extension for vector similarity search
* **NumPy / Pandas** – Data processing

### DevOps

* **Docker** – Containerization
* **Docker Compose** – Container orchestration
* **GitHub Actions** – Automated CI/CD
* **GitHub Container Registry** – Docker image registry

## Prerequisites

### For Local Development

* **Docker** 20.10+
* **Docker Compose** 2.0+
* **Node.js** 20+ and **npm** (optional, for local frontend development)

### For Production

* **Docker** 20.10+
* **Docker Compose** 2.0+

## Installation and Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd movies
```

### 2. Configure environment variables

Create a `.env` file at the project root based on `.env.example`:

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```env
# Database
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=movies_db

# Application
JWT_SECRET=your_very_secure_jwt_secret_here

# Optional (production)
DOCKER_REGISTRY=ghcr.io/your-user
VERSION=latest
APP_PORT=8080
```

## Development

### Option 1: Everything in Docker (Recommended)

This is the simplest and recommended way to develop:

```bash
# Start all services
docker-compose up -d --build
```

This will start:

* **Go backend** on port `8080` with hot reload (Air)
* **Vite frontend** in watch mode, automatically building to `public/`
* **PostgreSQL** on port `5432`

The application will be available at `http://localhost:8080`.

#### Initialize the database

On first run, you must populate the database:

```bash
docker exec movies-app-1 go run ./database/import/install.go
```

#### Train the recommendation model (optional)

To retrain the NCF model and regenerate embeddings:

```bash
docker compose --profile training up recommender
```

This trains the NeuMF model on user interaction data and exports 128-dim embeddings into PostgreSQL via pgvector. The Go backend uses these embeddings at runtime as part of its hybrid scoring algorithm.

#### Development workflow

* **Backend**: Changes to `.go` files automatically restart the server (Air)
* **Frontend**: Changes in `web/` are automatically built to `public/` (Vite watch mode)

### Option 2: Hybrid Development

For local frontend development (without Docker):
#### 1. Install frontend dependencies

```bash
cd web
npm install
```

#### 2. Run build in watch mode

```bash
npm run dev
```

This will watch changes in `web/` and automatically build to `public/`.

#### 3. Start only backend and database via Docker

```bash
docker-compose up postgres app -d
```

#### 4. Initialize the database

```bash
docker exec movies-app-1 go run ./database/import/install.go
```

---

## Production

### Production Docker Architecture

The project uses an optimized **multi-stage Dockerfile** for production:

```
┌─────────────────────────────────────────────────────────────┐
│                    MULTI-STAGE BUILD                        │
├─────────────────────────────────────────────────────────────┤
│  Stage 1: dev              │ Development environment        │
│  Stage 2: frontend-builder │ Frontend build (Vite)          │
│  Stage 3: backend-builder  │ Go compilation                 │
│  Stage 4: prod             │ Final image (~20MB)            │
└─────────────────────────────────────────────────────────────┘
```

### Security Features

| Feature                 | Description                                      |
| ----------------------- | ------------------------------------------------ |
| 🔒 Non-root user        | Container runs as `appuser` (UID 10001)          |
| 📁 Read-only filesystem | Filesystem in read-only mode                     |
| 🚫 no-new-privileges    | Prevents privilege escalation                    |
| 🗑️ CAP_DROP ALL        | Removes all Linux capabilities                   |
| 🌐 Isolated network     | Internal service network with no external access |
| 📊 Resource limits      | CPU and memory limits per container              |
| 🩺 Health checks        | Continuous service health checks                 |
| 📝 Structured logging   | Logs with automatic rotation                     |

### Manual Deploy with Docker Compose

```bash
# Build and start production containers
docker compose -f docker-compose.prod.yaml up -d --build

# Check container status
docker compose -f docker-compose.prod.yaml ps

# View logs in real time
docker compose -f docker-compose.prod.yaml logs -f

# Stop services
docker compose -f docker-compose.prod.yaml down
```

### Database Initialization

> **Automatic in Production**: The database is automatically initialized on first run!

`docker-compose.prod.yaml` mounts the `database-dump.sql` file into PostgreSQL’s `/docker-entrypoint-initdb.d/` directory. This causes the SQL script to run **automatically** when the database volume is created for the first time.

```yaml
# Configuration in docker-compose.prod.yaml
volumes:
  - ./server/database/import/database-dump.sql:/docker-entrypoint-initdb.d/01-init.sql:ro
```

**Behavior:**

* **First deploy**: Database is created and populated with ~4,800 movies
* **Subsequent deploys**: Volume persists and data is preserved
* **Database reset**: Use `docker compose -f docker-compose.prod.yaml down -v` to remove the volume and reinitialize

**Verify database initialization:**

```bash
# Check if tables exist
docker exec movies-postgres psql -U $POSTGRES_USER -d $POSTGRES_DB -c "\dt"

# Count records
docker exec movies-postgres psql -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT COUNT(*) FROM movies;"
```

### Production Environment Variables

Create a `.env` file with the following variables:

```env
# === REQUIRED ===
POSTGRES_USER=movies_prod
POSTGRES_PASSWORD=<strong-password-here>
POSTGRES_DB=movies_production
JWT_SECRET=<strong-256-bit-jwt-secret>

# === OPTIONAL ===
# Docker registry (for CI/CD)
DOCKER_REGISTRY=ghcr.io/your-user

# Image version (commit SHA or semantic tag)
VERSION=latest

# Application port (default: 8080)
APP_PORT=8080
```

### Container Health Check

```bash
# Check application health
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy"}

# Check all containers health
docker compose -f docker-compose.prod.yaml ps
```

---

## CI/CD with GitHub Actions

The project includes a complete CI/CD pipeline configured in `.github/workflows/ci-cd.yaml`.

### Pipeline Overview

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Test      │───▶│   Build     │───▶│   Scan      │───▶│   Deploy    │
│  Backend    │    │   Docker    │    │   Trivy     │    │  (manual)   │
│  Frontend   │    │   Image     │    │   SARIF     │    │             │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### Pipeline Features

| Stage             | Description                                             |
| ----------------- | ------------------------------------------------------- |
| **Test Backend**  | Go tests, linting, formatting checks                    |
| **Test Frontend** | Vite verification build                                 |
| **Build & Push**  | Multi-stage build and push to GitHub Container Registry |
| **Security Scan** | Vulnerability scanning with Trivy                       |
| **Deploy**        | Automatic deploy (configurable)                         |

### Pipeline Triggers

* **Push to `main`**: Runs full pipeline with deploy
* **Pull Request**: Runs tests only (no deploy)
* **Manual**: Can be triggered via GitHub UI

### GitHub Secrets Configuration

To enable the pipeline, configure the following secrets:

| Secret           | Description                         |
| ---------------- | ----------------------------------- |
| `GITHUB_TOKEN`   | Automatic (no setup required)       |
| `SERVER_HOST`    | Server IP/hostname (for SSH deploy) |
| `SERVER_USER`    | SSH user                            |
| `SERVER_SSH_KEY` | Private SSH key for deploy          |

### Automatic Deploy

The pipeline supports multiple deployment options:

#### Option 1: SSH Deploy to VPS/VM

Uncomment the section in the workflow and configure secrets:

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

#### Option 2: Deploy to PaaS Platforms

The pipeline can be adapted for:

* **Fly.io**: `flyctl deploy`
* **Railway**: Deploy API
* **Render**: Deploy webhook
* **DigitalOcean App Platform**: Deploy API

### Run Pipeline Manually

1. Go to **Actions** in the GitHub repository
2. Select **CI/CD Pipeline**
3. Click **Run workflow**
4. Choose the deploy environment

---

## NPM Scripts

In the `web/` folder:

| Script            | Description                               |
| ----------------- | ----------------------------------------- |
| `npm run dev`     | Watch-mode build for development          |
| `npm run build`   | Production build (optimized and minified) |
| `npm run preview` | Local preview of production build         |

## API Endpoints

### Health Check

* `GET /health` – Checks application health and database connection

### Authentication

* `POST /api/account/register/` – Register new user
* `POST /api/account/authenticate/` – Authenticate user (login)

### Movies

* `GET /api/movies/top` – List top 10 most popular movies
* `GET /api/movies/random` – List random movies
* `GET /api/movies/search?q={query}&order={order}&genre={genre}` – Search movies
* `GET /api/movies/{id}` – Get movie details
* `GET /api/genres` – List all genres

### Recommendations (Authentication required)

* `POST /api/movies/recommendations` – Get personalized recommendations for the authenticated user

### Collections (Authentication required)

* `GET /api/account/favorites/` – List favorite movies
* `GET /api/account/watchlist/` – List watchlist
* `POST /api/account/save-to-collection/` – Add movie to collection (also triggers recommendation recomputation)

**Authentication**: Protected endpoints require header `Authorization: Bearer {token}`

## Tests

*Section reserved for future test implementation*

## Additional Documentation

* [Project Architecture](docs/PROJECT_ARCHITECTURE.MD) – Clean Architecture and DDD details
* [Entity-Relationship Diagram](docs/ENTITY_RELATION_DIAGRAM.MD) – Database structure
* [Frontend Performance Guide](docs/FRONTEND_PERFORMANCE_GUIDE.md) – Optimizations and best practices

## Useful Commands
### Development

```bash
# Start development environment
docker compose up -d --build

# View logs in real time
docker compose logs -f app

# Run command inside container
docker exec -it movies-app-1 sh

# Rebuild backend only
docker compose up -d --build app
```

### Production

```bash
# Production build
docker compose -f docker-compose.prod.yaml build

# Deploy new version
VERSION=v1.0.0 docker compose -f docker-compose.prod.yaml up -d

# Check container resource usage
docker stats

# Database backup
docker exec movies-postgres pg_dump -U $POSTGRES_USER $POSTGRES_DB > backup.sql
```

### Maintenance

```bash
# Clean unused images
docker image prune -a

# Clean orphan volumes
docker volume prune

# Check disk usage
docker system df

# System logs
docker compose -f docker-compose.prod.yaml logs --tail=100
```

## Stop the Application

To stop and remove containers:

```bash
# Development
docker compose down

# Production
docker compose -f docker-compose.prod.yaml down
```

To also remove volumes (database data):

```bash
docker compose down -v
```

## Data Structure

### Main Entities

* **Movie** – Movie information (title, synopsis, cast, etc.)
* **User** – System users
* **Actor** – Actors/actresses
* **Genre** – Movie genres
* **UserMovie** – Relationship between users and movies (favorites/watchlist)
* **MovieEmbedding** – 128-dim vector embeddings for movies (pgvector)
* **UserEmbedding** – Aggregated user taste vectors (pgvector)
* **UserRecommendation** – Cached personalized recommendations per user

## Security

### Application

* Passwords hashed with bcrypt
* JWT-based authentication (JSON Web Tokens)
* Backend data validation (Value Objects)
* Input sanitization

### Containers (Production)

* Non-root user in all containers
* Read-only filesystem
* Linux capabilities removed
* Resource limits (CPU/memory)
* Isolated service network
* Active health checks
* Logging with automatic rotation
