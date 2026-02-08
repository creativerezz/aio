FROM golang:1.25-alpine AS builder

WORKDIR /src

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /aio ./cmd/aio

FROM alpine:latest

LABEL org.opencontainers.image.description="Aio AI augmentation CLI and REST API server. See https://github.com/creativerezz/aio for details."

RUN apk add --no-cache ca-certificates curl yt-dlp \
    && mkdir -p /root/.config/aio

COPY --from=builder /aio /usr/local/bin/aio

# Bake default patterns and strategies into the image
COPY data/patterns/ /root/.config/aio/patterns/
COPY data/strategies/ /root/.config/aio/strategies/

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:${PORT:-8080}/patterns || exit 1

ENTRYPOINT ["aio"]
CMD ["--serve"]
