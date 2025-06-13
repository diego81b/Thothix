FROM golang:1.23-alpine AS builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage finale con UBI Micro (più piccola)
FROM registry.access.redhat.com/ubi9/ubi-micro

# Metadata dell'immagine
LABEL org.opencontainers.image.title="Thothix API"
LABEL org.opencontainers.image.description="API Backend per la piattaforma di messaggistica aziendale Thothix"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.vendor="Thothix"
LABEL org.opencontainers.image.authors="Thothix Team"
LABEL org.opencontainers.image.source="https://github.com/thothix/thothix"
LABEL org.opencontainers.image.documentation="https://docs.thothix.com"

WORKDIR /app

# Copia il binario dal builder stage
COPY --from=builder /app/backend/main .

# Espone la porta
EXPOSE 30000
USER 1001
CMD ["./main"]
