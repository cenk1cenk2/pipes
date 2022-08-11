GO_VERSION=1.19

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run .
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_VENDOR=$(GO_CMD) mod vendor

GO_OPTION_C=0

install:
	$(GO_VENDOR)

update:
	$(GO_GET) -u all
	$(GO_VENDOR)
	$(GO_CMD) mod tidy -compat=$(GO_VERSION)

lint:
	CGO_ENABLED=$(GO_OPTION_C)	golangci-lint run ./...

tidy:
	$(GO_CMD) mod tidy -compat=$(GO_VERSION)
	$(GO_VENDOR)
