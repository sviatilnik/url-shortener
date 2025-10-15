#!/bin/bash

# Скрипт для сравнения профилей памяти

set -e

echo "Comparing memory profiles..."

echo "=== Base Profile (Before Optimization) ==="
go tool pprof -alloc_space -top profiles/base.pprof

echo ""
echo "=== Result Profile (After Optimization) ==="
go tool pprof -alloc_space -top profiles/result.pprof

echo ""
echo "=== Comparison (Base vs Result) ==="
go tool pprof -alloc_space -top -diff_base=profiles/base.pprof profiles/result.pprof

echo ""
echo "=== Memory Usage Comparison ==="
echo "Base profile total memory:"
go tool pprof -alloc_space -top profiles/base.pprof | head -1

echo "Result profile total memory:"
go tool pprof -alloc_space -top profiles/result.pprof | head -1
