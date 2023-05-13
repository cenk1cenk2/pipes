package setup

var (
	PackageManagers = AvailablePackageManagerCommands{
		"yarn": {
			Install:         []string{"install"},
			InstallWithLock: []string{"install", "--frozen-lock-file"},
			Run:             []string{"run"},
			RunDelimitter:   []string{},
			Add:             []string{"add"},
			Global:          []string{"global"},
			Cache:           []string{"--prefer-offline", "--cache-folder"},
			Version:         []string{"--version"},
		},

		"npm": {
			Install:         []string{"i", "--unsafe-perm"},
			InstallWithLock: []string{"ci", "--unsafe-perm"},
			Run:             []string{"run"},
			RunDelimitter:   []string{"--"},
			Add:             []string{"install", "--unsafe-perm"},
			Global:          []string{"-g"},
			Cache:           []string{"--prefer-offline", "--cache"},
			Version:         []string{"--version"},
		},

		"pnpm": {
			Install:         []string{"i", "--unsafe-perm"},
			InstallWithLock: []string{"i", "--frozen-lockfile"},
			Run:             []string{"run"},
			RunDelimitter:   []string{},
			Add:             []string{"add"},
			Global:          []string{"-g"},
			Cache:           []string{"--prefer-offline", "--store-dir"},
			Version:         []string{"--version"},
		},
	}
)

const (
	CONTAINER_USER  = "root"
	CONTAINER_GROUP = "root"
)
