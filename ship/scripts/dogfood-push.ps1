# Dogfood: plan → apply → push (fail closed). Refuses main/master.
$ErrorActionPreference = "Stop"
$env:Path = "D:\frontier\bin;C:\Users\waka\sdk\go\bin;C:\Program Files\Git\cmd;" + $env:Path
if (-not $env:FRONTIER_GIT_BIN) {
  $env:FRONTIER_GIT_BIN = "C:\Program Files\Git\cmd\git.exe"
}
$env:FRONTIER_SOFT = "0"
$env:FRONTIER_VERBOSE = "1"

$branch = (& git branch --show-current).Trim()
if ($branch -eq "main" -or $branch -eq "master") {
  Write-Error "Dogfood refuses to push from $branch. Use: git checkout -b frontier/…"
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

Write-Host "=== dogfood OK — open a PR into main ===" -ForegroundColor Green
gh pr create --fill 2>$null
Write-Host "Done."
