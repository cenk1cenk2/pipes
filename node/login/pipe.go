// Package login writes the npmrc the package manager reads its credentials back from.
package login

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

var Login = &node.Login{}

func New(p *Plumber) *TaskList {
	return node.LoginTaskList(p, Login)
}
