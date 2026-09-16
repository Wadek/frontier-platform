// Package registry discovers projects and their agent packs under a workspace root.
// Spec: english/LEARNING.md (layout) — a project may carry a sibling .agent_<name>/ pack.
package registry

import (
	"os"
	"path/filepath"
	"strings"
)

// Project is one workspace entry and whether it has an agent pack.
type Project struct {
	Name     string `json:"project"`
	AgentDir string `json:"agent_dir"`
	HasAgent bool   `json:"has_agent"`
}

// DefaultExclude are directories that are infrastructure, not projects.
var DefaultExclude = map[string]bool{
	"frontier-platform": true,
	"runtime":           true,
	"node_modules":      true,
	".git":              true,
}

// AgentPrefix is the directory prefix for an agent pack.
const AgentPrefix = ".agent_"

// Scan lists the projects directly under root, sorted by name, skipping hidden
// directories and anything in exclude. A nil exclude uses DefaultExclude.
func Scan(root string, exclude map[string]bool) ([]Project, error) {
	if exclude == nil {
		exclude = DefaultExclude
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	out := make([]Project, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() || strings.HasPrefix(name, ".") || exclude[name] {
			continue
		}
		agentDir := filepath.Join(root, AgentPrefix+name)
		info, err := os.Stat(agentDir)
		out = append(out, Project{
			Name:     name,
			AgentDir: agentDir,
			HasAgent: err == nil && info.IsDir(),
		})
	}
	return out, nil
}
