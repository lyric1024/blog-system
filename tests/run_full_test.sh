#!/usr/bin/env bash
set -euo pipefail

# This script starts the server (detached), waits for it to become healthy,
# runs the curl test script, collects outputs, and then stops the server.
# Run locally: bash tests/run_full_test.sh

BASE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BASE_URL="http://127.0.0.1:8080"

echo "Starting server in background..."
nohup env GIN_MODE=release go run "$BASE_DIR"/main.go > /tmp/blog_system.log 2>&1 &
PID=$!
echo "Server PID: $PID"

echo "Waiting up to 30s for server to become available..."
for i in {1..30}; do
  # Accept any HTTP response as success (server listening). Avoid -f which fails on 404.
  if curl -sS --connect-timeout 2 "$BASE_URL/healthz" >/dev/null 2>&1; then
    echo "Server is reachable"
    break
  fi
  sleep 1
done

if ! kill -0 "$PID" >/dev/null 2>&1; then
  echo "Server process is not running. Check /tmp/blog_system.log for details." >&2
  exit 1
fi

echo "Running curl tests..."
bash "$BASE_DIR/tests/run_curl_tests.sh" || true

echo "Tests finished. Server log tail (last 200 lines):"
tail -n 200 /tmp/blog_system.log || true

echo "Stopping server PID $PID..."
kill "$PID" || true
sleep 1
if kill -0 "$PID" >/dev/null 2>&1; then
  echo "Server did not stop; sending SIGKILL"
  kill -9 "$PID" || true
fi

echo "Done. Test artifacts: /tmp/_login.json /tmp/_reg.json /tmp/_body (last request body)"
