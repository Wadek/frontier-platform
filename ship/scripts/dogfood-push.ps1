# Dogfood: plan → apply → push (fail closed). Refuses dev/main/master.
$ErrorActionPreference = "Stop"
$env:FRONTIER_SOFT = "0"
$env:FRONTIER_VERBOSE = "1"

$branch = (& git branch --show-current).Trim()
if ($branch -eq "main" -or $branch -eq "master" -or $branch -eq "dev") {
  Write-Error "Dogfood refuses to push from $branch. Use: git checkout -b feat/..."
}

Write-Host "=== DOGFOOD on $branch ===" -ForegroundColor Cyan
Write-Host "=== frontier plan ===" -ForegroundColor Cyan
frontier plan
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "=== frontier apply ===" -ForegroundColor Cyan
frontier apply
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "=== git push ===" -ForegroundColor Cyan
git push -u origin HEAD
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "=== dogfood OK -> open a PR into dev ===" -ForegroundColor Green
gh pr create --base dev --fill 2>$null
Write-Host "Done."
