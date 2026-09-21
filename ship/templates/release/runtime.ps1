# Container runtime adapter: Docker or Podman (compose-compatible).
function Get-ComposeCommand {
    $rt = if ($env:RUNTIME) { $env:RUNTIME.ToLowerInvariant() } else { "docker" }
    switch ($rt) {
        "podman" { return @("podman", "compose") }
        default { return @("docker", "compose") }
    }
}

function Invoke-Compose {
    param([Parameter(ValueFromRemainingArguments = $true)][string[]]$ComposeArgs)
    $cmd = Get-ComposeCommand
    & $cmd[0] $cmd[1] @ComposeArgs
    if ($LASTEXITCODE -ne 0) { throw "compose failed: $($ComposeArgs -join ' ')" }
}
