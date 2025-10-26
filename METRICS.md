# Prometheus Metrics

Bolt-proxy exposes Prometheus metrics for monitoring and observability.

## Accessing Metrics

Metrics are exposed on the `/metrics` HTTP endpoint. By default, the metrics server listens on port `9090`.

### Configuration

Configure the metrics port using:

**Environment Variable:**
```bash
export BOLT_PROXY_METRICS_PORT=9090
```

**Command Line Flag:**
```bash
./bolt-proxy -metrics-port 9090
```

**Docker Compose:**
```yaml
environment:
  BOLT_PROXY_METRICS_PORT: "9090"
ports:
  - "9090:9090"
```

### Accessing Metrics Endpoint

```bash
# View all metrics
curl http://localhost:9090/metrics

# View only bolt-proxy metrics
curl http://localhost:9090/metrics | grep bolt_proxy_
```

## Available Metrics

### Connection Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `bolt_proxy_active_connections` | Gauge | Number of currently active client connections |
| `bolt_proxy_connections_total` | Counter | Total number of client connections since startup |
| `bolt_proxy_connection_duration_seconds` | Histogram | Duration of client connections in seconds |

### Authentication Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `bolt_proxy_auth_attempts_total` | Counter | Total number of authentication attempts by status (success/failure) |
| `bolt_proxy_auth_duration_seconds` | Histogram | Duration of authentication attempts in seconds |

### Message Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `bolt_proxy_messages_forwarded_total` | Counter | Total number of messages forwarded by direction and type |
| `bolt_proxy_message_bytes_total` | Counter | Total bytes of messages forwarded by direction |
| `bolt_proxy_message_processing_duration_seconds` | Histogram | Duration of message processing in seconds by type |

### Backend Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `bolt_proxy_backend_connections` | Gauge | Number of active connections to backend servers |
| `bolt_proxy_backend_connection_errors_total` | Counter | Total number of backend connection errors |
| `bolt_proxy_backend_latency_seconds` | Histogram | Latency of backend connections in seconds |

### Transaction Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `bolt_proxy_active_transactions` | Gauge | Number of currently active transactions |
| `bolt_proxy_transactions_total` | Counter | Total number of transactions by status (committed/rolled_back/failed) |

### Health Check Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `bolt_proxy_health_checks_total` | Counter | Total number of health check requests by status |

### Protocol Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `bolt_proxy_bolt_version_negotiations_total` | Counter | Total number of Bolt version negotiations by version |

### Error Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `bolt_proxy_errors_total` | Counter | Total number of errors by error_type and component |

## Example Queries

### Prometheus Queries

```promql
# Current active connections
bolt_proxy_active_connections

# Connection rate (connections per second)
rate(bolt_proxy_connections_total[5m])

# Average connection duration
rate(bolt_proxy_connection_duration_seconds_sum[5m]) / rate(bolt_proxy_connection_duration_seconds_count[5m])

# Authentication success rate
rate(bolt_proxy_auth_attempts_total{status="success"}[5m]) / rate(bolt_proxy_auth_attempts_total[5m])

# Message throughput (messages per second)
rate(bolt_proxy_messages_forwarded_total[5m])

# Backend error rate
rate(bolt_proxy_backend_connection_errors_total[5m])

# 95th percentile connection duration
histogram_quantile(0.95, rate(bolt_proxy_connection_duration_seconds_bucket[5m]))
```

## Grafana Dashboard

Here's a sample Grafana dashboard configuration:

### Connection Panel
```json
{
  "title": "Active Connections",
  "targets": [
    {
      "expr": "bolt_proxy_active_connections"
    }
  ]
}
```

### Connection Rate Panel
```json
{
  "title": "Connection Rate",
  "targets": [
    {
      "expr": "rate(bolt_proxy_connections_total[5m])"
    }
  ]
}
```

### Authentication Success Rate Panel
```json
{
  "title": "Auth Success Rate",
  "targets": [
    {
      "expr": "rate(bolt_proxy_auth_attempts_total{status='success'}[5m]) / rate(bolt_proxy_auth_attempts_total[5m]) * 100"
    }
  ]
}
```

## Integration with Prometheus

### Prometheus Configuration

Add bolt-proxy to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'bolt-proxy'
    static_configs:
      - targets: ['localhost:9090']
    scrape_interval: 15s
    scrape_timeout: 10s
```

### Docker Compose with Prometheus

```yaml
version: '3'

services:
  bolt-proxy:
    image: bolt-proxy
    ports:
      - "8080:8080"
      - "9090:9090"
    environment:
      BOLT_PROXY_METRICS_PORT: "9090"

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9091:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  prometheus_data:
  grafana_data:
```

## Alerting Rules

Example Prometheus alerting rules:

```yaml
groups:
  - name: bolt_proxy
    rules:
      - alert: HighConnectionFailureRate
        expr: rate(bolt_proxy_auth_attempts_total{status="failure"}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High connection failure rate"
          description: "Bolt-proxy authentication failure rate is {{ $value }} per second"

      - alert: NoActiveConnections
        expr: bolt_proxy_active_connections == 0
        for: 10m
        labels:
          severity: info
        annotations:
          summary: "No active connections"
          description: "Bolt-proxy has had no active connections for 10 minutes"

      - alert: HighBackendErrors
        expr: rate(bolt_proxy_backend_connection_errors_total[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High backend error rate"
          description: "Backend connection error rate is {{ $value }} per second"

      - alert: HighConnectionDuration
        expr: histogram_quantile(0.95, rate(bolt_proxy_connection_duration_seconds_bucket[5m])) > 60
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High connection duration"
          description: "95th percentile connection duration is {{ $value }} seconds"
```

## Testing Metrics

```bash
# Generate some test traffic
for i in {1..10}; do
  echo "RETURN $i" | docker exec -i bolt-proxy-memgraph-1 \
    mgconsole --host bolt-proxy --port 8080 --use-ssl=False
done

# Check metrics
curl -s http://localhost:9090/metrics | grep bolt_proxy_connections_total
```

## Troubleshooting

### Metrics endpoint not accessible

1. **Check port configuration:**
   ```bash
   docker logs bolt-proxy-bolt-proxy-1 | grep "starting metrics server"
   ```

2. **Verify port is exposed:**
   ```bash
   docker ps | grep bolt-proxy
   ```

3. **Test from inside container:**
   ```bash
   docker exec bolt-proxy-bolt-proxy-1 wget -O- http://localhost:9090/metrics
   ```

### No metrics being recorded

1. **Verify connections are reaching the proxy:**
   ```bash
   docker logs bolt-proxy-bolt-proxy-1 | tail -20
   ```

2. **Check for errors:**
   ```bash
   curl http://localhost:9090/metrics | grep -i error
   ```

### Metrics server fails to start

Check the logs for port conflicts:
```bash
docker logs bolt-proxy-bolt-proxy-1 | grep "metrics server failed"
```

Common causes:
- Port 9090 already in use
- Invalid port number configuration
- Network issues

## Best Practices

1. **Scrape Interval**: Use 15-30 second scrape intervals for bolt-proxy metrics
2. **Retention**: Keep at least 15 days of metrics data for trend analysis
3. **Alerting**: Set up alerts for high error rates and connection failures
4. **Dashboards**: Create separate dashboards for different teams (ops, dev, business)
5. **Labels**: Use Prometheus relabeling to add environment and region labels
