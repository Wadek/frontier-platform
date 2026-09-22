//go:build !windows

package frontierai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func keystoreDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "frontier", "secrets")
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
	return filepath.Join(dir, safe+".secret")
}

func keystoreSet(name, value string) error {
	path := keystorePath(name)
	if path == "" {
		return fmt.Errorf("keystore path unavailable")
	}
	return os.WriteFile(path, []byte(value+"\n"), 0o600)
}

func keystoreGet(name string) (string, error) {
	path := keystorePath(name)
	if path == "" {
		return "", fmt.Errorf("keystore path unavailable")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
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
