# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVENDOR=$(GOCMD) mod vendor

install:
	$(GOVENDOR)

update:
	$(GOGET) -u all
	$(GOVENDOR)
	$(GOCMD) mod tidy -compat=1.17
