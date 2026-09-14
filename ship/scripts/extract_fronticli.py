from pathlib import Path

src = Path(r"C:\Users\waka\src\frontier-ship\cmd\frontier-git\main.go")
dst = Path(r"C:\Users\waka\src\frontier-ship\internal\fronticli\cli.go")
t = src.read_text(encoding="utf-8")
t = t.replace("package main", "package fronticli", 1)

start = t.find("func main() {")
if start < 0:
    raise SystemExit("main not found")
# find matching closing brace for main at column 0
i = start
depth = 0
end = None
for j, ch in enumerate(t[start:], start):
    if ch == "{":
        depth += 1
    elif ch == "}":
        depth -= 1
        if depth == 0:
            end = j + 1
            break
if end is None:
    raise SystemExit("end not found")

new_run = '''
// Version/Commit set via ldflags from cmd wrappers.
var (
	Version = "dev"
	Commit  = "none"
)

// Run executes frontier subcommands (V, S, plan, apply, ...).
// Git passthrough stays in cmd/frontier-git.
func Run(args []string) {
	if Version != "dev" || Commit != "none" {
		version, commit = Version, Commit
	}
	handleMeta(args)
}

// GuardPush is used by the git shim.
func GuardPush(soft bool) error { return guardPush(soft) }

// GuardCommit is used by the git shim.
func GuardCommit(args []string, soft, strict bool) error {
	return guardCommit(args, soft, strict)
}

// FindRealGit exposes real git path for the shim.
func FindRealGit() string { return findRealGit() }

// RunGitPassthrough runs real git and returns exit code.
func RunGitPassthrough(git string, args []string) int {
	return runPassthrough(git, args)
}
'''

t = t[:start] + new_run + t[end:]
# fail() currently calls os.Exit — keep it
dst.write_text(t, encoding="utf-8")
print("wrote", dst, "bytes", dst.stat().st_size)
