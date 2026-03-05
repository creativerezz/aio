#!/bin/sh
set -e

# Create .env file if it doesn't exist
if [ ! -f /root/.config/aio/.env ]; then
  echo "Creating empty .env file..."
  touch /root/.config/aio/.env
fi

# Check if patterns are installed, if not, download them
if [ ! -d /root/.config/aio/patterns ] || [ -z "$(find /root/.config/aio/patterns -maxdepth 1 -type d | tail -n +2)" ]; then
  echo "Patterns not found. Downloading patterns..."
  exec /app/aio -U
  echo "Patterns downloaded successfully!"
fi

# Start the REST API server
exec /app/aio --serve --address :8080
