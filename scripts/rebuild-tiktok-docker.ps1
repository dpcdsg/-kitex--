# 代码（含埋点）改完后：重新构建 tiktok 镜像并重启业务容器，使新二进制生效。
# 前置：已安装 Docker Desktop；基础设施已起（docker compose up -d）。
#
# 仅 API：
#   .\scripts\rebuild-tiktok-docker.ps1
# 四个服务一起重建（api / user / product / order）：
#   .\scripts\rebuild-tiktok-docker.ps1 -AllServices
#
# 之后建议重启 Prometheus 以重读配置（若改过 prometheus.yml）：
#   docker compose restart prometheus

param(
    [switch]$AllServices
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $root

Write-Host "[rebuild-tiktok-docker] Project root: $root"

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Error "未找到 docker 命令，请先安装 Docker Desktop 并确保在 PATH 中。"
}

Write-Host "[rebuild-tiktok-docker] docker build -t tiktok:latest ..."
docker build -t tiktok:latest .

$composeBase = @("-f", "docker-compose.yml", "-f", "docker-compose.apps.yml")

if ($AllServices) {
    Write-Host "[rebuild-tiktok-docker] Recreating tiktok-api, tiktok-user, tiktok-product, tiktok-order ..."
    docker compose @composeBase up -d --force-recreate tiktok-api tiktok-user tiktok-product tiktok-order
}
else {
    Write-Host "[rebuild-tiktok-docker] Recreating tiktok-api only ..."
    docker compose @composeBase up -d --force-recreate tiktok-api
}

Write-Host "[rebuild-tiktok-docker] Done."
Write-Host "  - 容器内 API 端口映射: http://localhost:10001/metrics"
Write-Host "  - 若 Prometheus 改过 prometheus.yml: docker compose restart prometheus"
Write-Host "  - 若 API 跑在宿主机 8080 而非容器: 本机执行 go build 后重启 API 进程即可，不必 docker build。"
