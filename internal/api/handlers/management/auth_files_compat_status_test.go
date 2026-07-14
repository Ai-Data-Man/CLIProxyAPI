package management

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/config"
	coreauth "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/auth"
)

// writeCompatTestConfig writes a config.yaml with an openai-compatibility block
// containing two api-key-entries and returns its path plus the loaded Config.
func writeCompatTestConfig(t *testing.T) (string, *config.Config) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `openai-compatibility:
  - name: MixedProvider
    base-url: "https://mixed.api.com"
    api-key-entries:
      - api-key: "active-key"
      - api-key: "disabled-key"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
	return path, cfg
}

// TestPatchAuthFileStatus_OpenAICompatConfigManaged_PersistsDisabled verifies
// that toggling disabled on an openai-compatibility config-managed credential
// writes back to config.yaml at the correct api-key-entry index.
func TestPatchAuthFileStatus_OpenAICompatConfigManaged_PersistsDisabled(t *testing.T) {
	t.Setenv("MANAGEMENT_PASSWORD", "")
	gin.SetMode(gin.TestMode)

	path, cfg := writeCompatTestConfig(t)
	if len(cfg.OpenAICompatibility) != 1 || len(cfg.OpenAICompatibility[0].APIKeyEntries) != 2 {
		t.Fatalf("expected 1 provider with 2 entries, got %+v", cfg.OpenAICompatibility)
	}

	store := &memoryAuthStore{}
	manager := coreauth.NewManager(store, nil, nil)
	record := &coreauth.Auth{
		ID:       "compat-key-1",
		Provider: "mixedprovider",
		FileName: "",
		Attributes: map[string]string{
			"compat_name":  "MixedProvider",
			"entry_index":  "1",
			"provider_key": "mixedprovider",
		},
		Status: coreauth.StatusActive,
	}
	if _, err := manager.Register(context.Background(), record); err != nil {
		t.Fatalf("failed to register auth: %v", err)
	}

	h := NewHandler(cfg, path, manager)

	body := `{"name":"compat-key-1","disabled":true}`
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/v0/management/auth-files/status", strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	h.PatchAuthFileStatus(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	// Reload config from disk and verify entry index 1 is disabled.
	reloaded, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("failed to reload config: %v", err)
	}
	entries := reloaded.OpenAICompatibility[0].APIKeyEntries
	if !entries[1].Disabled {
		t.Errorf("expected entries[1].Disabled=true after toggle, got false")
	}
	if entries[0].Disabled {
		t.Errorf("expected entries[0].Disabled=false (untouched), got true")
	}
}

// TestBuildAuthFileEntry_OpenAICompatConfigManaged_Visible verifies that an
// openai-compatibility config-managed auth (no path, not runtime_only) is NOT
// filtered out of the auth-files list, even when active.
func TestBuildAuthFileEntry_OpenAICompatConfigManaged_Visible(t *testing.T) {
	auth := &coreauth.Auth{
		ID:       "compat-key-0",
		Provider: "mixedprovider",
		FileName: "",
		Attributes: map[string]string{
			"compat_name":  "MixedProvider",
			"entry_index":  "0",
			"provider_key": "mixedprovider",
		},
		Status: coreauth.StatusActive,
	}
	h := &Handler{cfg: &config.Config{}}
	entry := h.buildAuthFileEntry(auth)
	if entry == nil {
		t.Fatal("expected config-managed auth to be visible, got nil (filtered out)")
	}
	if entry["compat_name"] != "MixedProvider" {
		t.Errorf("expected compat_name=MixedProvider, got %v", entry["compat_name"])
	}
	if entry["entry_index"] != "0" {
		t.Errorf("expected entry_index=0, got %v", entry["entry_index"])
	}
	if cm, _ := entry["config_managed"].(bool); !cm {
		t.Errorf("expected config_managed=true, got %v", entry["config_managed"])
	}
}

// TestPatchAuthFileStatus_OpenAICompatConfigManaged_RoundTrip verifies the
// re-enable path: disable then re-enable a config-managed credential, and
// confirm the disabled flag is removed from config.yaml (omitempty).
func TestPatchAuthFileStatus_OpenAICompatConfigManaged_RoundTrip(t *testing.T) {
	t.Setenv("MANAGEMENT_PASSWORD", "")
	gin.SetMode(gin.TestMode)

	path, cfg := writeCompatTestConfig(t)

	store := &memoryAuthStore{}
	manager := coreauth.NewManager(store, nil, nil)
	record := &coreauth.Auth{
		ID:       "compat-key-0",
		Provider: "mixedprovider",
		FileName: "",
		Attributes: map[string]string{
			"compat_name":  "MixedProvider",
			"entry_index":  "0",
			"provider_key": "mixedprovider",
		},
		Status: coreauth.StatusActive,
	}
	if _, err := manager.Register(context.Background(), record); err != nil {
		t.Fatalf("failed to register auth: %v", err)
	}
	h := NewHandler(cfg, path, manager)

	// Step 1: disable
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/v0/management/auth-files/status", strings.NewReader(`{"name":"compat-key-0","disabled":true}`))
	ctx.Request.Header.Set("Content-Type", "application/json")
	h.PatchAuthFileStatus(ctx)
	if rec.Code != http.StatusOK {
		t.Fatalf("disable: expected status %d, got %d with body %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	reloaded, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("failed to reload config after disable: %v", err)
	}
	if !reloaded.OpenAICompatibility[0].APIKeyEntries[0].Disabled {
		t.Errorf("expected entries[0].Disabled=true after disable, got false")
	}

	// Step 2: re-enable
	rec2 := httptest.NewRecorder()
	ctx2, _ := gin.CreateTestContext(rec2)
	ctx2.Request = httptest.NewRequest(http.MethodPatch, "/v0/management/auth-files/status", strings.NewReader(`{"name":"compat-key-0","disabled":false}`))
	ctx2.Request.Header.Set("Content-Type", "application/json")
	h.PatchAuthFileStatus(ctx2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("re-enable: expected status %d, got %d with body %s", http.StatusOK, rec2.Code, rec2.Body.String())
	}

	// Verify disabled flag removed from disk (omitempty should omit false).
	// Search for the YAML key 'disabled:' as a map key, not as a value substring
	// (e.g. 'disabled-key' as an api-key value would falsely match a naive strings.Contains).
	rawBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	if isDisabledKeyPresentInYAML(string(rawBytes)) {
		t.Errorf("expected 'disabled:' key to be omitted from YAML after re-enable (omitempty), but file contains it:\n%s", string(rawBytes))
	}
}

// isDisabledKeyPresentInYAML checks whether any line in the YAML has
// 'disabled' as a map key (i.e. starts with optional whitespace then 'disabled:').
// This avoids matching 'disabled' appearing inside a value string.
func isDisabledKeyPresentInYAML(raw string) bool {
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmed, "disabled:") || strings.HasPrefix(trimmed, "disabled ") {
			return true
		}
	}
	return false
}
