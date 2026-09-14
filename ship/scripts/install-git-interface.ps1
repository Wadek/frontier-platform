# Install Frontier so the interface is plain `git`.
# Run once in an elevated-enough user PowerShell (no admin required).

$ErrorActionPreference = "Stop"
$GoBin = "C:\Users\waka\sdk\go\bin"
$Repo = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$Bin = "D:\frontier\bin"
$RealGit = "C:\Program Files\Git\cmd\git.exe"

if (-not (Test-Path $RealGit)) { throw "Real git not found at $RealGit" }
if (-not (Test-Path "$GoBin\go.exe")) { throw "Go not found at $GoBin\go.exe" }

New-Item -ItemType Directory -Force -Path $Bin | Out-Null
$env:Path = "$GoBin;" + $env:Path
$env:GOROOT = "C:\Users\waka\sdk\go"

Push-Location $Repo
go build -o "$Bin\frontier-git.exe" ./cmd/frontier-git
go build -o "$Bin\frontier.exe" ./cmd/frontier
Copy-Item -Force "$Bin\frontier-git.exe" "$Bin\git.exe"
Pop-Location

[Environment]::SetEnvironmentVariable("FRONTIER_GIT_BIN", $RealGit, "User")
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$parts = @($userPath -split ';' | Where-Object { $_ -and ($_ -ne $Bin) })
$newPath = ($Bin + ';' + ($parts -join ';')).TrimEnd(';')
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")

Write-Host "Installed."
Write-Host "  git shim:  $Bin\git.exe"
Write-Host "  frontier:  $Bin\frontier.exe   (standalone: frontier V | plan | apply)"
Write-Host "  engine:    $RealGit"
Write-Host "Open a NEW terminal, then run:"
Write-Host "  git frontier explain"
Write-Host "  frontier V"
Write-Host "  Get-Command git, frontier"
