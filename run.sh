#!/usr/bin/env bash
# Rebuild and (re)start the SSH site. Usage: ./run.sh
set -e

# Make sure `go` is reachable even from a fresh shell.
export PATH="/opt/homebrew/bin:$PATH"

cd "$(dirname "$0")"

echo "› stopping any running server…"
pkill -f '/sshsite$' 2>/dev/null || true
sleep 1

echo "› building…"
go build -o sshsite ./...

echo "› starting on port 2222 (ssh -p 2222 localhost)…"
./sshsite &.
