# Diagnosis commands

Read-only. Run these before changing anything.

## GitHub side (`gh` authenticated)

```sh
# Recent runs of one workflow, with status and timing
gh run list --workflow <file.yml> --limit 10 \
  --json databaseId,status,conclusion,createdAt,startedAt,updatedAt

# Per-step durations for one run
gh run view <run-id> --json jobs \
  --jq '.jobs[] | .name as $j | .steps[] |
        [$j, .name, .startedAt, .completedAt] | @tsv'

# Jobs currently waiting for a runner
gh run list --status queued --limit 20

# Self-hosted runners and whether they are busy
gh api repos/<owner>/<repo>/actions/runners \
  --jq '.runners[] | [.name, .status, .busy, ([.labels[].name] | join(","))] | @tsv'
```

Queue time = `startedAt - createdAt`. Run time = `updatedAt - startedAt`.
Report them separately.

## Host side

Linux:

```sh
nproc; lscpu | grep -E 'Model name|MHz'
free -h; swapon --show
lsblk -d -o NAME,ROTA,SIZE,MODEL      # ROTA=0 means SSD
iostat -x 5 3                          # %iowait, device utilization
docker system df                       # build cache size
```

Windows (PowerShell):

```powershell
Get-CimInstance Win32_Processor | Select Name,NumberOfCores,MaxClockSpeed
Get-CimInstance Win32_OperatingSystem | Select TotalVisibleMemorySize,FreePhysicalMemory
Get-PhysicalDisk | Select FriendlyName,MediaType,BusType
Get-Service actions.runner.*
```
