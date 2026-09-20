#!/bin/sh
# Every upload sitting in the bucket. R2 has no object listing outside a worker binding, so the
# worker grew a token-gated /list endpoint rather than this needing wrangler and a dev preview.
set -e
curl -sf -H "Authorization: Bearer $(cat ~/.openclaw-alfred/secrets/forever_list_token.txt)" \
    https://forever-uploads.gigaflare-elliot.workers.dev/list |
    python -c "import json,sys; [print(f'{o[\"uploaded\"][:16]}  {o[\"size\"]:>9}  {o[\"key\"]}  {o[\"note\"]}') for o in json.load(sys.stdin)['objects']]"
