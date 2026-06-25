#!/usr/bin/env bash
# Claude Code PostToolUse hook: gofmt the file Claude just edited.
# Reads the hook JSON payload from stdin, formats it if it is a Go file.
set -euo pipefail

payload=$(cat)
file=$(printf '%s' "$payload" | jq -r '.tool_input.file_path // empty')

[ -n "$file" ] || exit 0
[ "${file##*.}" = "go" ] || exit 0
[ -f "$file" ] || exit 0

gofmt -w "$file"
