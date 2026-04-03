# 删除「悬空」镜像（无标签、多为旧 build 层），保留所有仍带 tag 的镜像（含当前 tiktok:latest）。
# 在项目根目录执行: powershell -ExecutionPolicy Bypass -File .\scripts\docker-prune-dangling.ps1
#
# 可选: -WithBuilder  同时清理 docker build 缓存
# 可选: -DryRun       只打印将执行的命令

param(
    [switch]$WithBuilder,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

function Run($cmd) {
    Write-Host ">> $cmd" -ForegroundColor Cyan
    if (-not $DryRun) {
        Invoke-Expression $cmd
    }
}

Write-Host "=== docker image prune (dangling only) ===" -ForegroundColor Yellow
Write-Host "Removes untagged layers from previous builds; keeps tagged images like tiktok:latest."
Run "docker image prune -f"

if ($WithBuilder) {
    Write-Host "`n=== docker builder prune ===" -ForegroundColor Yellow
    Run "docker builder prune -f"
}

Write-Host "`nDone. Show disk use:" -ForegroundColor Green
if (-not $DryRun) {
    docker system df
}
