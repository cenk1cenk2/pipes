// Package login writes the npmrc the package manager reads its credentials
// back from, so the commands that reach a registry compose it after setup.
package login

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

var Login = &node.Login{}

func New(p *Plumber) *TaskList {
	return node.LoginTaskList(p, Login)
}
