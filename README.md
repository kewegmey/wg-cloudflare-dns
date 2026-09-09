# wg-cloudflare-dns

Automatically update Cloudflare DNS A records with WireGuard peer endpoint IP addresses.

## Features

- Reads WireGuard peers from interface `wg0` (or configurable interface)
- Maps peer public keys to DNS names from YAML config
- Creates or updates Cloudflare DNS A records to match current peer endpoint IP
- Runs updates on a configurable interval

## Configuration

Create `/etc/wg-cloudflare-dns/config.yaml`:

```yaml
cloudflare:
  api_token: "your-api-token"
  zone_id: "your-zone-id"

peers:
  - public_key: "peer-public-key"
    dns_name: "peer.example.com"

update_interval: "5m"
```

`update_interval` defaults to `5m` when omitted.

## Build and run

```bash
go test ./...
go build -o wg-cloudflare-dns .
./wg-cloudflare-dns -config config.yaml -interface wg0
```

## Debian package (Ubuntu)

Build a `.deb` package:

```bash
make package-deb
```

Install:

```bash
sudo dpkg -i build/deb/wg-cloudflare-dns_<version>_amd64.deb
sudo systemctl daemon-reload
sudo systemctl enable --now wg-cloudflare-dns
```

## Systemd service

Installed service file:

- `/lib/systemd/system/wg-cloudflare-dns.service`

Service command:

- `/usr/local/bin/wg-cloudflare-dns -config /etc/wg-cloudflare-dns/config.yaml -interface wg0`

## GitHub Actions

Workflow: `.github/workflows/build.yml`

- Runs tests
- Builds Linux binary
- Builds Debian package
- Uploads artifacts
- Publishes release artifacts on tags matching `v*`

## License

GNU General Public License v3.0. See `LICENSE`.
