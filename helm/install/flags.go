package install

import (
	"github.com/urfave/cli/v3"
)

// The install command is driven entirely by the setup and login flags around
// it, but the file stays so that a flag of its own lands where every other
// command keeps them.
var Flags = []cli.Flag{}
