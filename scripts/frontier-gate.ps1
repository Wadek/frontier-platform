# Frontier pre-push/pre-commit gate for the frontier-platform monorepo.
# Resolves the frontier binary from FRONTIER_RUNTIME or ./runtime/bin.

param(
    [Parameter(Mandatory = $false)]
    [ValidateSet('pre-commit', 'pre-push')]
    [string]$Hook = 'pre-push'
)

$ErrorActionPreference = 'Stop'
$RepoRoot = Split-Path -Parent $PSScriptRoot
$Runtime = if ($env:FRONTIER_RUNTIME) { $env:FRONTIER_RUNTIME } else { Join-Path $RepoRoot 'runtime' }
$FrontierExe = Join-Path $Runtime 'bin\frontier.exe'
$PipelineDoc = Join-Path $RepoRoot 'PIPELINE.md'

if (-not (Test-Path $FrontierExe)) {
    Write-Host "frontier binary not found: $FrontierExe" -ForegroundColor Red
    Write-Host "Set FRONTIER_RUNTIME or place binaries under ./runtime/bin" -ForegroundColor Yellow
    Write-Host "pipeline: $PipelineDoc" -ForegroundColor Yellow
    exit 1
}

$env:FRONTIER_SOFT = '0'
& $FrontierExe plan
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
& $FrontierExe apply
exit $LASTEXITCODE
