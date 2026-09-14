package optimize

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type hotspot struct {
	Path     string
	FuncName string
	Start    int
	End      int
	Snippet  string
}

var (
	reGoFunc    = regexp.MustCompile(`^func\s+(\([^)]+\)\s*)?([A-Za-z0-9_]+)\s*\(`)
	rePyDef     = regexp.MustCompile(`^(\s*)def\s+([A-Za-z0-9_]+)\s*\(`)
	reJSFunc    = regexp.MustCompile(`^\s*(async\s+)?function\s+([A-Za-z0-9_]+)\s*\(`)
	reJSMethod  = regexp.MustCompile(`^\s*(async\s+)?([A-Za-z0-9_]+)\s*\([^)]*\)\s*\{`)
	reJSArrow   = regexp.MustCompile(`^\s*(?:export\s+)?(?:const|let|var)\s+([A-Za-z0-9_]+)\s*=\s*(async\s*)?\(`)
)

func scanLargeFunctions(root string) ([]hotspot, error) {
	var out []hotspot
	skipDir := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, ".frontier": true,
		"__pycache__": true, "dist": true, "build": true, "testdata": true,
		"library": true, "covers": true, "postgres": true, "bin": true, "ledgers": true,
		"venv": true, ".venv": true, "site-packages": true,
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skipDir[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go", ".py", ".js", ".ts", ".tsx", ".jsx":
		default:
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		hs, _ := findHotspotsInFile(path, filepath.ToSlash(rel), ext)
		out = append(out, hs...)
		return nil
	})
	// Prefer largest first
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			si := out[i].End - out[i].Start
			sj := out[j].End - out[j].Start
			if sj > si {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, err
}

func findHotspotsInFile(abs, rel, ext string) ([]hotspot, error) {
	f, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	type start struct {
		name string
		line int // 1-based
		indent int
	}
	var opens []start
	var out []hotspot

	flush := func(endLine int) {
		if len(opens) == 0 {
			return
		}
		s := opens[len(opens)-1]
		opens = opens[:len(opens)-1]
		span := endLine - s.line + 1
		if span < MaxFuncLines {
			return
		}
		snippet := sliceLines(lines, s.line-1, min(endLine, s.line-1+15))
		out = append(out, hotspot{
			Path: rel, FuncName: s.name, Start: s.line, End: endLine, Snippet: snippet,
		})
	}

	for i, line := range lines {
		lineNo := i + 1
		trim := strings.TrimRight(line, "\r")

		switch ext {
		case ".go":
			if m := reGoFunc.FindStringSubmatch(trim); m != nil {
				// close previous at previous line
				if len(opens) > 0 {
					flush(lineNo - 1)
				}
				opens = append(opens, start{name: m[2], line: lineNo})
			}
		case ".py":
			if m := rePyDef.FindStringSubmatch(trim); m != nil {
				ind := len(m[1])
				// close nested/outer defs that are at indent >= this
				for len(opens) > 0 && opens[len(opens)-1].indent >= ind {
					flush(lineNo - 1)
				}
				opens = append(opens, start{name: m[2], line: lineNo, indent: ind})
				continue
			}
			// End current def when a non-empty, non-comment line is at indent <= def indent
			if len(opens) > 0 {
				stripped := strings.TrimSpace(trim)
				if stripped != "" && !strings.HasPrefix(stripped, "#") {
					ind := len(trim) - len(strings.TrimLeft(trim, " \t"))
					for len(opens) > 0 && ind <= opens[len(opens)-1].indent {
						flush(lineNo - 1)
					}
				}
			}
		case ".js", ".ts", ".tsx", ".jsx":
			name := ""
			if m := reJSFunc.FindStringSubmatch(trim); m != nil {
				name = m[2]
			} else if m := reJSArrow.FindStringSubmatch(trim); m != nil {
				name = m[1]
			} else if m := reJSMethod.FindStringSubmatch(trim); m != nil && !strings.HasPrefix(strings.TrimSpace(trim), "if") &&
				!strings.HasPrefix(strings.TrimSpace(trim), "for") && !strings.HasPrefix(strings.TrimSpace(trim), "while") &&
				!strings.HasPrefix(strings.TrimSpace(trim), "switch") {
				name = m[2]
			}
			if name != "" && name != "if" && name != "for" && name != "while" && name != "switch" && name != "catch" {
				if len(opens) > 0 {
					flush(lineNo - 1)
				}
				opens = append(opens, start{name: name, line: lineNo})
			}
		}
	}
	// EOF flush
	for len(opens) > 0 {
		flush(len(lines))
	}
	return out, nil
}

func sliceLines(lines []string, from, toExclusive int) string {
	if from < 0 {
		from = 0
	}
	if toExclusive > len(lines) {
		toExclusive = len(lines)
	}
	if from >= toExclusive {
		return ""
	}
	return strings.Join(lines[from:toExclusive], "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
