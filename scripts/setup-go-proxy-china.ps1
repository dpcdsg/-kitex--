# 国内网络下 Go 拉依赖常连不上 proxy.golang.org，可执行本脚本写入用户级 GOPROXY。
# 用法：powershell -ExecutionPolicy Bypass -File .\scripts\setup-go-proxy-china.ps1

$ErrorActionPreference = "Stop"

$machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
$userPath2 = [Environment]::GetEnvironmentVariable("Path", "User")
$env:Path = "$machinePath;$userPath2"

& go env -w GOPROXY=https://goproxy.cn,direct
& go env -w GOSUMDB=sum.golang.google.cn

Write-Host "GOPROXY / GOSUMDB 已设置。当前："
& go env GOPROXY
& go env GOSUMDB
