package pipe

type (
	PackageManagerCommands struct {
		Install         []string
		InstallWithLock []string
		Run             []string
		RunDelimitter   []string
	}

	AvailablePackageManagerCommands map[string]PackageManagerCommands

	PackageManager struct {
		Exe      string
		Commands PackageManagerCommands
	}
)
