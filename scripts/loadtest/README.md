# k6 Load Test

## 0) 与仓库内 Docker Compose 对齐

本项目 **基础设施** 在根目录 `docker-compose.yml`，**业务四服务** 在 `docker-compose.apps.yml`（与主文件叠加使用）。

| 组件 | 宿主机端口（默认） | 说明 |
|------|------------------|------|
| **API 网关** `tiktok-api` | **10001** | 压测、浏览器访问 HTTP 入口 |
| `tiktok-user` | 10002 | Kitex |
| `tiktok-product` | 10003 | Kitex |
| `tiktok-order` | 10004 | Kitex |
| MySQL | **43306** → 容器 3306 | 与 `config.docker.yaml` 中 `mysql:3306`（容器内）对应 |
| Redis | 6379 | 密码见 compose / `config.docker.yaml` |
| RabbitMQ | 5672 / 管理台 15672 | |
| etcd | 2379 | 业务容器 `ETCD_ADDR=etcd:2379` |
| Prometheus | 9090 | 已配抓取 `tiktok-api:10001/metrics` |
| Grafana | 3000 | |

**启动顺序（项目根目录）：**

1. 镜像来自 `docker.xuanyuan.run` 时需先：`docker login docker.xuanyuan.run`（见 `docker-compose.yml` 顶部注释）。
2. 构建业务镜像（仅首次或代码变更后）。推荐用脚本顺带 **清悬空镜像、压缩 C 盘 `docker_data.vhdx` 再 build**（会关 Docker 并可能弹 UAC）：

   ```powershell
   .\scripts\redeploy-docker.ps1
   ```

   仅快速 build、不压盘：`.\scripts\redeploy-docker.ps1 -Quick` 或 `docker build -t tiktok:latest .`

3. 先起基础设施，再起应用（推荐一条命令叠加两个文件）：

   ```bash
   docker compose -f docker-compose.yml -f docker-compose.apps.yml up -d
   ```

   `docker-compose.apps.yml` 会把 `config/config.docker.yaml` 挂进 etcd，供 `etcd-monitor` 写入 **容器网络版** 配置；业务容器通过服务名 `mysql`、`redis`、`etcd` 等互联。

4. 确认 API 可用：`curl -s http://127.0.0.1:10001/metrics | head`

**压测默认地址**：脚本里 `BASE_URL` 默认为 `http://127.0.0.1:10001`（与 `tiktok-api` 端口映射一致）。若你在本机直接 `go run` API 且监听 8080，再设 `BASE_URL=http://127.0.0.1:8080`。

---

## 1) 压测怎么跑（推荐：Docker，无需安装 k6）

**项目根目录**执行（API 为全栈 Docker、映射宿主机 `10001` 时，`BASE_URL` 在容器内须用 `host.docker.internal`）：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\loadtest\run-k6-docker.ps1
```

可选参数示例：

```powershell
.\scripts\loadtest\run-k6-docker.ps1 -Vus 100 -Duration 60s -ActivityId 1000001
```

等价的手动 `docker run`（需已能拉取 `docker.xuanyuan.run/grafana/k6:latest`）：

```powershell
docker run --rm -i `
  -e BASE_URL=http://host.docker.internal:10001 `
  -e VUS=50 -e DURATION=30s `
  -v "${PWD}/scripts/loadtest:/scripts:ro" `
  docker.xuanyuan.run/grafana/k6:latest run /scripts/seckill_k6.js
```

若 API 跑在**宿主机 8080**（非容器 10001），把上面 `BASE_URL` 改为 `http://host.docker.internal:8080`。

**说明**：`seckill_k6.js` 默认 `BASE_URL=http://127.0.0.1:10001` 是给**本机已安装 k6** 时用的；用 Docker 跑 k6 时不要依赖该默认值，须显式传 `host.docker.internal`（见脚本 `run-k6-docker.ps1`）。

### 全接口巡检：`all_api_k6.js`

