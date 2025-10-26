#!/bin/bash
#
# Simple shell script to test bolt-proxy connection using cypher-shell
#
# Requirements:
#   - cypher-shell (from Neo4j)
#   - bolt-proxy running
#
# Usage:
#   ./test_bolt_proxy.sh [HOST] [PORT] [USERNAME] [PASSWORD]
#

set -e

HOST="${1:-localhost}"
PORT="${2:-8080}"
USERNAME="${3:-}"
PASSWORD="${4:-}"

BOLT_URI="bolt://${HOST}:${PORT}"

echo "=========================================="
echo "Bolt Proxy Connection Test"
echo "=========================================="
echo "URI: ${BOLT_URI}"
echo "User: ${USERNAME:-(none)}"
echo ""

# Check if cypher-shell is available
if ! command -v cypher-shell &> /dev/null; then
    echo "Error: cypher-shell not found. Please install Neo4j or cypher-shell."
    echo "Visit: https://neo4j.com/docs/operations-manual/current/tools/cypher-shell/"
    exit 1
fi

# Build auth parameter
AUTH_PARAM=""
if [ -n "$USERNAME" ]; then
    AUTH_PARAM="-u ${USERNAME} -p ${PASSWORD}"
fi

# Test 1: Basic connection
echo "Test 1: Basic connection..."
if echo "RETURN 1 AS num;" | cypher-shell -a "${BOLT_URI}" ${AUTH_PARAM} --format plain 2>/dev/null | grep -q "1"; then
    echo "✓ Connection successful"
else
    echo "✗ Connection failed"
    exit 1
fi

# Test 2: Create node
echo ""
echo "Test 2: Creating test node..."
QUERY="CREATE (n:TestNode {name: 'BoltProxyTest', timestamp: timestamp()}) RETURN n;"
if echo "${QUERY}" | cypher-shell -a "${BOLT_URI}" ${AUTH_PARAM} --format plain &>/dev/null; then
    echo "✓ Node created successfully"
else
    echo "✗ Node creation failed"
    exit 1
fi

# Test 3: Query node
echo ""
echo "Test 3: Querying test node..."
QUERY="MATCH (n:TestNode {name: 'BoltProxyTest'}) RETURN n LIMIT 5;"
if echo "${QUERY}" | cypher-shell -a "${BOLT_URI}" ${AUTH_PARAM} --format plain; then
    echo "✓ Query successful"
else
    echo "✗ Query failed"
    exit 1
fi

# Test 4: Cleanup
echo ""
echo "Test 4: Cleaning up test data..."
QUERY="MATCH (n:TestNode {name: 'BoltProxyTest'}) DELETE n;"
if echo "${QUERY}" | cypher-shell -a "${BOLT_URI}" ${AUTH_PARAM} --format plain &>/dev/null; then
    echo "✓ Cleanup successful"
else
    echo "✗ Cleanup failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "All tests passed!"
echo "=========================================="
