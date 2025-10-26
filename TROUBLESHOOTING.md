# Troubleshooting Guide

## Common Issues

### Error: "failed to read expected Hello from client <nil> false"

**Cause:** The client connection closed after version negotiation but before sending the HELLO message.

**Possible Reasons:**

1. **TLS/SSL Mismatch**
   - Client expects encrypted connection but bolt-proxy is not configured with TLS
   - Or vice versa: proxy has TLS but client doesn't

   **Solution:**
   ```bash
   # For non-TLS connections, ensure client uses bolt:// not bolt+s:// or bolt+ssc://
   # Python example:
   driver = GraphDatabase.driver("bolt://localhost:8080", auth=("user", "pass"))

   # If proxy needs TLS, configure cert/key:
   export BOLT_PROXY_CERT=/path/to/cert.pem
   export BOLT_PROXY_KEY=/path/to/key.pem
   ```

2. **Client Timeout**
   - Client timeout is too short
   - Network latency between client and proxy

   **Solution:**
   - Increase client connection timeout
   - Check network connectivity
   ```bash
   # Test basic connectivity
   nc -zv localhost 8080

   # Test bolt handshake
   echo -ne '\x60\x60\xb0\x17\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00' | nc localhost 8080 | xxd
   ```

3. **Client Library Version Incompatibility**
   - Very old or very new client drivers may not be compatible

   **Solution:**
   - Use recommended driver versions:
     - Python: neo4j >= 5.0.0
     - Go: neo4j-go-driver/v4 >= 4.4.0
     - Java: neo4j-java-driver >= 4.3.0

4. **Backend Server Not Ready**
   - Neo4j or Memgraph backend is not fully started
   - Backend is refusing connections

   **Solution:**
   ```bash
   # Check backend is running
   docker-compose ps

   # Check backend logs
   docker-compose logs neo4j
   # or
   docker-compose logs memgraph

   # Test direct connection to backend
   # For Neo4j:
   cypher-shell -a bolt://localhost:7687 -u neo4j -p password

   # For Memgraph:
   mgconsole --host localhost --port 7687
   ```

5. **Health Check Probes**
   - Load balancers or monitoring tools sending probes
   - This is usually harmless and can be ignored if occasional

   **Solution:**
   - These show up in DEBUG logs but don't affect real clients
   - Consider filtering these logs if too verbose

6. **Authentication Method Mismatch**
   - Proxy configured for auth but client not sending credentials
   - Auth provider is down or misconfigured

   **Solution:**
   ```bash
   # Check AUTH_METHOD setting
   echo $AUTH_METHOD

   # If using BASIC_AUTH, verify URL is reachable:
   curl -v $BASIC_AUTH_URL

   # If using AAD_TOKEN_AUTH, verify provider:
   curl -v $AAD_TOKEN_PROVIDER
   ```

## Debug Steps

### 1. Enable Debug Logging

```bash
# In docker-compose.yml or docker-compose.neo4j.yml
environment:
    BOLT_PROXY_DEBUG: "1"

# Or when running directly:
export BOLT_PROXY_DEBUG=1
./bolt-proxy
```

### 2. Check Connection Flow

The typical Bolt connection flow is:
1. Client sends Bolt signature: `0x6060b017`
2. Client sends version handshake (16 bytes)
3. Proxy responds with selected version (4 bytes)
4. Client sends HELLO message
5. Proxy authenticates (if enabled)
6. Proxy forwards to backend

If failing at step 4, the error "failed to read expected Hello" appears.

### 3. Test with Minimal Client

Create a simple test to isolate the issue:

```python
# test_simple.py
from neo4j import GraphDatabase
import sys

uri = "bolt://localhost:8080"
try:
    driver = GraphDatabase.driver(uri, auth=("", ""))
    with driver.session() as session:
        result = session.run("RETURN 1")
        print("✓ Success:", result.single()[0])
    driver.close()
except Exception as e:
    print("✗ Failed:", e)
    sys.exit(1)
```

### 4. Capture Network Traffic

```bash
# Capture traffic on proxy port
sudo tcpdump -i lo0 -A -s 0 'tcp port 8080'

# Or use Wireshark for detailed analysis
```

### 5. Check Docker Networking

```bash
# Verify services can communicate
docker-compose exec bolt-proxy ping neo4j
docker-compose exec bolt-proxy nc -zv neo4j 7687

# Check port bindings
docker-compose ps
netstat -an | grep 8080
```

## Testing

### Quick Connection Test

```bash
# Python quick test
python -c "from neo4j import GraphDatabase; d = GraphDatabase.driver('bolt://localhost:8080', auth=('','')); print(d.verify_connectivity()); d.close()"

# Using cypher-shell through proxy
cypher-shell -a bolt://localhost:8080

# Using examples
python examples/test_bolt_proxy.py --host localhost --port 8080
```

### Verify Proxy Configuration

```bash
# Check environment variables
docker-compose exec bolt-proxy env | grep BOLT_PROXY

# Verify backend URI is correct
# Should point to backend service, not localhost
# ✓ BOLT_PROXY_URI=bolt://neo4j:7687
# ✗ BOLT_PROXY_URI=bolt://localhost:7687 (wrong in container)
```

## Common Fixes

### Fix 1: Restart Services in Correct Order

```bash
# Stop everything
docker-compose down

# Start backend first, wait for it to be ready
docker-compose up -d neo4j
sleep 30

# Then start proxy
docker-compose up -d bolt-proxy

# Check logs
docker-compose logs -f
```

### Fix 2: Disable Authentication for Testing

```bash
# In docker-compose.yml, remove or comment out AUTH_METHOD
# environment:
#     # AUTH_METHOD: BASIC_AUTH

# Restart
docker-compose restart bolt-proxy
```

### Fix 3: Update Client Libraries

```bash
# Python
pip install --upgrade neo4j

# Go
go get -u github.com/neo4j/neo4j-go-driver/v4@latest
```

## Still Having Issues?

1. Check GitHub issues: https://github.com/memgraph/bolt-proxy/issues
2. Provide debug logs when reporting:
   ```bash
   docker-compose logs bolt-proxy > debug.log
   ```
3. Include:
   - Client library and version
   - Backend (Neo4j/Memgraph) version
   - Connection code snippet
   - Full error message
