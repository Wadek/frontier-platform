//go:build windows

package frontierai

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func keystoreDir() (string, error) {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		return "", fmt.Errorf("LOCALAPPDATA not set")
	}
	dir := filepath.Join(base, "frontier", "secrets")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

func keystorePath(name string) string {
	dir, err := keystoreDir()
	if err != nil {
		return ""
	}
	safe := strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, name)
	return filepath.Join(dir, safe+".dpapi")
}

func runSecretPS(script string, envPath string, stdin string) (string, error) {
	f, err := os.CreateTemp("", "frontier-secret-*.ps1")
	if err != nil {
		return "", err
	}
	ps1 := f.Name()
	_, _ = f.WriteString(script)
	_ = f.Close()
	defer os.Remove(ps1)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", ps1)
	cmd.Env = append(os.Environ(), "FRONTIER_SECRET_PATH="+envPath)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func keystoreSet(name, value string) error {
	path := keystorePath(name)
	if path == "" {
		return fmt.Errorf("keystore path unavailable")
	}
	script := `
$ErrorActionPreference = 'Stop'
$Path = $env:FRONTIER_SECRET_PATH
$Plain = [Console]::In.ReadToEnd()
if (-not $Plain) { throw 'empty stdin' }
$Plain = $Plain.Trim()
$sec = ConvertTo-SecureString $Plain -AsPlainText -Force
$enc = ConvertFrom-SecureString $sec
Set-Content -LiteralPath $Path -Value $enc -Encoding ASCII -NoNewline
`
	_, err := runSecretPS(script, path, value)
	if err != nil {
		return fmt.Errorf("dpapi store: %w", err)
	}
	if st, err := os.Stat(path); err != nil || st.Size() == 0 {
		return fmt.Errorf("dpapi store: file missing or empty after write")
	}
	return nil
}

func keystoreGet(name string) (string, error) {
	path := keystorePath(name)
	if path == "" {
		return "", fmt.Errorf("keystore path unavailable")
	}
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	script := `
$ErrorActionPreference = 'Stop'
$Path = $env:FRONTIER_SECRET_PATH
$enc = Get-Content -LiteralPath $Path -Raw
$sec = ConvertTo-SecureString $enc
$bstr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($sec)
try {
  [Runtime.InteropServices.Marshal]::PtrToStringBSTR($bstr)
} finally {
  [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($bstr)
}
`
	out, err := runSecretPS(script, path, "")
	if err != nil {
		return "", fmt.Errorf("dpapi read: %w", err)
	}
	if out == "" {
		return "", fmt.Errorf("dpapi read: empty")
	}
	return out, nil
}

func keystoreDelete(name string) error {
	path := keystorePath(name)
	if path == "" {
		return nil
	}
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
