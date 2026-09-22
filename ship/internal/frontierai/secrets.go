package frontierai

import (
	"fmt"
	"os"
	"strings"
)

// Secret name used in the local keystore and GitHub Actions secrets.
const DeepSeekSecretName = "DEEPSEEK_API_KEY"

// KeySource describes where ResolveDeepSeekKey found the value.
type KeySource string

const (
	KeySourceEnv      KeySource = "env"       // process env (includes GitHub Actions secrets on runners)
	KeySourceKeystore KeySource = "keystore"  // local OS keystore / DPAPI file
	KeySourceMissing  KeySource = "missing"
)

// ResolveDeepSeekKey returns the API key and where it came from.
//
// Order (wherever the runner/process is):
//  1. Process env DEEPSEEK_API_KEY — GitHub Actions injects repo/org secrets here on both
//     hosted and self-hosted runners when the workflow declares env/secrets.
//  2. Local keystore — Windows DPAPI file under %LOCALAPPDATA%\frontier, or
//     ~/.config/frontier on Unix (for interactive laptop `frontier ai` outside Actions).
func ResolveDeepSeekKey() (key string, src KeySource) {
	if v := strings.TrimSpace(os.Getenv(DeepSeekSecretName)); v != "" {
		return v, KeySourceEnv
	}
	if v, err := keystoreGet(DeepSeekSecretName); err == nil && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v), KeySourceKeystore
	}
	return "", KeySourceMissing
}

// SetDeepSeekKey stores the key in the local keystore (not GitHub).
func SetDeepSeekKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("empty key")
	}
	return keystoreSet(DeepSeekSecretName, key)
}

// DeleteDeepSeekKey removes the local keystore entry.
func DeleteDeepSeekKey() error {
	return keystoreDelete(DeepSeekSecretName)
}

// SecretsStatus is a safe (no secret value) status snapshot.
type SecretsStatus struct {
	EnvSet       bool   `json:"env_set"`
	KeystoreSet  bool   `json:"keystore_set"`
	KeystorePath string `json:"keystore_path"`
	ResolvedFrom string `json:"resolved_from"`
	InActions    bool   `json:"in_github_actions"`
	Hint         string `json:"hint"`
}

// StatusDeepSeek reports whether a key is available without printing it.
func StatusDeepSeek() SecretsStatus {
	_, fromEnv := os.LookupEnv(DeepSeekSecretName)
	envSet := fromEnv && strings.TrimSpace(os.Getenv(DeepSeekSecretName)) != ""
	ksPath := keystorePath(DeepSeekSecretName)
	_, ksErr := keystoreGet(DeepSeekSecretName)
	ksSet := ksErr == nil
	_, src := ResolveDeepSeekKey()
	inActions := os.Getenv("GITHUB_ACTIONS") == "true"
	hint := "Local laptop: frontier ai secrets set deepseek (stores in local keystore). " +
		"GitHub Actions: gh secret set DEEPSEEK_API_KEY — injected as env on the runner (hosted or self-hosted)."
	if src == KeySourceMissing {
		hint = "No key found. DeepSeek does not ship a key — you provide one once. " + hint
	}
	return SecretsStatus{
		EnvSet:       envSet,
		KeystoreSet:  ksSet,
		KeystorePath: ksPath,
		ResolvedFrom: string(src),
		InActions:    inActions,
		Hint:         hint,
	}
}
