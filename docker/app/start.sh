#!/bin/sh

if [ "$DEBUG" = "true" ]; then
  echo "⚙️Starting Go app at $PORT and (debug mode) at 40000..."
  mkdir -p .bin

  exec dlv debug ./ \
    --headless \
    --listen=:40000 \
    --accept-multiclient \
    --api-version=2 \
    --log
else
  echo "🚀 Starting Go app at $PORT..."
  exec go run .
fi