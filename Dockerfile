# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /aio ./cmd/aio

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates curl yt-dlp && \
    mkdir -p /root/.config/aio

COPY --from=builder /aio /usr/local/bin/aio

# Bake patterns and strategies into the image
COPY data/patterns/ /root/.config/aio/patterns/
COPY data/strategies/ /root/.config/aio/strategies/

ENV AIO_CONFIG_DIR=/root/.config/aio

EXPOSE 8080

CMD ["aio", "--serve", "--address", ":8080"]
