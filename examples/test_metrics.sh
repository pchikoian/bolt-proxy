#!/bin/bash
#
# Test script to demonstrate bolt-proxy Prometheus metrics
#
# Usage: ./test_metrics.sh
#

set -e

METRICS_URL="http://localhost:9090/metrics"
HEALTH_URL="http://localhost:8080/health"

echo "=========================================="
echo "Bolt Proxy Metrics Test"
echo "=========================================="
echo ""

# Check if metrics endpoint is accessible
echo "1. Checking metrics endpoint..."
if curl -s -f "$METRICS_URL" > /dev/null; then
    echo "✓ Metrics endpoint is accessible at $METRICS_URL"
else
    echo "✗ Metrics endpoint is not accessible"
    exit 1
fi
echo ""

# Show current metrics
echo "2. Current connection metrics:"
curl -s "$METRICS_URL" | grep -E "bolt_proxy_connections_total|bolt_proxy_active_connections" | grep -v "^#"
echo ""

echo "3. Current health check metrics:"
curl -s "$METRICS_URL" | grep "bolt_proxy_health_checks_total" | grep -v "^#"
echo ""

# Generate some traffic
echo "4. Generating test traffic (5 health checks)..."
for i in {1..5}; do
    curl -s "$HEALTH_URL" > /dev/null
    echo "  Health check $i/5"
done
echo ""

# Show updated metrics
echo "5. Updated health check metrics:"
curl -s "$METRICS_URL" | grep "bolt_proxy_health_checks_total" | grep -v "^#"
echo ""

# Show all available metrics
echo "6. All available bolt-proxy metrics:"
curl -s "$METRICS_URL" | grep "^bolt_proxy_" | grep -v "^#" | wc -l | awk '{print "  Total metric lines: " $1}'
echo ""

# Show metric types
echo "7. Metric categories:"
curl -s "$METRICS_URL" | grep "^# HELP bolt_proxy_" | sed 's/# HELP /  - /' | cut -d' ' -f1-2
echo ""

echo "=========================================="
echo "Metrics test completed!"
echo ""
echo "View all metrics at: $METRICS_URL"
echo "=========================================="
