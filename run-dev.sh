#!/bin/bash
set -e

# Always rollback shell options before exiting or returning
trap "set +e" EXIT RETURN

ENV_FILE=.env

if [ -f "$ENV_FILE" ]; then
  export $(cat $ENV_FILE | xargs)
  echo "$ENV_FILE file configured"
else
  echo "$ENV_FILE file does not exist."
  exit
fi

echo "-------"
echo "starting stori-test service"
echo "-------"
echo "[+] Run containers ${@} "
docker compose up ${@}

echo "[+] Cleaning up stopped containers..."
docker ps --all --filter status=exited -q | xargs docker rm -v;
