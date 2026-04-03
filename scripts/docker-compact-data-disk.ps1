# 压缩 Docker Desktop 在 C 盘上的 WSL 数据盘（释放 NTFS 占用）。
# 须管理员运行；Docker 须已退出，建议先在本机执行 wsl --shutdown。
#
# 由 redeploy-docker.ps1 以「提升权限」方式调用，并传入当前用户的 VHD 绝对路径。
#
# 手动执行（管理员 PowerShell）示例：
#   wsl --shutdown
#   powershell -ExecutionPolicy Bypass -File .\scripts\docker-compact-data-disk.ps1 `
#     -DataVhd "$env:LOCALAPPDATA\Docker\wsl\disk\docker_data.vhdx"

param(
    [string]$DataVhd = "",
    [string]$MainVhd = ""
)

$ErrorActionPreference = "Stop"

$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole(
    [Security.Principal.WindowsBuiltinRole]::Administrator)
if (-not $isAdmin) {
    Write-Error "请以管理员身份运行本脚本（压缩 VHD 需要 Hyper-V Optimize-VHD）。"
}

if (-not (Get-Command Optimize-VHD -ErrorAction SilentlyContinue)) {
    Write-Error "未找到 Optimize-VHD。请启用 Hyper-V PowerShell 或相关 Windows 功能。"
}

if (-not $DataVhd) {
    $DataVhd = Join-Path $env:LOCALAPPDATA "Docker\wsl\disk\docker_data.vhdx"
}
if (-not $MainVhd) {
    $MainVhd = Join-Path $env:LOCALAPPDATA "Docker\wsl\main\ext4.vhdx"
}

Write-Host "[docker-compact] DataVhd=$DataVhd"
if (Test-Path $DataVhd) {
    $before = (Get-Item $DataVhd).Length / 1GB
    Write-Host "[docker-compact] docker_data.vhdx before: $([math]::Round($before, 2)) GB"
    Optimize-VHD -Path $DataVhd -Mode Full
    $after = (Get-Item $DataVhd).Length / 1GB
    Write-Host "[docker-compact] docker_data.vhdx after:  $([math]::Round($after, 2)) GB"
}
else {
    Write-Warning "[docker-compact] 未找到 $DataVhd，跳过（若 Docker 数据盘路径不同请传 -DataVhd）。"
}

if (Test-Path $MainVhd) {
    Write-Host "[docker-compact] compact main ext4.vhdx ..."
    Optimize-VHD -Path $MainVhd -Mode Full
}

Write-Host "[docker-compact] Done."
