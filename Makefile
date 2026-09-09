APP_NAME := wg-cloudflare-dns
VERSION ?= $(shell git describe --tags --always --dirty)
DEB_VERSION := $(shell if echo "$(VERSION)" | grep -Eq '^[0-9]'; then echo "$(VERSION)"; else echo "0.0.0+$(VERSION)"; fi)
BUILD_DIR := build
BINARY := $(BUILD_DIR)/$(APP_NAME)
ARCH := amd64
DEB_ROOT := $(BUILD_DIR)/deb/$(APP_NAME)_$(DEB_VERSION)_$(ARCH)

.PHONY: build test package-deb clean

build:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) .

test:
	go test ./...

package-deb: build
	rm -rf $(DEB_ROOT)
	mkdir -p $(DEB_ROOT)/DEBIAN
	mkdir -p $(DEB_ROOT)/usr/local/bin
	mkdir -p $(DEB_ROOT)/etc/wg-cloudflare-dns
	mkdir -p $(DEB_ROOT)/lib/systemd/system
	cp $(BINARY) $(DEB_ROOT)/usr/local/bin/$(APP_NAME)
	cp config.yaml $(DEB_ROOT)/etc/wg-cloudflare-dns/config.yaml
	cp packaging/$(APP_NAME).service $(DEB_ROOT)/lib/systemd/system/$(APP_NAME).service
	printf "Package: $(APP_NAME)\nVersion: $(DEB_VERSION)\nSection: net\nPriority: optional\nArchitecture: $(ARCH)\nMaintainer: wg-cloudflare-dns maintainers\nDescription: WireGuard peer to Cloudflare DNS updater\n" > $(DEB_ROOT)/DEBIAN/control
	dpkg-deb --build $(DEB_ROOT)

clean:
	rm -rf $(BUILD_DIR)
