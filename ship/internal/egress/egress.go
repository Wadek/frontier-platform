package egress

import "strings"

// SummarizeDiff returns a short, source-light summary suitable for cloud/tier-2 models.
// It strips patch hunks and keeps only stat-like lines.
func SummarizeDiff(diffStat string) string {
	lines := strings.Split(diffStat, "\n")
	var keep []string
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		// drop unified diff lines if present
		if strings.HasPrefix(ln, "+++") || strings.HasPrefix(ln, "---") || strings.HasPrefix(ln, "@@") || strings.HasPrefix(ln, "+") || strings.HasPrefix(ln, "-") {
			continue
		}
		keep = append(keep, ln)
		if len(keep) >= 40 {
			break
		}
	}
	if len(keep) == 0 {
		return "(no file-level changes summarized)"
	}
	return strings.Join(keep, "\n")
}

func AdviceFromStat(stat string) string {
	s := strings.ToLower(stat)
	switch {
	case strings.Contains(s, "go.mod"), strings.Contains(s, "package.json"):
		return "dependency manifest changed — review supply chain before push"
	case strings.Contains(s, ".env"):
		return "BLOCK: .env-like path in diff — remove secrets before any remote"
	case stat == "" || stat == "(no file-level changes summarized)":
		return "clean tree or no diff — nothing to push"
	default:
		return "diff looks routine — elevate to operator to commit, then gate before push"
	}
}
