// Package frontierai is the cloud fallback LLM path (DeepSeek) with a local token ledger.
// Cascade: process (caller) → laptop Ollama → habitat Qwen (mock until live) → DeepSeek.
package frontierai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	DefaultDeepSeekBase = "https://api.deepseek.com"
	DefaultDeepSeekModel = "deepseek-chat"
	DefaultOllamaBase   = "http://127.0.0.1:11434"
	DefaultHabitatBase  = "http://127.0.0.1:11435" // mock / future Mac mini endpoint
)

// Route names recorded in the token ledger.
const (
	RouteProcess = "process"
	RouteOllama  = "laptop_ollama"
	RouteHabitat = "habitat_qwen"
	RouteDeepSeek = "frontier_ai_deepseek"
)

// Usage is one completion's token counts.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Result is a completed chat turn plus routing metadata.
type Result struct {
	Route   string `json:"route"`
	Model   string `json:"model"`
	Content string `json:"content"`
	Usage   Usage  `json:"usage"`
}

// Client talks to OpenAI-compatible chat APIs (DeepSeek, Ollama).
type Client struct {
	HTTP       *http.Client
	DeepSeekKey string
	DeepSeekBase string
	DeepSeekModel string
	OllamaBase   string
	HabitatBase  string
	HabitatMock  bool // when true, habitat probe always fails (predict fallthrough)
}

// NewFromEnv builds a Client from environment.
// DEEPSEEK_API_KEY (required for DeepSeek), optional DEEPSEEK_BASE_URL, DEEPSEEK_MODEL,
// OLLAMA_HOST, FRONTIER_HABITAT_BASE, FRONTIER_HABITAT_MOCK=1 (default mock until live).
func NewFromEnv() *Client {
	habitatMock := true
	if v := strings.TrimSpace(os.Getenv("FRONTIER_HABITAT_MOCK")); v == "0" || strings.EqualFold(v, "false") {
		habitatMock = false
	}
	base := strings.TrimSpace(os.Getenv("DEEPSEEK_BASE_URL"))
	if base == "" {
		base = DefaultDeepSeekBase
	}
	model := strings.TrimSpace(os.Getenv("DEEPSEEK_MODEL"))
	if model == "" {
		model = strings.TrimSpace(os.Getenv("TIPS_DS_MODEL"))
	}
	if model == "" {
		model = DefaultDeepSeekModel
	}
	ollama := strings.TrimSpace(os.Getenv("OLLAMA_HOST"))
	if ollama == "" {
		ollama = DefaultOllamaBase
	}
	habitat := strings.TrimSpace(os.Getenv("FRONTIER_HABITAT_BASE"))
	if habitat == "" {
		habitat = DefaultHabitatBase
	}
	key, _ := ResolveDeepSeekKey()
	return &Client{
		HTTP:          &http.Client{Timeout: 120 * time.Second},
		DeepSeekKey:   key,
		DeepSeekBase:  strings.TrimRight(base, "/"),
		DeepSeekModel: model,
		OllamaBase:    strings.TrimRight(ollama, "/"),
		HabitatBase:   strings.TrimRight(habitat, "/"),
		HabitatMock:   habitatMock,
	}
}

// Probe reports which cascade tier would answer (without spending completion tokens).
func (c *Client) Probe(ctx context.Context) (route string, detail string) {
	if c == nil {
		return RouteDeepSeek, "nil client"
	}
	if ok, why := c.probeOllama(ctx); ok {
		return RouteOllama, why
	} else {
		detail = "ollama: " + why
	}
	if ok, why := c.probeHabitat(ctx); ok {
		return RouteHabitat, why
	} else {
		detail += "; habitat: " + why
	}
	if c.DeepSeekKey == "" {
		_, src := ResolveDeepSeekKey()
		return RouteDeepSeek, detail + "; deepseek: missing key (source=" + string(src) + "; run: frontier ai secrets status)"
	}
	_, src := ResolveDeepSeekKey()
	return RouteDeepSeek, detail + "; deepseek: key present (" + string(src) + ")"
}

func (c *Client) probeOllama(ctx context.Context) (bool, string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.OllamaBase+"/api/tags", nil)
	if err != nil {
		return false, err.Error()
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return false, "unreachable"
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, "ok"
	}
	return false, fmt.Sprintf("http %d", resp.StatusCode)
}

func (c *Client) probeHabitat(ctx context.Context) (bool, string) {
	if c.HabitatMock {
		return false, "mocked fallthrough"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.HabitatBase+"/api/tags", nil)
	if err != nil {
		return false, err.Error()
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return false, "unreachable"
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, "ok"
	}
	return false, fmt.Sprintf("http %d", resp.StatusCode)
}

type chatReq struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResp struct {
	Model   string `json:"model"`
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

// Complete runs the cascade and returns the first successful completion.
// Laptop Ollama is preferred when healthy; habitat is mocked by default; else DeepSeek.
func (c *Client) Complete(ctx context.Context, system, user string) (*Result, error) {
	if c == nil {
		return nil, fmt.Errorf("frontierai: nil client")
	}
	if ok, _ := c.probeOllama(ctx); ok {
		res, err := c.completeOpenAICompat(ctx, c.OllamaBase+"/v1/chat/completions", "", pickOllamaModel(), system, user)
		if err == nil {
			res.Route = RouteOllama
			return res, nil
		}
	}
	if ok, _ := c.probeHabitat(ctx); ok {
		res, err := c.completeOpenAICompat(ctx, c.HabitatBase+"/v1/chat/completions", "", "qwen", system, user)
		if err == nil {
			res.Route = RouteHabitat
			return res, nil
		}
	}
	if c.DeepSeekKey == "" {
		return nil, fmt.Errorf("frontierai: no local model and DEEPSEEK_API_KEY is empty")
	}
	res, err := c.completeOpenAICompat(ctx, c.DeepSeekBase+"/chat/completions", c.DeepSeekKey, c.DeepSeekModel, system, user)
	if err != nil {
		return nil, err
	}
	res.Route = RouteDeepSeek
	return res, nil
}

func pickOllamaModel() string {
	if m := strings.TrimSpace(os.Getenv("OLLAMA_MODEL")); m != "" {
		return m
	}
	return "qwen2.5-coder:latest"
}

func (c *Client) completeOpenAICompat(ctx context.Context, url, bearer, model, system, user string) (*Result, error) {
	body, _ := json.Marshal(chatReq{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Stream: false,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("frontierai: http %d: %s", resp.StatusCode, trimErr(raw))
	}
	var parsed chatResp
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("frontierai: decode: %w", err)
	}
	content := ""
	if len(parsed.Choices) > 0 {
		content = parsed.Choices[0].Message.Content
	}
	if parsed.Usage.TotalTokens == 0 {
		parsed.Usage.TotalTokens = parsed.Usage.PromptTokens + parsed.Usage.CompletionTokens
	}
	m := parsed.Model
	if m == "" {
		m = model
	}
	return &Result{Model: m, Content: content, Usage: parsed.Usage}, nil
}

func trimErr(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 240 {
		return s[:240] + "…"
	}
	return s
}
