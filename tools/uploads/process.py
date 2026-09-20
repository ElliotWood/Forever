#!/usr/bin/python

# Reads whatever players have uploaded and says what is in it, so a contribution does not sit
# in a bucket until somebody remembers to go and look. Run it on a timer.
#
#   python tools/uploads/process.py             # new uploads only
#   python tools/uploads/process.py --all       # everything, ignoring what has been seen
#   python tools/uploads/process.py --quiet     # no Discord post, just stdout
#
# Needs the worker's list token (~/.openclaw-alfred/secrets/forever_list_token.txt) and, to
# announce anything, the webhook (forever_webhook.url). Neither lives in this repo.
#
# The files themselves are kept rather than processed and dropped: a hotfix cache is the only
# evidence of what Blizzard changed on a given day, and re-downloading one is not possible
# once a sender's client has moved on.

import argparse
import json
import os
import subprocess
import sys
import urllib.parse
import urllib.request

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
WORKER = 'https://forever-uploads.gigaflare-elliot.workers.dev'
SECRETS = os.path.expanduser('~/.openclaw-alfred/secrets')
STATE = os.path.expanduser('~/.openclaw-alfred/workspace/triage-state/forever_uploads_seen.json')
KEEP = os.path.expanduser('~/.openclaw-alfred/workspace/scratch/forever-uploads')
AGENT = 'DiscordBot (https://github.com/ElliotWood/Forever, 1.0)'

# What to run on each kind, and how much of its output is worth repeating in a chat message.
READERS = {
    'dbcache': ([sys.executable, os.path.join(REPO, 'tools/data_watch/hotfix_cache.py'), '--cache'], 4),
    'damagemeter': ([sys.executable, os.path.join(REPO, 'tools/data_watch/damage_meter.py'), '--file'], 8),
    # A screenshot is for a human to look at; there is nothing to parse.
    'screenshot': (None, 0),
}


def secret(name):
    path = os.path.join(SECRETS, name)
    return open(path).read().strip() if os.path.exists(path) else ''


def fetch(path, token):
    request = urllib.request.Request(f'{WORKER}{path}', headers={'Authorization': f'Bearer {token}', 'User-Agent': AGENT})
    with urllib.request.urlopen(request, timeout=120) as response:
        return response.read()


def announce(text):
    url = secret('forever_webhook.url')
    if not url:
        return
    request = urllib.request.Request(
        url, data=json.dumps({'content': text[:1900]}).encode(),
        headers={'Content-Type': 'application/json', 'User-Agent': AGENT}, method='POST')
    try:
        urllib.request.urlopen(request, timeout=20)
    except Exception as error:
        print(f'discord: {error}', file=sys.stderr)


def read_state():
    try:
        return set(json.load(open(STATE))['seen'])
    except Exception:
        return set()


def write_state(seen):
    os.makedirs(os.path.dirname(STATE), exist_ok=True)
    json.dump({'seen': sorted(seen)}, open(STATE, 'w'), indent='\t')


def summarise(kind, path):
    """The first few lines of the matching reader, or why there are none."""
    command, lines = READERS.get(kind, (None, 0))
    if not command:
        return '(nothing to parse)'
    done = subprocess.run(command + [path], capture_output=True, text=True, cwd=REPO)
    if done.returncode != 0:
        return f'reader failed: {done.stderr.strip().splitlines()[-1] if done.stderr.strip() else "no output"}'
    body = [line for line in done.stdout.splitlines() if line.strip()]
    # The first line is the path we just passed in, which the reader echoes back.
    return '\n'.join(body[1 : 1 + lines]) or '(reader said nothing)'


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--all', action='store_true', help='reprocess everything, not just what is new')
    parser.add_argument('--quiet', action='store_true', help='do not post to Discord')
    args = parser.parse_args()

    token = secret('forever_list_token.txt')
    if not token:
        sys.exit('no list token: expected ~/.openclaw-alfred/secrets/forever_list_token.txt')

    objects = json.loads(fetch('/list', token))['objects']
    seen = set() if args.all else read_state()
    fresh = [o for o in objects if o['key'] not in seen]
    if not fresh:
        print('nothing new')
        return

    os.makedirs(KEEP, exist_ok=True)
    reports = []
    for obj in sorted(fresh, key=lambda o: o['uploaded']):
        kind = obj['key'].split('/')[0]
        path = os.path.join(KEEP, obj['key'].replace('/', '_'))
        if not os.path.exists(path):
            open(path, 'wb').write(fetch(f'/get?key={urllib.parse.quote(obj["key"])}', token))
        note = f'\nnote: {obj["note"]}' if obj['note'] else ''
        reports.append(f'**{kind}** · {obj["size"] / 1e6:.1f} MB · {obj["uploaded"][:16]}{note}\n```\n{summarise(kind, path)}\n```')
        seen.add(obj['key'])

    body = f'**{len(fresh)} new upload{"s" if len(fresh) != 1 else ""}**\n' + '\n'.join(reports)
    print(body)
    if not args.quiet:
        announce(body)
    write_state(seen)


if __name__ == '__main__':
    main()
