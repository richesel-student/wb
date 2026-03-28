#!/bin/bash

set -e

PORT=8080

# освободить порт
free_port() {
  PID=$(lsof -ti :$PORT || true)
  if [ -n "$PID" ]; then
    echo "⚠️ Port $PORT busy. Killing process $PID..."
    kill -9 $PID
  fi
}

case "$1" in
  run)
    make run
    ;;

  test)
    echo "🧪 Running API tests..."
    make test-api
    ;;

  check)
    echo "🔍 Running checks (vet + lint + unit tests)..."
    make check
    ;;

  dev)
    echo "🔥 Full dev cycle..."

    free_port

    echo "🔍 Step 1: Code checks..."
    make check

    echo "🚀 Step 2: Starting server..."
    make run &
    SERVER_PID=$!

    sleep 1

    echo "🧪 Step 3: API tests..."
    make test-api

    echo "🛑 Stopping server..."
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true

    echo "✅ Dev cycle complete!"
    ;;

  *)
    echo "Usage:"
    echo "  ./dev.sh run     # run server"
    echo "  ./dev.sh test    # API tests"
    echo "  ./dev.sh check   # vet + lint + unit tests"
    echo "  ./dev.sh dev     # full pipeline"
    ;;
esac