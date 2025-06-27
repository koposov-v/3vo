#!/bin/bash

services=("api-gateway" "auth-service" "kitchen-queue" "order-core")

for service in "${services[@]}"; do
  env_example="$service/.env.example"
  env_file="$service/.env"

  if [[ -f "$env_example" ]]; then
    if [[ -f "$env_file" ]]; then
      echo "[✓] $env_file уже существует"
    else
      cp "$env_example" "$env_file"
      echo "[+] Создан $env_file из $env_example"
    fi
  else
    echo "[!] Внимание: $env_example не найден, пропускаем $service"
  fi
done