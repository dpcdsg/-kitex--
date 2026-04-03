# Refresh PATH from Machine + User, then ensure Go bin is on User PATH.
# Run after installing Go (winget: GoLang.Go) or if `go` is not found in PowerShell.

$ErrorActionPreference = "Stop"

$goBin = "C:\Program Files\Go\bin"
if (-not (Test-Path "$goBin\go.exe")) {
    Write-Host "Go not found at $goBin\go.exe"
    Write-Host "Install with: winget install GoLang.Go --accept-package-agreements --accept-source-agreements"
    exit 1
}

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$goBin*") {
    $newUserPath = if ([string]::IsNullOrEmpty($userPath)) { $goBin } else { "$userPath;$goBin" }
    [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
    Write-Host "Appended to User PATH: $goBin"
} else {
    Write-Host "User PATH already contains: $goBin"
}

# Reload PATH in this session (Machine + User)
$machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
$userPath2 = [Environment]::GetEnvironmentVariable("Path", "User")
$env:Path = "$machinePath;$userPath2"

& "$goBin\go.exe" version
