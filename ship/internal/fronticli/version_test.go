package fronticli

import "testing"

func TestProductVersion(t *testing.T) {
	if ProductVersion != "1.0.0.0" {
		t.Fatalf("ProductVersion=%q want 1.0.0.0", ProductVersion)
	}
	if Version != ProductVersion {
		t.Fatalf("Version=%q want %q", Version, ProductVersion)
	}
	if version != ProductVersion {
		t.Fatalf("internal version=%q want %q", version, ProductVersion)
	}
}
