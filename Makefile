BIN_NAME?=brev
BIN_VERSION?=0.1.2
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

darwin-homebrew:
	env GOOS=darwin GOARCH=arm64 $(BUILDCMD) \
			-o $(PATH_BIN)/homebrew/$(BIN_NAME)-arm64_big_sur -v $(PATH_MAIN)
	env CGO_CFLAGS="-mmacosx-version-min=11.2" CGO_LDFLAGS="-mmacosx-version-min=11.2" GOOS=darwin GOARCH=amd64 $(BUILDCMD) \
			-o $(PATH_BIN)/homebrew/$(BIN_NAME)-big_sur -v $(PATH_MAIN)
	env CGO_CFLAGS="-mmacosx-version-min=10.15" CGO_LDFLAGS="-mmacosx-version-min=10.15" GOOS=darwin GOARCH=amd64 $(BUILDCMD) \
			-o $(PATH_BIN)/homebrew/$(BIN_NAME)-catalina -v $(PATH_MAIN)

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

dist-homebrew: darwin-homebrew
	mkdir -p $(PATH_DIST)/homebrew
	tar -C $(PATH_BIN)/homebrew/ -czf $(PATH_DIST)/homebrew/brev-homebrew-bundle.tar.gz .
	shasum -a 256 $(PATH_DIST)/homebrew/brev-homebrew-bundle.tar.gz | awk '{print $$1}' > $(PATH_DIST)/homebrew/brev-homebrew-bundle.tar.gz.sha256
	@echo "\nsha256:"
	@cat $(PATH_DIST)/homebrew/brev-homebrew-bundle.tar.gz.sha256

clean:
	$(GOCLEAN)
	rm -rf $(PATH_BIN)
