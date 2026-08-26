package environment

import (
	"github.com/joho/godotenv"
)

// WriteFile writes the variables to path as a dotenv file, for the jobs that
// read the selected environment back through their own shell.
func WriteFile(path string, vars map[string]string) error {
	return godotenv.Write(vars, path)
}
