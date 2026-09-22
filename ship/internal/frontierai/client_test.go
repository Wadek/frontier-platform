package frontierai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompleteDeepSeekFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/tags":
			http.Error(w, "down", http.StatusBadGateway)
		case r.URL.Path == "/chat/completions":
			_ = json.NewEncoder(w).Encode(chatResp{
				Model: "deepseek-chat",
				Choices: []struct {
					Message chatMessage `json:"message"`
				}{{Message: chatMessage{Role: "assistant", Content: "ok"}}},
				Usage: Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := &Client{
		HTTP:          srv.Client(),
		DeepSeekKey:   "test-key",
		DeepSeekBase:  srv.URL,
		DeepSeekModel: "deepseek-chat",
		OllamaBase:    srv.URL, // /api/tags will 502
		HabitatBase:   srv.URL,
		HabitatMock:   true,
	}
	res, err := c.Complete(context.Background(), "sys", "user")
	if err != nil {
		t.Fatal(err)
	}
	if res.Route != RouteDeepSeek || res.Usage.TotalTokens != 15 || res.Content != "ok" {
		t.Fatalf("%+v", res)
	}
}

func TestProbeHabitatMock(t *testing.T) {
	c := &Client{
		HTTP:         http.DefaultClient,
		DeepSeekKey:  "x",
		OllamaBase:   "http://127.0.0.1:1", // closed port → unreachable
		HabitatBase:  "http://127.0.0.1:1",
		HabitatMock:  true,
		DeepSeekBase: DefaultDeepSeekBase,
	}
	route, detail := c.Probe(context.Background())
	if route != RouteDeepSeek {
		t.Fatalf("want deepseek got %s (%s)", route, detail)
	}
	if detail == "" || !contains(detail, "mocked") {
		t.Fatalf("detail=%q", detail)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
