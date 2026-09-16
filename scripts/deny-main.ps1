# Refuse commits on protected branches (main, master, frontier/public-root).
$ErrorActionPreference = 'Stop'
$branch = (& git branch --show-current | Out-String).Trim()
if ($branch -eq 'main' -or $branch -eq 'master' -or $branch -eq 'frontier/public-root') {
    Write-Host "frontier deny commit on $branch - use a feature branch" -ForegroundColor Red
    exit 1
}
exit 0
