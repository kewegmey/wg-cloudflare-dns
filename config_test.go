package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfigParsesIntervalAndPeers(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	content := []byte(`cloudflare:
  api_token: "token"
  zone_id: "zone"
peers:
  - public_key: "peer-key"
    dns_name: "peer.example.com"
update_interval: "10m"
`)
	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.UpdateInterval != 10*time.Minute {
		t.Fatalf("expected 10m interval, got %s", cfg.UpdateInterval)
	}
	if len(cfg.Peers) != 1 {
		t.Fatalf("expected 1 peer mapping, got %d", len(cfg.Peers))
	}
}

func TestLoadConfigDefaultsInterval(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	content := []byte(`cloudflare:
  api_token: "token"
  zone_id: "zone"
peers:
  - public_key: "peer-key"
    dns_name: "peer.example.com"
`)
	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.UpdateInterval != 5*time.Minute {
		t.Fatalf("expected default 5m interval, got %s", cfg.UpdateInterval)
	}
}
