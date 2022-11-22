package setup

type (
	PackageManagerCommands struct {
		Install         []string
		InstallWithLock []string
		Run             []string
		RunDelimitter   []string
		Add             []string
		Global          []string
	}

	AvailablePackageManagerCommands map[string]PackageManagerCommands

	PackageManager struct {
		Exe      string
		Commands PackageManagerCommands
	}
)
