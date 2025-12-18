FROM golang:1.24-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git bash && \
    go install github.com/cosmtrek/air@v1.49.0

COPY server/go.mod server/go.sum ./

RUN go mod download

COPY .air.toml ./

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

FROM node:20-alpine AS frontend-builder

LABEL stage="frontend-builder"
LABEL maintainer="Movies App Team"

WORKDIR /build-frontend

COPY web/package*.json ./

RUN npm ci --frozen-lockfile

COPY web/ ./

ENV NODE_ENV=production

RUN mkdir -p dist && \
    npx vite build --outDir ./dist


FROM golang:1.24-alpine AS backend-builder

LABEL stage="backend-builder"

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata git

COPY server/go.mod server/go.sum ./

RUN go mod download && \
    go mod verify

COPY server/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -o /app/server \
    ./cmd/api/main.go

COPY --from=frontend-builder /build-frontend/dist ./public

FROM alpine:3.21 AS prod

LABEL org.opencontainers.image.title="Movies App"
LABEL org.opencontainers.image.description="Full-stack movie listing application"
LABEL org.opencontainers.image.authors="Movies App Team"
LABEL org.opencontainers.image.source="https://github.com/jgamaraalv/movies"
LABEL org.opencontainers.image.licenses="MIT"

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S -G appgroup -h /app -s /sbin/nologin appuser

COPY --from=backend-builder --chown=appuser:appgroup /app/server .

COPY --from=backend-builder --chown=appuser:appgroup /app/public ./public

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
