#!/usr/bin/env sh

# entrypoint to Dockerfile-git

cd /app

chmod 0600 /root/.ssh/ed25519
eval $(ssh-agent -s)

ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
ssh-add /root/.ssh/ed25519

git "$@"