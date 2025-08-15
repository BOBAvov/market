# Stage 1: build
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /src

# Сначала зависимости (быстрее кэш)
COPY go.mod go.sum ./
RUN go mod download

# Дальше код
COPY . .

# Сборка статического бинарника
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# Stage 2: runtime
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/api /app/api
COPY --from=builder /src/config /app/config/
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/api"]