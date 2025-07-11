package setup

type (
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

	PackageManager struct {
		Exe      string
		Commands PackageManagerCommands
	}
)
