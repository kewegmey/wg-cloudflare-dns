package main

import "testing"

func TestRecordTypeForIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		wantType string
		wantErr  bool
	}{
		{name: "ipv4", ip: "198.51.100.10", wantType: "A"},
		{name: "ipv6", ip: "2001:db8::1", wantType: "AAAA"},
		{name: "invalid", ip: "not-an-ip", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := recordTypeForIP(tt.ip)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.ip)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantType {
				t.Fatalf("expected %s, got %s", tt.wantType, got)
			}
		})
	}
}