对 **当前已挂载的 Hertz 路由** 各请求至少一次（含 `GET /metrics`）。**未包含**：handler 里有、但未在 `router`/`main` 注册的路径（例如 `POST /seckill/product/publish/`）。

```powershell
.\scripts\loadtest\run-k6-docker.ps1 -K6Script all_api_k6.js -Vus 3 -Duration 60s
```

可选：`PRODUCT_ID`、`ORDER_NO` 传给脚本（见 `run-k6-docker.ps1` 的 `-ProductId` / `-OrderNo`）。脚本会在每轮迭代内 **创建商品与秒杀活动再软删**，并发不要过大，以免拖慢库或触发 Sentinel。

本机 k6：

```powershell
$env:BASE_URL="http://127.0.0.1:10001"
k6 run scripts/loadtest/all_api_k6.js
```

### 压测主要在测什么（和「接口测一遍」的区别）

| 侧重 | 说明 |
|------|------|
| **HTTP 与网关** | Hertz 路由、中间件（Recovery、Gzip、**Prometheus 埋点**、**Sentinel 限流**）在并发下是否稳定；是否出现连接错误、非预期 5xx。 |
| **整条调用链** | API → Kitex → user/product/order → MySQL/Redis 等在压力下的延迟与错误率（很多业务错误仍是 **HTTP 200 + JSON `status_code!=0`**）。 |
| **可观测性** | 配合 Prometheus/Grafana 看 QPS、延迟、业务计数器 `api_business_requests_total` 等是否随压测变化。 |

**不保证**：每个接口在业务语义上都「成功」（例如非卖家创建商品、库存不足）；全接口脚本默认 **只校验 HTTP 状态为 200**，不把 JSON 业务码当作失败，除非你自己改 `check`。

**秒杀专项脚本** `seckill_k6.js` 更贴近「列表 + 秒杀下单」真实流量形状；**全接口脚本**更贴近「路由冒烟 + 链路透传」与回归。

---

## 2) 其他前置

- **Go 编译/测试**（宿主机）：若 `go test` 访问 `proxy.golang.org` 超时：

  ```powershell
  powershell -ExecutionPolicy Bypass -File .\scripts\setup-go-proxy-china.ps1
  ```

- **压测账号与活动**：脚本会在 setup 里先注册再登录；`ACTIVITY_ID` 等需与库里活动一致（默认 `1000001` 为示例）。

---

## 3) 本机已安装 k6 时（可选）

```bash
k6 run scripts/loadtest/seckill_k6.js
# 或
k6 run scripts/loadtest/all_api_k6.js
```

```powershell
$env:BASE_URL="http://127.0.0.1:10001"
k6 run scripts/loadtest/seckill_k6.js
k6 run scripts/loadtest/all_api_k6.js
```

自定义变量（bash）：

```bash
BASE_URL=http://127.0.0.1:10001 \
USERNAME=dpc PASSWORD=123 ACTIVITY_ID=1000001 VUS=100 DURATION=60s \
k6 run scripts/loadtest/seckill_k6.js
```

---

## 4) Prometheus / Grafana

API 暴露：

```text
GET http://127.0.0.1:10001/metrics
```

`config/prometheus/prometheus.yml` 里 `job_name: api` 已包含：

- `tiktok-api:10001` — 与 Prometheus 同 Docker 网络时抓取容器内 API（**全栈 Docker 推荐**）
- `host.docker.internal:8080` — 宿主机本地跑 API 时用

指标示例：

- `api_http_requests_total{method,path,code}`
- `api_http_request_duration_seconds{method,path,code}`
- `api_http_requests_in_flight{method,path}`
- `api_business_requests_total{handler,result,error_code}`（各 handler 在调用 `RecordBusinessResult` 时递增）

Grafana：导入 `config/grafana/dashboards/tiktok-api-seckill.import.json` 或 `tiktok-api-seckill.json`，数据源选 Prometheus。
