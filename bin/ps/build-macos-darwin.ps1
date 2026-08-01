# Build & package wol-panel for macOS darwin (Apple Silicon)
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

Write-Host "Building wol-panel v$Version darwin/arm64" -ForegroundColor Cyan
$env:CGO_ENABLED = "0"; $env:GOOS = "darwin"; $env:GOARCH = "arm64"
go build -ldflags "-s -w -X wol-panel/version.Version=$Version -X wol-panel/version.Arch=arm64 -X 'wol-panel/version.BuildTime=$BuildTime'" -o build\wol-panel .
Remove-Item Env:\CGO_ENABLED, Env:\GOOS, Env:\GOARCH -ErrorAction SilentlyContinue

$ReleaseDir = "release"
New-Item -ItemType Directory -Force -Path $ReleaseDir | Out-Null

$PkgName = "wol-panel-$Version-macos-darwin"
$PkgDir = Join-Path $ReleaseDir $PkgName
if (Test-Path $PkgDir) { Remove-Item -Recurse -Force $PkgDir }
New-Item -ItemType Directory -Force -Path $PkgDir | Out-Null

Copy-Item "build\wol-panel" $PkgDir
Copy-Item "config.template.json" $PkgDir
Copy-Item "wol-panel.service" $PkgDir

$ZipPath = Join-Path $ReleaseDir "$PkgName.zip"
if (Test-Path $ZipPath) { Remove-Item -Force $ZipPath }
Compress-Archive -Path (Join-Path $PkgName "*") -DestinationPath $ZipPath
Remove-Item -Recurse -Force $PkgDir

Write-Host "Done: $ZipPath" -ForegroundColor Green
