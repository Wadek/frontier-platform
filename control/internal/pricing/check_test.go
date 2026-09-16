package pricing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const goodPage = `<html><body><p>Standard price: 01:00 - 04:00 and 06:00 - 10:00 UTC</p></body></html>`

func TestFindWindow(t *testing.T) {
	if got := FindWindow(goodPage); got != ExpectedWindow {
		t.Fatalf("FindWindow = %q, want %q", got, ExpectedWindow)
	}
	if got := FindWindow("<p>no window here</p>"); got != "" {
		t.Fatalf("FindWindow on windowless page = %q, want empty", got)
	}
}

func TestVerifyWindow(t *testing.T) {
	if err := VerifyWindow(goodPage); err != nil {
		t.Fatalf("valid page rejected: %v", err)
	}
	if err := VerifyWindow("<p>nothing</p>"); err == nil {
		t.Error("page without a window accepted")
	}
	drifted := `<p>01:00 - 05:00 and 06:00 - 10:00 UTC</p>`
	if err := VerifyWindow(drifted); err == nil {
		t.Error("drifted window accepted")
	}
}

func TestVerifyWindowToleratesWhitespace(t *testing.T) {
	spaced := "<p>01:00  -  04:00 and 06:00  -  10:00   UTC</p>"
	if err := VerifyWindow(spaced); err != nil {
		t.Errorf("whitespace variant should not count as drift: %v", err)
	}
}

func TestCheckAgainstServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Error("User-Agent not set")
		}
		_, _ = w.Write([]byte(goodPage))
	}))
	defer srv.Close()

	got, err := Check(context.Background(), srv.Client(), srv.URL)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if got != ExpectedWindow {
		t.Fatalf("Check window = %q", got)
	}
}

func TestCheckServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	if _, err := Check(context.Background(), srv.Client(), srv.URL); err == nil {
		t.Error("500 response accepted")
	}
}

func TestAppendLogCreatesDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "pricing-check.log")
	if err := AppendLog(path, "line one"); err != nil {
		t.Fatal(err)
	}
	if err := AppendLog(path, "line two"); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(raw), "\n") != 2 {
		t.Fatalf("log = %q, want 2 lines", string(raw))
	}
}

func TestLogLine(t *testing.T) {
	now := time.Date(2026, 9, 14, 7, 0, 0, 0, time.UTC) // Monday peak
	line := LogLine(now, ExpectedWindow, nil)
	if !strings.Contains(line, "expected_ok=true") {
		t.Errorf("line = %q", line)
	}
	if !strings.Contains(line, "peak_now=true") {
		t.Errorf("peak not detected at Mon 07:00: %q", line)
	}
	bad := LogLine(now, "", errFake)
	if !strings.Contains(bad, "expected_ok=false") {
		t.Errorf("failed line = %q", bad)
	}
}

type fakeErr struct{}

func (fakeErr) Error() string { return "boom" }

var errFake = fakeErr{}
