# Install Frontier so the interface is plain `git`.
# Requires Go and Git on PATH (or set GO_BIN / FRONTIER_RUNTIME).

$ErrorActionPreference = "Stop"
$Repo = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$GoBin = if ($env:GO_BIN) { $env:GO_BIN } else { Split-Path -Parent (Get-Command go -ErrorAction Stop).Source }
$Bin = if ($env:FRONTIER_RUNTIME) { Join-Path $env:FRONTIER_RUNTIME 'bin' } else { Join-Path $Repo 'runtime\bin' }
$RealGit = if ($env:FRONTIER_GIT_BIN) { $env:FRONTIER_GIT_BIN } else { (Get-Command git -ErrorAction Stop).Source }

if (-not (Test-Path $RealGit)) { throw "Real git not found" }
if (-not (Test-Path (Join-Path $GoBin 'go.exe')) -and -not (Get-Command go -ErrorAction SilentlyContinue)) {
  throw "Go not found (set GO_BIN or put go on PATH)"
}

New-Item -ItemType Directory -Force -Path $Bin | Out-Null
$env:Path = "$GoBin;" + $env:Path

Push-Location $Repo
go build -o (Join-Path $Bin 'frontier-git.exe') ./cmd/frontier-git
go build -o (Join-Path $Bin 'frontier.exe') ./cmd/frontier
Copy-Item -Force (Join-Path $Bin 'frontier-git.exe') (Join-Path $Bin 'git.exe')
Pop-Location

[Environment]::SetEnvironmentVariable("FRONTIER_GIT_BIN", $RealGit, "User")
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$parts = @($userPath -split ';' | Where-Object { $_ -and ($_ -ne $Bin) })
$newPath = ($Bin + ';' + ($parts -join ';')).TrimEnd(';')
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")

Write-Host "Installed."
Write-Host "  git shim:  $(Join-Path $Bin 'git.exe')"
Write-Host "  frontier:  $(Join-Path $Bin 'frontier.exe')"
Write-Host "  engine:    $RealGit"
Write-Host "Open a NEW terminal, then run:"
Write-Host "  git frontier explain"
Write-Host "  frontier V"
Write-Host "  Get-Command git, frontier"
