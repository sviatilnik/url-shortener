#!/bin/bash

# Скрипт для запуска всех бенчмарков и сохранения результатов

set -e

echo "Running comprehensive benchmarks..."

echo "=== Shortener Benchmarks ==="
go test -bench=. -benchmem ./internal/app/shortener/ > benchmarks/shortener_benchmarks.txt

echo "=== Generator Benchmarks ==="
go test -bench=. -benchmem ./internal/app/generators/ > benchmarks/generators_benchmarks.txt

echo "=== Storage Benchmarks ==="
go test -bench=. -benchmem ./internal/app/storages/ > benchmarks/storage_benchmarks.txt

echo "=== Memory Usage Analysis ==="
go test -bench=. -benchmem -memprofile=profiles/memory.pprof ./internal/app/shortener/
go test -bench=. -benchmem -memprofile=profiles/generators_memory.pprof ./internal/app/generators/
go test -bench=. -benchmem -memprofile=profiles/storage_memory.pprof ./internal/app/storages/

echo "Benchmarks completed and saved to benchmarks/ directory"
