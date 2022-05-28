package setup

var (
	PackageManagers = AvailablePackageManagerCommands{
		"yarn": {
			Install:         []string{"install"},
			InstallWithLock: []string{"install", "--frozen-lock-file"},
			Run:             []string{"run"},
			RunDelimitter:   []string{},
		},

		"npm": {
			Install:         []string{"i", "--unsafe-perm"},
			InstallWithLock: []string{"ci", "--unsafe-perm"},
			Run:             []string{"run"},
			RunDelimitter:   []string{"--"},
		},
	}
)

const (
	CONTAINER_USER  = "root"
	CONTAINER_GROUP = "root"
)
