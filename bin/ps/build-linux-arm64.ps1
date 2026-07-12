# Build wol_admin for Linux arm64
param([string]$Version)
$ErrorActionPreference = "Stop"
Set-Location (Join-Path (Split-Path -Parent $MyInvocation.MyCommand.Path) "..\..")

if (-not $Version) {
    $Version = (Select-String 'Version\s*=\s*"([^"]+)"' "version\version.go").Matches.Groups[1].Value
}
$BuildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-dd HH:mm:ss")

Write-Host "Building frontend..." -ForegroundColor Cyan
Set-Location "frontend"; npm run build; if ($LASTEXITCODE -ne 0) { exit 1 }; Set-Location ".."

Write-Host "Building wol_admin v$Version linux/arm64" -ForegroundColor Cyan
$env:CGO_ENABLED = "0"; $env:GOOS = "linux"; $env:GOARCH = "arm64"
go build -ldflags "-s -w -X wol_admin/version.Version=$Version -X wol_admin/version.Arch=arm64 -X 'wol_admin/version.BuildTime=$BuildTime'" -o build\wol_admin .
Remove-Item Env:\CGO_ENABLED, Env:\GOOS, Env:\GOARCH -ErrorAction SilentlyContinue
