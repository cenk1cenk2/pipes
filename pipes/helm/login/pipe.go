package login

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

var P = &registry.Credentials{}

func New(p *plumber.Plumber) *plumber.TaskList {
	return registry.LoginTaskList(p, P, "helm", "registry", "login")
}
