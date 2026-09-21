package vscan

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const MaxSecretSurfaces = 40

// ListSecretSurfaces returns path names that look like secret/credential material.
// Names only — never file contents. Part of Guard (G), not Learn.
func ListSecretSurfaces(root string) ([]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	skipDir := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, ".frontier": true,
		"__pycache__": true, "dist": true, "build": true, "testdata": true,
		"venv": true, ".venv": true, "site-packages": true,
	}
	var out []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skipDir[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		base := strings.ToLower(filepath.Base(path))
		if looksSecretSurface(base) || looksSecretSurface(rel) {
			out = append(out, rel)
		}
		if len(out) >= MaxSecretSurfaces*2 {
			return filepath.SkipAll
		}
		return nil
	})
	sort.Strings(out)
	if len(out) > MaxSecretSurfaces {
		out = out[:MaxSecretSurfaces]
	}
	return out, nil
}

func looksSecretSurface(s string) bool {
	low := strings.ToLower(filepath.ToSlash(s))
	base := filepath.Base(low)
	if base == ".env" || strings.HasPrefix(base, ".env.") {
		return true
	}
	if strings.HasSuffix(base, ".pem") || strings.HasSuffix(base, ".key") {
		return true
	}
	switch base {
	case "credentials.json", "secrets.json", "auth.json", "id_rsa", "id_ed25519":
		return true
	}
	if strings.Contains(base, "id_rsa") || strings.Contains(base, "id_ed25519") {
		return true
	}
	return false
}
