# Build wol_admin for Windows arm
param([string]$Version)
$ErrorActionPreference = "Stop"
Set-Location (Join-Path (Split-Path -Parent $MyInvocation.MyCommand.Path) "..\..")

if (-not $Version) {
    $Version = (Select-String 'Version\s*=\s*"([^"]+)"' "version\version.go").Matches.Groups[1].Value
}
$BuildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

Write-Host "Building frontend..." -ForegroundColor Cyan
Set-Location "frontend"; npm run build; if ($LASTEXITCODE -ne 0) { exit 1 }; Set-Location ".."

Write-Host "Building wol_admin v$Version windows/arm" -ForegroundColor Cyan
$env:CGO_ENABLED = "0"; $env:GOOS = "windows"; $env:GOARCH = "arm"
go build -ldflags "-s -w -X wol_admin/version.Version=$Version -X wol_admin/version.Arch=arm -X wol_admin/version.BuildTime=$BuildTime" -o build\wol_admin.exe .
Remove-Item Env:\CGO_ENABLED, Env:\GOOS, Env:\GOARCH -ErrorAction SilentlyContinue
