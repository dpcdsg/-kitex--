# 重新部署业务镜像到 Docker（默认：清悬空镜像与 build 缓存 + 压缩 C 盘 docker_data.vhdx + build + recreate）。
# 实际逻辑在 redeploy-docker.ps1；本文件为兼容旧命令的入口。
#
# 完整说明见 redeploy-docker.ps1 顶部注释。
#
# 仅快速 build（不清理、不压盘、无 UAC）：
#   .\scripts\rebuild-tiktok-docker.ps1 -Quick
#
# 四个 Kitex 容器一起重建：
#   .\scripts\rebuild-tiktok-docker.ps1 -AllServices

param(
    [switch]$AllServices,
    [switch]$Quick,
    [switch]$SkipPrune,
    [switch]$SkipCompact,
    [switch]$DeepClean
)

$ErrorActionPreference = "Stop"
& (Join-Path $PSScriptRoot "redeploy-docker.ps1") @PSBoundParameters
