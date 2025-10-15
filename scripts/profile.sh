#!/bin/bash

# Скрипт для профилирования и анализа производительности

set -e

echo "Building profiler..."
go build -o profiler ./cmd/profiler

echo "Starting profiler..."
./profiler &
PROFILER_PID=$!

# Ждем запуска profiler
sleep 2

echo "Collecting base profile..."
go tool pprof -proto http://localhost:6060/debug/pprof/heap > profiles/base.pprof

echo "Stopping profiler..."
kill $PROFILER_PID

echo "Running benchmarks..."
go test -bench=. -benchmem ./internal/app/shortener/
go test -bench=. -benchmem ./internal/app/generators/
go test -bench=. -benchmem ./internal/app/storages/

echo "Profile collection completed"
