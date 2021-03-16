BIN_NAME?=brev
BIN_VERSION?=0.1.1
API_KEY_COTTER?=unknown

GOCMD=GO

GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOCLEAN=$(GOCMD) clean

PATH_BIN=bin
PATH_DIST=dist
PATH_PROJECT=github.com/brevdev/brev-go-cli
PATH_MAIN=$(PATH_PROJECT)/cmd

FIELD_VERSION=$(PATH_PROJECT)/internal/config.Version
FIELD_COTTER_API_KEY=$(PATH_PROJECT)/internal/config.CotterAPIKey

BUILDCMD=$(GOBUILD) -ldflags "-X $(FIELD_VERSION)=$(BIN_VERSION) -X $(FIELD_COTTER_API_KEY)=$(API_KEY_COTTER)"

build: linux darwin

linux:
	env GOOS=linux GOARCH=amd64 $(BUILDCMD) \
			-o $(PATH_BIN)/nix/$(BIN_NAME) -v $(PATH_MAIN)

darwin:
	env GOOS=darwin GOARCH=amd64 $(BUILDCMD) \
			-o $(PATH_BIN)/osx/$(BIN_NAME) -v $(PATH_MAIN)

test_unit:
	$(GOTEST) -v $(PATH_PROJECT)/cmd/...
	$(GOTEST) -v $(PATH_PROJECT)/internal/...

fmt:
	$(GOFMT) ./...

dist: build dist-linux dist-darwin

dist-linux:
	mkdir -p $(PATH_DIST)/nix
	tar -C $(PATH_BIN)/nix/ -czf $(PATH_DIST)/nix/brev-nix-64.tar.gz $(BIN_NAME)
	shasum -a 256 $(PATH_DIST)/nix/brev-nix-64.tar.gz | awk '{print $$1}' > $(PATH_DIST)/nix/brev-nix-64.tar.gz.sha256

dist-darwin:
	mkdir -p $(PATH_DIST)/osx
	tar -C $(PATH_BIN)/osx/ -czf $(PATH_DIST)/osx/brev-osx-64.tar.gz $(BIN_NAME)
	shasum -a 256 $(PATH_DIST)/osx/brev-osx-64.tar.gz | awk '{print $$1}' > $(PATH_DIST)/osx/brev-osx-64.tar.gz.sha256

clean:
	$(GOCLEAN)
	rm -rf $(PATH_BIN)
