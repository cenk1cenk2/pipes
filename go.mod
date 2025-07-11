module gitlab.kilic.dev/devops/pipes

go 1.24.4

require (
	github.com/bmatcuk/doublestar/v4 v4.8.1
	github.com/cenk1cenk2/plumber/v6 v6.0.0-00010101000000-000000000000
	github.com/docker/docker v27.5.1+incompatible
	github.com/ekalinin/github-markdown-toc.go v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/nochso/gomd v0.0.0-20160625161351-1785d26cc410
	github.com/sirupsen/logrus v1.9.3
	github.com/urfave/cli/v3 v3.3.8
	gitlab.kilic.dev/libraries/go-utils/v2 v2.1.3
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/workanator/go-floc/v3 v3.0.1 // indirect
	gitlab.kilic.dev/libraries/go-broadcaster v1.1.3 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gotest.tools/v3 v3.2.0 // indirect
)

replace github.com/cenk1cenk2/plumber/v6 => ../plumber
