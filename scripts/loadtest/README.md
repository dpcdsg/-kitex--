# k6 Load Test

## 1) Run

```bash
k6 run scripts/loadtest/seckill_k6.js
```

## 2) Custom arguments

```bash
BASE_URL=http://127.0.0.1:8080 \
USERNAME=dpc \
PASSWORD=123 \
ACTIVITY_ID=1000001 \
VUS=100 \
DURATION=60s \
k6 run scripts/loadtest/seckill_k6.js
```

## 3) Prometheus metrics endpoint

After starting API service, scrape:

```text
GET /metrics
```

Core custom metrics:

- `api_http_requests_total{method,path,code}`
- `api_http_request_duration_seconds{method,path,code}`
- `api_http_requests_in_flight{method,path}`
