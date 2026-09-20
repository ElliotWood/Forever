#!/bin/sh
# Prints every upload sitting in the bucket: key, size, when it arrived, the sender's note.
# Nothing consumes these automatically - this is the "go and look" the upload worker assumes.
set -e
cd "$(dirname "$0")"
npx wrangler dev --remote --port 8799 --ip 127.0.0.1 >/dev/null 2>&1 &
pid=$!
trap 'kill $pid 2>/dev/null' EXIT
until curl -sf http://127.0.0.1:8799/ >/dev/null 2>&1; do sleep 2; done
curl -s http://127.0.0.1:8799/
