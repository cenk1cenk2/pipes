package setup

import (
	"fmt"
	"path/filepath"

	"github.com/cenk1cenk2/plumber/v6"
)

func GoEnv(tl *plumber.TaskList) *plumber.Task {
	return tl.CreateTask("env").
		Set(func(t *plumber.Task) error {
			if P.Cache != "" {
				cache, err := filepath.Abs(P.Cache)
				if err != nil {
					return fmt.Errorf("Cannot get absolute path of cache dir: %s -> %w", P.Cache, err)
				}

				C.Env["GOPATH"] = cache
				C.Env["GOCACHE"] = filepath.Join(cache, "go-build")
				C.Env["GOLANGCI_LINT_CACHE"] = filepath.Join(cache, "golangci-lint")
			}

			return nil
		})
}
