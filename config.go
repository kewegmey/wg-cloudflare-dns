package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Cloudflare        CloudflareConfig `yaml:"cloudflare"`
	Peers             []PeerMapping    `yaml:"peers"`
	UpdateIntervalRaw string           `yaml:"update_interval"`
	UpdateInterval    time.Duration    `yaml:"-"`
}

type CloudflareConfig struct {
	APIToken string `yaml:"api_token"`
	ZoneID   string `yaml:"zone_id"`
}

type PeerMapping struct {
	PublicKey string `yaml:"public_key"`
	DNSName   string `yaml:"dns_name"`
}

func LoadConfig(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Cloudflare.APIToken == "" {
		return fmt.Errorf("cloudflare.api_token is required")
	}
	if c.Cloudflare.ZoneID == "" {
		return fmt.Errorf("cloudflare.zone_id is required")
	}
	if len(c.Peers) == 0 {
		return fmt.Errorf("at least one peer mapping is required")
	}

	seen := make(map[string]struct{}, len(c.Peers))
	for i, peer := range c.Peers {
		if peer.PublicKey == "" {
			return fmt.Errorf("peers[%d].public_key is required", i)
		}
		if peer.DNSName == "" {
			return fmt.Errorf("peers[%d].dns_name is required", i)
		}
		if _, ok := seen[peer.PublicKey]; ok {
			return fmt.Errorf("duplicate public_key in peers: %s", peer.PublicKey)
		}
		seen[peer.PublicKey] = struct{}{}
	}

	if c.UpdateIntervalRaw == "" {
		c.UpdateInterval = 5 * time.Minute
		return nil
	}

	interval, err := time.ParseDuration(c.UpdateIntervalRaw)
	if err != nil {
		return fmt.Errorf("invalid update_interval: %w", err)
	}
	if interval <= 0 {
		return fmt.Errorf("update_interval must be greater than zero")
	}

	c.UpdateInterval = interval
	return nil
}
