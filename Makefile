# Go parameters
GOCMD=$(GOROOT)/bin/go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

# Cliapp
APPNAME = ratelimiter
APPPATH = ./cmd/$(APPNAME)
APPMAIN = $(APPPATH)/main.go

# Others
LINTER=golangci-lint

all: run

build: clean
	@echo "Build application $(APPNAME)"
	$(GOBUILD) -o $(APPPATH)/$(APPNAME) -v $(APPPATH)

clean:
	rm -f $(APPPATH)/$(APPNAME)

run: build
	@echo "Run $(APPNAME)"
	$(APPPATH)/$(APPNAME) -n=$(n) -x=$(x)

lint:
	@echo "Run linter  $(LINTER)"
	$(LINTER) run

test:
	@echo "Run tests $(APPNAME)"
	$(GOTEST) ./... -v

godoc:    
	godoc -http=:6060  -goroot=$(GOPATH)
