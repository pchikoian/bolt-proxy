# Bolt-Proxy Examples

This directory contains example scripts to test and demonstrate the bolt-proxy functionality.

## Test Scripts

### 1. Python Test Script

**File:** `test_bolt_proxy.py`

**Requirements:**
```bash
pip install neo4j
```

**Usage:**
```bash
# Basic test (no auth)
python test_bolt_proxy.py

# With custom host/port
python test_bolt_proxy.py --host localhost --port 8080

# With authentication
python test_bolt_proxy.py --user myuser --password mypass

# Full example
python test_bolt_proxy.py --host localhost --port 8080 --user neo4j --password password
```

**What it tests:**
- Basic connection to bolt-proxy
- Node creation through the proxy
- Query execution
- Cleanup operations

---

### 2. Shell Script Test

**File:** `test_bolt_proxy.sh`

**Requirements:**
- `cypher-shell` (installed with Neo4j)

**Installation (macOS):**
```bash
brew install neo4j
```

**Usage:**
```bash
# Make executable
chmod +x test_bolt_proxy.sh

# Basic test
./test_bolt_proxy.sh

# With parameters: HOST PORT USERNAME PASSWORD
./test_bolt_proxy.sh localhost 8080 neo4j password
```

---

### 3. Go Test Script

**File:** `test_bolt_proxy.go`

**Requirements:**
- Go 1.22 or later
- Neo4j Go Driver (will be downloaded automatically)

**Usage:**
```bash
# Build and run
go run test_bolt_proxy.go

# With custom parameters
go run test_bolt_proxy.go -host localhost -port 8080 -user neo4j -password password

# Or build first, then run
go build -o test_bolt_proxy test_bolt_proxy.go
./test_bolt_proxy -host localhost -port 8080
```

---

## Quick Start with Docker Compose

### Testing with Memgraph

```bash
# Start bolt-proxy with Memgraph
cd ..
docker-compose up -d

# Wait a few seconds for services to start
sleep 5

# Run Python test
python examples/test_bolt_proxy.py --host localhost --port 8080
```

### Testing with Neo4j

```bash
# Start bolt-proxy with Neo4j
cd ..
docker-compose -f docker-compose.neo4j.yml up -d

# Wait for Neo4j to be ready (may take 30-60 seconds)
sleep 30

# Run Python test with Neo4j credentials
python examples/test_bolt_proxy.py --host localhost --port 8080 --user neo4j --password password
```

---

## Expected Output

All test scripts should produce output similar to:

```
==========================================================
Bolt Proxy Test Suite
==========================================================
Testing connection...
✓ Connection successful

Testing node creation...
✓ Created node: {...}

Testing node query...
✓ Found 1 node(s)

Cleaning up test data...
✓ Cleanup successful

==========================================================
Results: 4 passed, 0 failed
==========================================================
```

---

## Troubleshooting

### Connection Refused
- Ensure bolt-proxy is running: `docker-compose ps`
- Check the port is correct (default: 8080)
- Verify the backend service (Memgraph/Neo4j) is healthy

### Authentication Failed
- For Neo4j: Use the credentials set in `docker-compose.neo4j.yml` (default: neo4j/password)
- For Memgraph: Usually no authentication required unless configured
- Check if AUTH_METHOD is set in bolt-proxy environment variables

### cypher-shell Not Found
Install Neo4j or just cypher-shell:
```bash
# macOS
brew install neo4j

# Ubuntu/Debian
wget -O - https://debian.neo4j.com/neotechnology.gpg.key | sudo apt-key add -
echo 'deb https://debian.neo4j.com stable latest' | sudo tee /etc/apt/sources.list.d/neo4j.list
sudo apt-get update
sudo apt-get install cypher-shell
```

---

## 4. Metrics Test Script

**File:** `test_metrics.sh`

A simple script to test and demonstrate the Prometheus metrics endpoint.

**Usage:**
```bash
# Make executable
chmod +x test_metrics.sh

# Run the test
./test_metrics.sh
```

**What it does:**
- Checks metrics endpoint availability
- Shows current connection and health check metrics
- Generates test traffic (health checks)
- Displays updated metrics
- Lists all available metric categories

**Output example:**
```
✓ Metrics endpoint is accessible
Current connection metrics:
  bolt_proxy_connections_total 5
  bolt_proxy_active_connections 0

Generating test traffic...
Updated health check metrics:
  bolt_proxy_health_checks_total{status="success"} 10
```

**Direct access:**
```bash
# View all metrics
curl http://localhost:9090/metrics

# View only bolt-proxy metrics
curl http://localhost:9090/metrics | grep "bolt_proxy_"
```

---

## Additional Examples

### Using with mgconsole (Memgraph CLI)

```bash
# Connect through bolt-proxy
mgconsole --host localhost --port 8080 --use-ssl=False
```

### Using with Neo4j Browser

While Neo4j Browser doesn't directly support custom ports, you can:
1. Access Neo4j directly at http://localhost:7474
2. Change connection to `bolt://localhost:8080`
3. Enter credentials if required

---

## Development

To add your own test scripts:
1. Follow the pattern in existing scripts
2. Test basic connectivity first
3. Clean up any test data created
4. Update this README with usage instructions
