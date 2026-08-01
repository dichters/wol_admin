# Build & package wol-panel for Windows x86-64
param([string]$Version)
$ErrorActionPreference = "Stop"
Set-Location (Join-Path (Split-Path -Parent $MyInvocation.MyCommand.Path) "..\..")

if (-not $Version) {
    $Version = (Select-String 'Version\s*=\s*"([^"]+)"' "version\version.go").Matches.Groups[1].Value
}
$BuildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-dd HH:mm:ss")

if (-not (Test-Path "dist")) {
    Write-Host "Building frontend..." -ForegroundColor Cyan
    Set-Location "frontend"; npm run build; if ($LASTEXITCODE -ne 0) { exit 1 }; Set-Location ".."
}

Write-Host "Building wol-panel v$Version windows/amd64" -ForegroundColor Cyan
$env:CGO_ENABLED = "0"; $env:GOOS = "windows"; $env:GOARCH = "amd64"
go build -ldflags "-s -w -X wol-panel/version.Version=$Version -X wol-panel/version.Arch=amd64 -X 'wol-panel/version.BuildTime=$BuildTime'" -o build\wol-panel.exe .
Remove-Item Env:\CGO_ENABLED, Env:\GOOS, Env:\GOARCH -ErrorAction SilentlyContinue

$ReleaseDir = "release"
New-Item -ItemType Directory -Force -Path $ReleaseDir | Out-Null

$PkgName = "wol-panel-$Version-windows-x86-64"
$PkgDir = Join-Path $ReleaseDir $PkgName
if (Test-Path $PkgDir) { Remove-Item -Recurse -Force $PkgDir }
New-Item -ItemType Directory -Force -Path $PkgDir | Out-Null

Copy-Item "build\wol-panel.exe" $PkgDir
Copy-Item "config.template.json" $PkgDir

$ZipPath = Join-Path $ReleaseDir "$PkgName.zip"
if (Test-Path $ZipPath) { Remove-Item -Force $ZipPath }
Compress-Archive -Path (Join-Path $PkgName "*") -DestinationPath $ZipPath
Remove-Item -Recurse -Force $PkgDir

Write-Host "Done: $ZipPath" -ForegroundColor Green
