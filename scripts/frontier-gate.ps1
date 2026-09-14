# Frontier pre-push gate for the frontier-fleet monorepo.
#
# The habitat installs a global core.hooksPath (D:/wakalabs/frontier-fleet/runtime/hooks)
# whose hook files are lefthook dispatchers, so a repo only runs hooks if it ships a
# lefthook.yml. Lefthook does not forward the hook's positional arguments (remote name /
# url), so the canonical pre-push.ps1 sees an empty remote and exits 0 ("remote is not
# GitHub - local gate skipped"). This wrapper restores that check and keeps the canonical
# fail-closed sequence: plan -> apply.
#
# Canonical doctrine: D:\wakalabs\frontier-fleet\PIPELINE.md

$ErrorActionPreference = "Stop"
$FrontierExe = "D:\wakalabs\frontier-fleet\runtime\bin\frontier.exe"
$env:FRONTIER_SOFT = "0"
$env:FRONTIER_VERBOSE = "1"

function Deny([string]$msg) {
    Write-Host "frontier deny push: $msg" -ForegroundColor Red
    Write-Host "pipeline: D:\wakalabs\frontier-fleet\PIPELINE.md" -ForegroundColor Yellow
    exit 1
}

# git sends "<local ref> <local sha> <remote ref> <remote sha>" lines on stdin.
$raw = ""
try { $raw = [Console]::In.ReadToEnd() } catch { }

$zero = "0" * 40
$hasCommits = $false
$pushingMain = $false
foreach ($line in ($raw -split "`r?`n")) {
    $line = $line.Trim()
    if ($line -eq "") { continue }
    $parts = $line -split "\s+"
    if ($parts.Count -lt 4) { continue }
    if ($parts[1] -eq $zero) { continue }   # branch deletion, nothing to ship
    $hasCommits = $true
    if ($parts[0] -match '(?i)refs/heads/(main|master)$' -or $parts[2] -match '(?i)refs/heads/(main|master)$') {
        $pushingMain = $true
    }
}

if ($pushingMain) { Deny "refusing push to main/master (use a feature branch + PR)" }
if (-not $hasCommits) { exit 0 }

$remotes = (& git remote -v 2>$null) -join "`n"
if ($remotes -notmatch '(?i)github\.com') {
    Write-Host "frontier: no GitHub remote in this repo - local gate skipped"
    exit 0
}

$branch = (& git branch --show-current 2>$null | Out-String).Trim()
if ($branch -eq "main" -or $branch -eq "master") { Deny "current branch is $branch" }
if (-not (Test-Path $FrontierExe)) { Deny "frontier.exe missing at $FrontierExe" }

Write-Host "=== FRONTIER GATE (GitHub) on $branch ===" -ForegroundColor Cyan
& $FrontierExe plan
if ($LASTEXITCODE -ne 0) { Deny "frontier plan failed (exit $LASTEXITCODE)" }
& $FrontierExe apply
if ($LASTEXITCODE -ne 0) { Deny "frontier apply failed (exit $LASTEXITCODE)" }
Write-Host "frontier: plan+apply sealed - allowing git push to GitHub" -ForegroundColor Green
exit 0
