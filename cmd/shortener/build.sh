#!/bin/bash

# Скрипт для сборки приложения с информацией о версии

VERSION=${1:-"1.0.0"}
DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")


go build -ldflags "-X main.buildVersion=$VERSION -X main.buildDate=$DATE -X main.buildCommit=$COMMIT" -o shortener main.go

echo "Build completed: bin/shortener"
