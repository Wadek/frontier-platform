package frontierai

import (
	"os"
	"testing"
)

func TestResolvePrefersEnvOverMissingKeystore(t *testing.T) {
	t.Setenv(DeepSeekSecretName, "sk-test-from-env")
	key, src := ResolveDeepSeekKey()
	if key != "sk-test-from-env" || src != KeySourceEnv {
		t.Fatalf("got %q %s", key, src)
	}
}

func TestStatusDoesNotLeakKey(t *testing.T) {
	t.Setenv(DeepSeekSecretName, "sk-should-not-appear-in-hint-as-value-only")
	st := StatusDeepSeek()
	if st.ResolvedFrom != string(KeySourceEnv) || !st.EnvSet {
		t.Fatalf("%+v", st)
	}
	if st.Hint == "" {
		t.Fatal("expected hint")
	}
}

func TestKeystoreRoundTrip(t *testing.T) {
	t.Setenv(DeepSeekSecretName, "")
	_ = os.Unsetenv(DeepSeekSecretName)
	_ = DeleteDeepSeekKey()

	if err := SetDeepSeekKey("sk-roundtrip-test-key"); err != nil {
		t.Fatalf("set: %v", err)
	}
	t.Cleanup(func() { _ = DeleteDeepSeekKey() })

	raw, err := keystoreGet(DeepSeekSecretName)
	if err != nil {
		t.Fatalf("direct get: %v path=%s", err, keystorePath(DeepSeekSecretName))
	}
	if raw != "sk-roundtrip-test-key" {
		t.Fatalf("direct get mismatch %q", raw)
	}

	key, src := ResolveDeepSeekKey()
	if src != KeySourceKeystore {
		t.Fatalf("want keystore got %s (key len %d)", src, len(key))
	}
	if key != "sk-roundtrip-test-key" {
		t.Fatalf("resolve mismatch")
	}
}
