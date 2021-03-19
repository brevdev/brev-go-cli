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

darwin_homebrew:
	env GOOS=darwin GOARCH=arm64 $(BUILDCMD) \
			-o $(PATH_BIN)/homebrew/arm64_big_sur/$(BIN_NAME)/$(BIN_VERSION)/bin/$(BIN_NAME) -v $(PATH_MAIN)
	env CGO_CFLAGS="-mmacosx-version-min=11.2" CGO_LDFLAGS="-mmacosx-version-min=11.2" GOOS=darwin GOARCH=amd64 $(BUILDCMD) \
			-o $(PATH_BIN)/homebrew/big_sur/$(BIN_NAME)/$(BIN_VERSION)/bin/$(BIN_NAME) -v $(PATH_MAIN)
	env CGO_CFLAGS="-mmacosx-version-min=10.15" CGO_LDFLAGS="-mmacosx-version-min=10.15" GOOS=darwin GOARCH=amd64 $(BUILDCMD) \
			-o $(PATH_BIN)/homebrew/catalina/$(BIN_NAME)/$(BIN_VERSION)/bin/$(BIN_NAME) -v $(PATH_MAIN)

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

dist-homebrew: darwin_homebrew
	mkdir -p $(PATH_DIST)/homebrew
	tar -C $(PATH_BIN)/homebrew/arm64_big_sur -czf $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).arm64_big_sur.bottle.tar.gz $(BIN_NAME)/$(BIN_VERSION)/bin/$(BIN_NAME)
	tar -C $(PATH_BIN)/homebrew/big_sur -czf $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).big_sur.bottle.tar.gz $(BIN_NAME)/$(BIN_VERSION)/bin/$(BIN_NAME)
	tar -C $(PATH_BIN)/homebrew/catalina -czf $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).catalina.bottle.tar.gz $(BIN_NAME)/$(BIN_VERSION)/bin/$(BIN_NAME)
	shasum -a 256 $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).arm64_big_sur.bottle.tar.gz | awk '{print $$1}' > $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).arm64_big_sur.bottle.tar.gz.sha256
	shasum -a 256 $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).big_sur.bottle.tar.gz | awk '{print $$1}' > $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).big_sur.bottle.tar.gz.sha256
	shasum -a 256 $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).catalina.bottle.tar.gz | awk '{print $$1}' > $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).catalina.bottle.tar.gz.sha256
	@echo "\nsha256 arm64_big_sur:"
	@cat $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).arm64_big_sur.bottle.tar.gz.sha256
	@echo "\nsha256 big_sur:"
	@cat $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).big_sur.bottle.tar.gz.sha256
	@echo "\nsha256 catalina:"
	@cat $(PATH_DIST)/homebrew/$(BIN_NAME)-$(BIN_VERSION).catalina.bottle.tar.gz.sha256

clean:
	$(GOCLEAN)
	rm -rf $(PATH_BIN)
