# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /app/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/migrate ./cmd/migrate

# Run stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl

RUN adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/migrate .
COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./api"]
