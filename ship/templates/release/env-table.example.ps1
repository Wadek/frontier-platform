# Example environment table — copy and edit per app.
# Gate: push = plan+apply; release = release-check only.

$DeployEnvs = @{
    test = @{
        Branch     = "dev"          # or a dedicated test branch
        Compose    = "docker-compose.test.yml"
        Gate       = "push"
        DataDir    = "./data.test"
        BackupData = $true
    }
    alpha = @{
        Branch     = "alpha"
        Compose    = "docker-compose.alpha.yml"
        Gate       = "push"
        DataDir    = $null          # wakagym alpha shares prod data today — document risk
        BackupData = $false
    }
    prod = @{
        Branch     = "main"
        Compose    = "docker-compose.yml"
        Gate       = "release"
        DataDir    = "./data"
        BackupData = $true
    }
}
