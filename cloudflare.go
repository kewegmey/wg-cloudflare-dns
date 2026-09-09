package main

import (
	"context"
	"fmt"
	"net"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareClient struct {
	api    *cloudflare.API
	zoneID string
}

func NewCloudflareClient(apiToken, zoneID string) (*CloudflareClient, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("create cloudflare client: %w", err)
	}

	return &CloudflareClient{api: api, zoneID: zoneID}, nil
}

func recordTypeForIP(ip string) (string, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", fmt.Errorf("invalid IP address: %s", ip)
	}
	if parsed.To4() != nil {
		return "A", nil
	}
	return "AAAA", nil
}

func (c *CloudflareClient) UpsertDNSRecord(ctx context.Context, name, ip string) error {
	recordType, err := recordTypeForIP(ip)
	if err != nil {
		return err
	}

	records, _, err := c.api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(c.zoneID), cloudflare.ListDNSRecordsParams{
		Type: recordType,
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("list DNS records: %w", err)
	}

	for _, record := range records {
		if record.Name != name || record.Type != recordType {
			continue
		}

		if record.Content == ip {
			return nil
		}

		_, err := c.api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(c.zoneID), cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Type:    recordType,
			Name:    name,
			Content: ip,
			TTL:     record.TTL,
			Proxied: record.Proxied,
		})
		if err != nil {
			return fmt.Errorf("update DNS record %s: %w", record.ID, err)
		}

		return nil
	}

	proxied := false
	_, err = c.api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(c.zoneID), cloudflare.CreateDNSRecordParams{
		Type:    recordType,
		Name:    name,
		Content: ip,
		TTL:     1,
		Proxied: &proxied,
	})
	if err != nil {
		return fmt.Errorf("create DNS record: %w", err)
	}

	return nil
}
