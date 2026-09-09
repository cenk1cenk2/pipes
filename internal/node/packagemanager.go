package node

type (
	// PackageManagerCommands is how one package manager spells each operation
	// the pipes run through it.
	PackageManagerCommands struct {
		Install         []string
		InstallWithLock []string
		Run             []string
		RunDelimiter    []string
		Add             []string
		Global          []string
		Cache           []string
		Version         []string
	}

	AvailablePackageManagerCommands map[string]PackageManagerCommands

	// PackageManager is the resolved package manager a task builds its commands
	// out of.
	PackageManager struct {
		Exe      string
		Commands PackageManagerCommands
	}
)

var PackageManagers = AvailablePackageManagerCommands{
	"yarn": {
		Install:         []string{"install"},
		InstallWithLock: []string{"install", "--frozen-lock-file"},
		Run:             []string{"run"},
		RunDelimiter:    []string{},
		Add:             []string{"add"},
		Global:          []string{"global"},
		Cache:           []string{"--prefer-offline", "--cache-folder"},
		Version:         []string{"--version"},
	},

	"npm": {
		Install:         []string{"i", "--unsafe-perm"},
		InstallWithLock: []string{"ci", "--unsafe-perm"},
		Run:             []string{"run"},
		RunDelimiter:    []string{"--"},
		Add:             []string{"install", "--unsafe-perm"},
		Global:          []string{"-g"},
		Cache:           []string{"--prefer-offline", "--cache"},
		Version:         []string{"--version"},
	},

	"pnpm": {
		Install:         []string{"i", "--unsafe-perm"},
		InstallWithLock: []string{"i", "--frozen-lockfile"},
		Run:             []string{"run"},
		RunDelimiter:    []string{},
		Add:             []string{"add"},
		Global:          []string{"-g"},
		Cache:           []string{"--prefer-offline", "--store-dir"},
		Version:         []string{"--version"},
	},
}
