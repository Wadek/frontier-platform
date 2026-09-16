// Official pricing-window verification.
//
// The peak windows are a platform constant (see pricing.go). This file checks the
// published DeepSeek page still advertises the same window, so a vendor price or
// schedule change is noticed instead of silently invalidating routing.
//
// Intended to run as an OS-scheduled one-shot: run, log one line, exit. Never polls.
package pricing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// OfficialURL is the published pricing page.
const OfficialURL = "https://api-docs.deepseek.com/quick_start/pricing"

// ExpectedWindow is the peak window this package's constants assume.
const ExpectedWindow = "01:00 - 04:00 and 06:00 - 10:00 UTC"

var windowRe = regexp.MustCompile(`\d{2}:00\s*-\s*\d{2}:00\s*and\s*\d{2}:00\s*-\s*\d{2}:00\s*UTC`)

// FindWindow returns the peak-window sentence found on the page, or "".
func FindWindow(page string) string {
	return windowRe.FindString(page)
}

// VerifyWindow reports whether the page still carries ExpectedWindow.
// Whitespace is normalised first: vendors re-render their pages, and a stray
// double space is not a price change.
func VerifyWindow(page string) error {
	found := FindWindow(page)
	if found == "" {
		return fmt.Errorf("no peak-window line found on page")
	}
	if normalize(found) != normalize(ExpectedWindow) {
		return fmt.Errorf("window drift: page says %q, platform assumes %q", found, ExpectedWindow)
	}
	return nil
}

func normalize(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// FetchPage retrieves a URL with a bounded timeout.
func FetchPage(ctx context.Context, client *http.Client, url string) (string, error) {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "frontier-platform/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %s", url, resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Check fetches the official page and verifies the window.
func Check(ctx context.Context, client *http.Client, url string) (string, error) {
	page, err := FetchPage(ctx, client, url)
	if err != nil {
		return "", err
	}
	found := FindWindow(page)
	if err := VerifyWindow(page); err != nil {
		return found, err
	}
	return found, nil
}

// LogPath resolves $FRONTIER_RUNTIME/pricing-check.log, else ./runtime/pricing-check.log.
func LogPath() string {
	base := os.Getenv("FRONTIER_RUNTIME")
	if base == "" {
		base = "runtime"
	}
	return filepath.Join(base, "pricing-check.log")
}

// AppendLog appends one line to path, creating parent directories.
func AppendLog(path, line string) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line + "\n")
	return err
}

// LogLine formats the one-line audit record.
func LogLine(now time.Time, window string, err error) string {
	if window == "" {
		window = "missing"
	}
	if err != nil && window == "missing" {
		window = "fetch error: " + err.Error()
	}
	return fmt.Sprintf("%s | window=%s | expected_ok=%t | peak_now=%t",
		now.UTC().Format(time.RFC3339), window, err == nil, IsPeak(now))
}
