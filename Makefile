BINARY := bin/pomodoro
PKG := ./cmd/pomodoro
GO := go

.PHONY: all build repro-build install test fmt vet clean run

all: build

build:
	@mkdir -p bin
	$(GO) build -o $(BINARY) $(PKG)

repro-build:
	@mkdir -p bin
	# reproducible-ish build
	SOURCE_DATE_EPOCH=1700000000 CGO_ENABLED=0 $(GO) build -trimpath -buildvcs=false -ldflags="-s -w" -o $(BINARY) $(PKG)

install:
	$(GO) install $(PKG)

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	-rm -f $(BINARY)

run: build
	./$(BINARY)
