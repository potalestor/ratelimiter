# Go parameters
GOCMD=$(GOROOT)/bin/go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

# Cliapp
APPNAME = ratelimiter
APPPATH = ./cmd/$(APPNAME)
APPMAIN = $(APPPATH)/main.go


all: run

build: clean
	@echo "Build application $(APPNAME)"
	$(GOBUILD) -o $(APPPATH)/$(APPNAME) -v $(APPPATH)

clean:
	rm -f $(APPPATH)/$(APPNAME)

run: build
	@echo "Run $(APPNAME)"
	$(APPPATH)/$(APPNAME) -n=$(n) -x=$(x)

linter: 
	golangci-lint run

test: 
	go test -v ./... -short	

godoc:    
	godoc -http=:6060  -goroot=$(GOPATH)
