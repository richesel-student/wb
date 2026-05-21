#!/usr/bin/env bash

set -e

echo "===> Building project..."
go build -o gc-monitor .

echo "===> Starting server..."
./gc-monitor > /tmp/gc-monitor.log 2>&1 &
PID=$!

cleanup() {
    echo ""
    echo "===> Stopping server..."
    kill $PID 2>/dev/null || true
}

trap cleanup EXIT

sleep 2

echo ""
echo "===> Testing /metrics endpoint..."

METRICS=$(curl -s http://localhost:8080/metrics)

if [[ -z "$METRICS" ]]; then
    echo "❌ /metrics returned empty response"
    exit 1
fi

echo "✅ /metrics is available"

echo ""
echo "===> Checking required metrics..."

REQUIRED_METRICS=(
    "gc_monitor_mem_alloc_bytes"
    "gc_monitor_mem_total_alloc_bytes_total"
    "gc_monitor_mem_sys_bytes"
    "gc_monitor_gc_cycles_total"
    "gc_monitor_gc_last_time_seconds"
)

for metric in "${REQUIRED_METRICS[@]}"; do
    if echo "$METRICS" | grep -q "$metric"; then
        echo "✅ Found metric: $metric"
    else
        echo "❌ Missing metric: $metric"
        exit 1
    fi
done

echo ""
echo "===> Testing GET /gc-percent..."

GC_RESPONSE=$(curl -s http://localhost:8080/gc-percent)

echo "Response:"
echo "$GC_RESPONSE"

echo ""
echo "===> Testing POST /gc-percent..."

POST_RESPONSE=$(curl -s -X POST "http://localhost:8080/gc-percent?value=50")

echo "Response:"
echo "$POST_RESPONSE"

echo ""
echo "===> Verifying updated GC percent..."

VERIFY_RESPONSE=$(curl -s http://localhost:8080/gc-percent)

echo "Response:"
echo "$VERIFY_RESPONSE"

if echo "$VERIFY_RESPONSE" | grep -q "50"; then
    echo "✅ GC percent updated successfully"
else
    echo "❌ GC percent update failed"
    exit 1
fi

echo ""
echo "===> Testing pprof endpoint..."

PPROF_RESPONSE=$(curl -s http://localhost:8080/debug/pprof/ | head -20)

if echo "$PPROF_RESPONSE" | grep -qi "profile"; then
    echo "✅ pprof endpoint works"
else
    echo "❌ pprof endpoint failed"
    exit 1
fi

echo ""
echo "===> Running go vet..."
go vet ./...

echo ""
echo "===> Running race detector..."
go test -race ./...

echo ""
echo "🎉 ALL TESTS PASSED SUCCESSFULLY"