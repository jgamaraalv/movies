# ============================================
# Stage 1: Base para desenvolvimento
# ============================================
FROM golang:1.24-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/cosmtrek/air@v1.49.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# ============================================
# Stage 2: Builder para produção
# ============================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build otimizado para produção
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./main.go

# ============================================
# Stage 3: Imagem final de produção (mínima)
# ============================================
FROM alpine:3.19 AS prod

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copia apenas o binário compilado
COPY --from=builder /app/server .
# Copia arquivos estáticos
COPY --from=builder /app/public ./public

# Usuário não-root para segurança
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]