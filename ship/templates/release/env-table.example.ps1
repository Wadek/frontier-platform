# Example environment table — copy and edit per app.
# Gate: feature-branch push = plan+apply; production deploy = release-check.

$DeployEnvs = @{
    test = @{
        Branch     = "dev"
        Compose    = "docker-compose.test.yml"
        Gate       = "push"
        DataDir    = "./data.test"
        BackupData = $true
    }
    prod = @{
        Branch     = "main"
        Compose    = "docker-compose.yml"
        Gate       = "release"
        DataDir    = "./data"
        BackupData = $true
    }
}
