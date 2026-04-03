#Requires -RunAsAdministrator
<#
.SYNOPSIS
  将 Docker Desktop 使用的 WSL2 发行版迁移到 E: 盘（释放 C: 空间）。

.DESCRIPTION
  当前环境检测到仅有发行版 "docker-desktop"（部分新版本不再单独列出 docker-desktop-data）。
  流程：退出 Docker → wsl --shutdown → export 到 E: → unregister → import 到 E:\WSL\docker-desktop

  使用前请关闭所有容器内工作。导出耗时与镜像/卷大小成正比，tar 会落在 E:\WSL\backup。

.EXAMPLE
  右键「以管理员身份运行 PowerShell」，执行：
  Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass -Force
  & "E:\360Downloads\Software\login_system\-kitex--\scripts\migrate-docker-wsl-to-e.ps1"
#>

$ErrorActionPreference = "Stop"

$Root = "E:\WSL"
$ImportDir = Join-Path $Root "docker-desktop"
$BackupDir = Join-Path $Root "backup"
$Tar = Join-Path $BackupDir "docker-desktop.tar"
$Distro = "docker-desktop"

function Test-Admin {
    $id = [Security.Principal.WindowsIdentity]::GetCurrent()
    $p = New-Object Security.Principal.WindowsPrincipal($id)
    return $p.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

if (-not (Test-Admin)) {
    Write-Error "请右键 PowerShell -> 以管理员身份运行，再执行本脚本。"
}

$freeE = (Get-PSDrive E).Free / 1GB
if ($freeE -lt 20) {
    Write-Warning "E: 剩余空间约 $([math]::Round($freeE,2)) GB，若 Docker 数据很大可能不够，建议先清理 E: 或扩容。"
}

New-Item -ItemType Directory -Force -Path $ImportDir, $BackupDir | Out-Null

Write-Host "== 1/5 结束 Docker Desktop 进程 =="
Get-Process "Docker Desktop" -ErrorAction SilentlyContinue | Stop-Process -Force
Get-Process "com.docker.backend" -ErrorAction SilentlyContinue | Stop-Process -Force
Start-Sleep -Seconds 4

Write-Host "== 2/5 wsl --shutdown =="
wsl.exe --shutdown
Start-Sleep -Seconds 3

$exists = & wsl.exe -l -q 2>$null
if ($exists -notcontains $Distro) {
    Write-Error "未找到 WSL 发行版 $Distro。请执行 wsl -l -v 确认名称后修改脚本中的 `$Distro。"
}

if (Test-Path $Tar) {
    Remove-Item $Tar -Force
}

Write-Host "== 3/5 导出 $Distro 到 $Tar （可能数十分钟，请勿中断）=="
wsl.exe --export $Distro $Tar
if (-not (Test-Path $Tar)) {
    Write-Error "导出失败：未生成 $Tar"
}
$gb = [math]::Round((Get-Item $Tar).Length / 1GB, 2)
Write-Host "导出完成，约 $gb GB"

Write-Host "== 4/5 注销原 C: 侧发行版（数据已在 tar 中）=="
wsl.exe --unregister $Distro

Write-Host "== 5/5 导入到 $ImportDir =="
if (Test-Path $ImportDir) {
    $left = Get-ChildItem $ImportDir -Force -ErrorAction SilentlyContinue
    if ($left.Count -gt 0) {
        Write-Error "目录非空: $ImportDir ，请手动清空后重试。"
    }
}
wsl.exe --import $Distro $ImportDir $Tar --version 2

Write-Host "== 删除临时 tar（可选，已导入成功）=="
Remove-Item $Tar -Force -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "迁移完成。请从开始菜单启动 Docker Desktop，等待引擎就绪后执行: docker run hello-world"
Write-Host "若 Docker 无法启动，可把本窗口完整输出与 Docker 诊断日志发给维护人员。"
