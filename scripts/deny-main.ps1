# Refuse commits on main/master.
$ErrorActionPreference = 'Stop'
$branch = (& git branch --show-current | Out-String).Trim()
if ($branch -eq 'main' -or $branch -eq 'master') {
    Write-Host "frontier deny commit on $branch - use a feature branch" -ForegroundColor Red
    exit 1
}
exit 0
