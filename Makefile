BIN_NAME?=brev
BIN_VERSION?=0.1.0

GOCMD=GO

GOBUILD=$(GOCMD) build
GOFMT=$(GOCMD) fmt
GOCLEAN=$(GOCMD) clean

PATH_BIN=bin
PATH_PROJECT=github.com/brevdev/brev-go-cli
PATH_MAIN=$(PATH_PROJECT)/cmd

FIELD_VERSION=$(PATH_PROJECT)/internal/config.Version

BUILDCMD=$(GOBUILD) -ldflags "-X $(FIELD_VERSION)=$(BIN_VERSION)"

build: linux darwin

linux:
	env GOOS=linux GOARCH=amd64 $(BUILDCMD) \
			-o $(PATH_BIN)/nix/$(BIN_NAME) -v $(PATH_MAIN)

darwin:
	env GOOS=darwin GOARCH=amd64 $(BUILDCMD) \
			-o $(PATH_BIN)/osx/$(BIN_NAME) -v $(PATH_MAIN)

fmt:
	$(GOFMT) ./...

clean:
	$(GOCLEAN)
	rm -rf $(PATH_BIN)