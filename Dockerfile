# ============================================
# Stage 1: Base para desenvolvimento
# ============================================
FROM golang:1.24-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git bash && \
    go install github.com/cosmtrek/air@v1.49.0

# Copy Go modules (for initial download, volumes will override)
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy .air.toml (needed for air to work)
COPY .air.toml ./

EXPOSE 8080

# Note: server/ and web/ are mounted via volumes in docker-compose
# The build script should be run manually or via a separate step
CMD ["air", "-c", ".air.toml"]

# ============================================
# Stage 2: Builder para produção
# ============================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy Go modules
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy server code
COPY server/ ./

# Build otimizado para produção
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/api/main.go

# Build frontend: copy web source and build script, then build
COPY web/ /build-frontend/web/
COPY build.sh /build-frontend/build.sh
RUN chmod +x /build-frontend/build.sh && \
    cd /build-frontend && \
    ./build.sh && \
    mv public /app/public

# ============================================
# Stage 3: Imagem final de produção (mínima)
# ============================================
FROM alpine:3.19 AS prod

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copia apenas o binário compilado
COPY --from=builder /app/server .
# Copia arquivos estáticos (build do frontend)
COPY --from=builder /app/public ./public

# Usuário não-root para segurança
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]