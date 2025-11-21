#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://127.0.0.1:8080"
USERNAME="lyric"
PASSWORD="123456"

echo "Wait for server ${BASE_URL} to be available (timeout 30s)..."
for i in {1..30}; do
  if curl -sS --connect-timeout 2 "$BASE_URL/healthz" >/dev/null 2>&1; then
    echo "Server is reachable"
    break
  fi
  sleep 1
done

echo "Attempt login as $USERNAME..."
LOGIN_RESP_FILE=/tmp/_login.json
HTTP_LOGIN=$(curl -s -o "$LOGIN_RESP_FILE" -w "%{http_code}" -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"userName\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

echo "Login HTTP status: $HTTP_LOGIN"
cat "$LOGIN_RESP_FILE" || true

TOKEN=$(python3 -c 'import sys,json; d=json.load(open("/tmp/_login.json")); print(d.get("data",{}).get("token",""))' 2>/dev/null || true)

if [ -z "$TOKEN" ]; then
  echo "❌ Login failed. Please ensure user '$USERNAME' is already registered."
  echo "Response was:"
  cat "$LOGIN_RESP_FILE" || true
  exit 1
fi

echo "✅ Token obtained (truncated): ${TOKEN:0:20}..."

run_req() {
  local name="$1"; shift
  echo -e "\n== $name =="
  curl -s -D /tmp/_hdr -o /tmp/_body -w "HTTP_STATUS:%{http_code}\n" "$@"
  cat /tmp/_hdr
  echo
  cat /tmp/_body || true
}

# Create a post
run_req "Create Post" -X POST "$BASE_URL/api/post/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"My First Post","content":"Hello world from script"}'

# List posts
run_req "List Posts" -X POST "$BASE_URL/api/post/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}'

# Create comment (attempt postID=1)
run_req "Create Comment" -X POST "$BASE_URL/api/comment/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"postID":1,"content":"Nice post"}'

# List comments
run_req "List Comments" -X POST "$BASE_URL/api/comment/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"postID":1}'

echo -e "\n✅ All tests completed. Last response body in /tmp/_body, headers in /tmp/_hdr."