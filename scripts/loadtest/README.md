# k6 Load Test

## 0) 前置条件

- **后端已启动**：至少 API 网关可访问（默认示例 `http://127.0.0.1:8080`，以你实际配置为准）。
- **Go 编译/测试**：若 `go test` 报 `proxy.golang.org` 超时，在项目根目录执行：

  ```powershell
  powershell -ExecutionPolicy Bypass -File .\scripts\setup-go-proxy-china.ps1
  ```

- **k6**：需本机已安装 [k6](https://k6.io/docs/get-started/installation/)（Windows 可用 `winget install GrafanaLabs.k6`；若网络失败可官网下载 MSI/ZIP）。

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
- `api_business_requests_total{handler,result,error_code}` — 业务成功/失败（与 JSON `status_code` 一致；HTTP 常为 200）

Prometheus 已配置抓取 `job=api`（默认 `host.docker.internal:8080`，见 `config/prometheus/prometheus.yml`）。

Grafana：**Dashboards → Import → Import via dashboard JSON model**，粘贴或上传 `config/grafana/dashboards/tiktok-api-seckill.import.json`（会提示选择 Prometheus）；若数据源 UID 固定为 `prometheus`，也可用 `tiktok-api-seckill.json`。Docker 挂载 `./config/grafana` 时会自动加载后者。
