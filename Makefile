## switchic Makefile
##
## Common targets:
##   make build           - compile binary into ./bin/switchic
##   make install         - build + install to PREFIX/bin (default /usr/local/bin, may sudo)
##   make user-install    - build + install to ~/.local/bin (no sudo)
##   make uninstall       - remove binary from PREFIX/bin
##   make test            - run go test ./...
##   make vet             - run go vet ./...
##   make fmt             - run gofmt -w on the tree
##   make clean           - remove ./bin
##   make version         - print VERSION used in builds

BINARY      ?= switchic
PREFIX      ?= /usr/local
BINDIR      ?= $(PREFIX)/bin
USER_BINDIR ?= $(HOME)/.local/bin
PKG          = github.com/Dandi-Pangestu/switchic
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS      = -s -w -X $(PKG)/cmd.Version=$(VERSION)

.PHONY: all build install user-install uninstall test vet fmt clean version help

all: build

help:
	@grep -E '^##' Makefile | sed 's/^## //'

version:
	@echo $(VERSION)

build:
	@mkdir -p bin
	@echo "Building $(BINARY) $(VERSION) -> bin/$(BINARY)"
	@CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' -o bin/$(BINARY) .

install: build
	@echo "Installing bin/$(BINARY) -> $(BINDIR)/$(BINARY)"
	@if [ -w "$(BINDIR)" ]; then \
		install -m 0755 bin/$(BINARY) $(BINDIR)/$(BINARY); \
	else \
		echo "Need elevated permissions to write to $(BINDIR); using sudo"; \
		sudo install -m 0755 bin/$(BINARY) $(BINDIR)/$(BINARY); \
	fi
	@echo ""
	@echo "Installed. Try: $(BINARY) --help"
	@case ":$$PATH:" in *":$(BINDIR):"*) ;; *) \
	  echo "Warning: $(BINDIR) is not on your PATH. Add it to ~/.zshrc or ~/.bashrc:"; \
	  echo "    export PATH=\"$(BINDIR):\$$PATH\"";; esac

user-install: build
	@mkdir -p $(USER_BINDIR)
	@echo "Installing bin/$(BINARY) -> $(USER_BINDIR)/$(BINARY)"
	@install -m 0755 bin/$(BINARY) $(USER_BINDIR)/$(BINARY)
	@echo ""
	@echo "Installed. Try: $(BINARY) --help"
	@case ":$$PATH:" in *":$(USER_BINDIR):"*) ;; *) \
	  echo "Warning: $(USER_BINDIR) is not on your PATH. Add it to ~/.zshrc or ~/.bashrc:"; \
	  echo "    export PATH=\"$(USER_BINDIR):\$$PATH\"";; esac

uninstall:
	@if [ -f "$(BINDIR)/$(BINARY)" ]; then \
		if [ -w "$(BINDIR)" ]; then rm -f $(BINDIR)/$(BINARY); \
		else sudo rm -f $(BINDIR)/$(BINARY); fi; \
		echo "Removed $(BINDIR)/$(BINARY)"; \
	fi
	@if [ -f "$(USER_BINDIR)/$(BINARY)" ]; then \
		rm -f $(USER_BINDIR)/$(BINARY); \
		echo "Removed $(USER_BINDIR)/$(BINARY)"; \
	fi

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -rf bin
