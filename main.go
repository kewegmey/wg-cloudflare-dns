package main

import (
	"context"
	"flag"
	"log"
	"time"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	interfaceName := flag.String("interface", "wg0", "WireGuard interface to monitor")
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client, err := NewCloudflareClient(cfg.Cloudflare.APIToken, cfg.Cloudflare.ZoneID)
	if err != nil {
		log.Fatalf("failed to initialize Cloudflare client: %v", err)
	}

	ticker := time.NewTicker(cfg.UpdateInterval)
	defer ticker.Stop()

	runUpdate := func() {
		peerIPs, err := GetWireGuardPeerIPs(*interfaceName)
		if err != nil {
			log.Printf("wireguard read failed: %v", err)
			return
		}

		for _, peer := range cfg.Peers {
			ip, ok := peerIPs[peer.PublicKey]
			if !ok {
				log.Printf("peer %s has no endpoint IP; skipping DNS update for %s", peer.PublicKey, peer.DNSName)
				continue
			}

			if err := client.UpsertDNSRecord(context.Background(), peer.DNSName, ip); err != nil {
				log.Printf("failed updating %s to %s: %v", peer.DNSName, ip, err)
				continue
			}

			log.Printf("updated %s to %s", peer.DNSName, ip)
		}
	}

	runUpdate()
	for range ticker.C {
		runUpdate()
	}
}
