FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /usr/local/bin/nfrecap cmd/nfrecap/main.go

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/bin/nfrecap /usr/local/bin/nfrecap

# Copy .env if it exists, though environment variables should officially be passed via compose
# COPY .env . 

EXPOSE 8080
CMD ["nfrecap", "serve"]
