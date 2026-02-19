#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"

# Ensure Fly CLI
if ! command -v fly &>/dev/null; then
  echo "Fly CLI not found. Install: https://fly.io/docs/hands-on/install-flyctl/"
  exit 1
fi

# Create app from fly.toml if needed (no-op if already exists)
fly launch --no-deploy

# Create volume if it doesn't exist (idempotent)
REGION=$(grep -E '^primary_region\s*=' fly.toml | sed -E 's/.*=\s*"([^"]+)".*/\1/' || echo "ord")
if ! fly volumes list 2>/dev/null | grep -q "data"; then
  fly volumes create data --size 1 --region "$REGION"
fi

# Remind about secrets (only first time)
if ! fly secrets list 2>/dev/null | grep -q "STEAM_API_KEY"; then
  echo "Set secrets before first deploy: fly secrets set STEAM_API_KEY=xxx SESSION_SECRET=xxx"
  read -p "Continue to deploy anyway? [y/N] " -n 1 -r; echo
  [[ $REPLY =~ ^[Yy]$ ]] || exit 0
fi

fly deploy
