# 用 Docker 跑 k6（无需本机安装 k6）。在项目根目录执行，或任意目录：
#   powershell -ExecutionPolicy Bypass -File .\scripts\loadtest\run-k6-docker.ps1
#
# 镜像与 compose 一致走轩辕源；直连 Docker Hub 失败时可用此脚本。

param(
    [string]$BaseUrl = "http://host.docker.internal:10001",
    [string]$Username = "dpc",
    [string]$Password = "123",
    [string]$ActivityId = "1000001",
    [string]$ProductId = "",
    [string]$OrderNo = "",
    [int]$Vus = 50,
    [string]$Duration = "30s",
    [ValidatePattern("^[a-zA-Z0-9][a-zA-Z0-9._-]*\.js$")]
    [string]$K6Script = "seckill_k6.js"
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$scriptDir = (Resolve-Path $PSScriptRoot).Path
$image = "docker.xuanyuan.run/grafana/k6:latest"

Write-Host "[run-k6-docker] script=$K6Script root=$root BASE_URL=$BaseUrl"

$envArgs = @(
    "-e", "BASE_URL=$BaseUrl",
    "-e", "USERNAME=$Username",
    "-e", "PASSWORD=$Password",
    "-e", "ACTIVITY_ID=$ActivityId",
    "-e", "VUS=$Vus",
    "-e", "DURATION=$Duration"
)
if ($ProductId) { $envArgs += @("-e", "PRODUCT_ID=$ProductId") }
if ($OrderNo) { $envArgs += @("-e", "ORDER_NO=$OrderNo") }

& docker run --rm -i @envArgs -v "${scriptDir}:/scripts:ro" $image run "/scripts/$K6Script"
