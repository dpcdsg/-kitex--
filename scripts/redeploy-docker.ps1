# 代码更新后：清理旧 Docker 层/缓存 → 压缩 C 盘 Docker 数据盘 → 重新 build 并拉起业务容器。
# 在项目根目录的上一级为仓库根：本脚本位于 scripts\，自动 cd 到仓库根。
#
# 默认行为（推荐，避免 C 盘被 docker_data.vhdx 撑满）：
#   - 删除悬空镜像（旧 build 遗留）
#   - 可选深度清理 build 缓存（-DeepClean）
#   - 退出 Docker Desktop、wsl --shutdown 后压缩 docker_data.vhdx（会弹出 UAC，请允许）
#   - 再启动 Docker、docker build、compose --force-recreate
#
# 用法：
#   .\scripts\redeploy-docker.ps1
#   .\scripts\redeploy-docker.ps1 -AllServices
#   .\scripts\redeploy-docker.ps1 -Quick              # 仅 build+recreate，不清理不压缩
#   .\scripts\redeploy-docker.ps1 -SkipCompact       # 清理镜像/缓存但不压缩 VHD（无需 UAC）
#   .\scripts\redeploy-docker.ps1 -DeepClean         # docker builder prune -a（下次 build 更慢）

param(
    [switch]$AllServices,
    [switch]$Quick,
    [switch]$SkipPrune,
    [switch]$SkipCompact,
    [switch]$DeepClean
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $root

$composeBase = @("-f", "docker-compose.yml", "-f", "docker-compose.apps.yml")

function Stop-DockerDesktopQuietly {
    Get-Process -Name "Docker Desktop" -ErrorAction SilentlyContinue | Stop-Process -Force
    Start-Sleep -Seconds 2
}

function Start-DockerDesktopAndWait {
    $exe = "${env:ProgramFiles}\Docker\Docker\Docker Desktop.exe"
    if (-not (Test-Path $exe)) {
        Write-Error "未找到 Docker Desktop：$exe"
    }
    Write-Host "[redeploy-docker] 正在启动 Docker Desktop ..."
    Start-Process -FilePath $exe
    $deadline = (Get-Date).AddMinutes(4)
    while ((Get-Date) -lt $deadline) {
        try {
            docker info 2>$null | Out-Null
            if ($LASTEXITCODE -eq 0) {
                Write-Host "[redeploy-docker] Docker 已就绪。"
                return
            }
        }
        catch { }
        Start-Sleep -Seconds 3
    }
    Write-Error "Docker 在 4 分钟内未就绪，请手动打开 Docker Desktop 后重试。"
}

function Test-Administrator {
    $p = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    return $p.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator)
}

Write-Host "[redeploy-docker] Project root: $root"

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Error "未找到 docker 命令，请先安装 Docker Desktop。"
}

if ($Quick) {
    $SkipPrune = $true
    $SkipCompact = $true
}

# --- 1) 清理（Docker 在跑时执行）---
if (-not $SkipPrune) {
    Write-Host "[redeploy-docker] docker image prune -f（悬空镜像）..."
    docker image prune -f
    if ($DeepClean) {
        Write-Host "[redeploy-docker] docker builder prune -a -f（全部 build 缓存，下次 build 较慢）..."
        docker builder prune -a -f
    }
    else {
        Write-Host "[redeploy-docker] docker builder prune -f（未使用的 build 缓存）..."
        docker builder prune -f
    }
}

# --- 2) 压缩 C 盘 VHD ---
if (-not $SkipCompact) {
    $dataVhd = Join-Path $env:LOCALAPPDATA "Docker\wsl\disk\docker_data.vhdx"
    $mainVhd = Join-Path $env:LOCALAPPDATA "Docker\wsl\main\ext4.vhdx"
    $compactScript = (Resolve-Path (Join-Path $PSScriptRoot "docker-compact-data-disk.ps1")).Path

    if (-not (Test-Path $dataVhd)) {
        Write-Warning "[redeploy-docker] 未找到 $dataVhd，跳过压缩（若 Docker 安装路径不同可改用 -SkipCompact）。"
    }
    else {
        Write-Host "[redeploy-docker] 关闭 Docker Desktop 并 wsl --shutdown，以便压缩 VHD ..."
        Stop-DockerDesktopQuietly
        wsl --shutdown 2>$null
        Start-Sleep -Seconds 2

        $dataVhdFull = (Resolve-Path $dataVhd).Path

        if (Test-Administrator) {
            Write-Host "[redeploy-docker] 当前已管理员，直接压缩 ..."
            if (Test-Path $mainVhd) {
                & $compactScript -DataVhd $dataVhdFull -MainVhd (Resolve-Path $mainVhd).Path
            }
            else {
                & $compactScript -DataVhd $dataVhdFull
            }
        }
        else {
            Write-Host "[redeploy-docker] 将弹出 UAC：请允许管理员窗口执行磁盘压缩（仅 Optimize-VHD）。"
            $argList = @(
                "-NoProfile", "-ExecutionPolicy", "Bypass",
                "-File", $compactScript,
                "-DataVhd", $dataVhdFull
            )
            if (Test-Path $mainVhd) {
                $argList += @("-MainVhd", (Resolve-Path $mainVhd).Path)
            }
            $p = Start-Process -FilePath "powershell.exe" -ArgumentList $argList -Verb RunAs -Wait -PassThru
            if ($p.ExitCode -ne 0) {
                Write-Warning "[redeploy-docker] 压缩脚本退出码 $($p.ExitCode)。若你点了「否」关闭 UAC，可改用 -SkipCompact 仅部署。"
            }
        }

        Start-DockerDesktopAndWait
    }
}

# --- 3) build + compose ---
Write-Host "[redeploy-docker] docker build -t tiktok:latest ..."
docker build -t tiktok:latest .

if ($AllServices) {
    Write-Host "[redeploy-docker] Recreating tiktok-api, tiktok-user, tiktok-product, tiktok-order ..."
    docker compose @composeBase up -d --force-recreate tiktok-api tiktok-user tiktok-product tiktok-order
}
else {
    Write-Host "[redeploy-docker] Recreating tiktok-api only ..."
    docker compose @composeBase up -d --force-recreate tiktok-api
}

Write-Host "[redeploy-docker] Done."
Write-Host "  - API: http://localhost:10001  /metrics"
Write-Host "  - 改过 prometheus.yml 可执行: docker compose -f docker-compose.yml -f docker-compose.apps.yml restart prometheus"
Write-Host "  - 快速仅重建（不清理 C 盘）: .\scripts\redeploy-docker.ps1 -Quick"
