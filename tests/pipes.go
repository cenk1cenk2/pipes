// Package tests is the conformance suite. It imports every pipe's command tree
// and asserts the conventions the pipes agreed on, so a rule that only lived in
// a review comment fails a build instead.
//
// Nothing imports this module, which is what lets it depend on all of them at
// once without closing a cycle.
package tests

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	buildahapp "gitlab.kilic.dev/devops/pipes/buildah/app"
	goapp "gitlab.kilic.dev/devops/pipes/go/app"
	helmapp "gitlab.kilic.dev/devops/pipes/helm/app"
	kustomizeapp "gitlab.kilic.dev/devops/pipes/kustomize/app"
	nodeapp "gitlab.kilic.dev/devops/pipes/node/app"
	pulumiapp "gitlab.kilic.dev/devops/pipes/pulumi/app"
	selectenvapp "gitlab.kilic.dev/devops/pipes/select-env/app"
	semanticreleaseapp "gitlab.kilic.dev/devops/pipes/semantic-release/app"
	terraformapp "gitlab.kilic.dev/devops/pipes/terraform/app"
	updatedockerhubreadmeapp "gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/app"
)

// VERSION is what the trees are built with. The specs only read the shape of a
// tree, so the value never reaches anything that cares what it is.
const VERSION = "test"

// pipe is one entry of the conformance table: the directory it lives in, the
// constants it publishes and the tree its app package builds.
type pipe struct {
	// Dir is the module directory, and the stem the image name and the manifest
	// entry are derived from.
	Dir         string
	Name        string
	Description string
	// New builds the tree. A nil plumber is enough for the shape, since a step
	// only reaches for one once the command runs.
	New func(p *plumber.Plumber) *ucli.Command
}

// Excluded are the directories that hold a Go module but no pipe. They are named
// here rather than skipped by a pattern, so adding one is a decision somebody
// writes down.
var Excluded = []string{
	// template is the scaffold a new pipe is copied from. It publishes no app
	// package and ships no image, so there is nothing to conform to.
	"template",
	// internal is the shared library the pipes are built out of.
	"internal",
	// tests is this module.
	"tests",
}

var pipes = []pipe{
	{"buildah", buildahapp.Name, buildahapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return buildahapp.New(p, VERSION, buildahapp.Options{})
	}},
	{"go", goapp.Name, goapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return goapp.New(p, VERSION, goapp.Options{})
	}},
	{"helm", helmapp.Name, helmapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return helmapp.New(p, VERSION, helmapp.Options{})
	}},
	{"kustomize", kustomizeapp.Name, kustomizeapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return kustomizeapp.New(p, VERSION, kustomizeapp.Options{})
	}},
	{"node", nodeapp.Name, nodeapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return nodeapp.New(p, VERSION, nodeapp.Options{})
	}},
	{"pulumi", pulumiapp.Name, pulumiapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return pulumiapp.New(p, VERSION, pulumiapp.Options{})
	}},
	{"select-env", selectenvapp.Name, selectenvapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return selectenvapp.New(p, VERSION, selectenvapp.Options{})
	}},
	{"semantic-release", semanticreleaseapp.Name, semanticreleaseapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return semanticreleaseapp.New(p, VERSION, semanticreleaseapp.Options{})
	}},
	{"terraform", terraformapp.Name, terraformapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return terraformapp.New(p, VERSION, terraformapp.Options{})
	}},
	{"update-docker-hub-readme", updatedockerhubreadmeapp.Name, updatedockerhubreadmeapp.Description, func(p *plumber.Plumber) *ucli.Command {
		return updatedockerhubreadmeapp.New(p, VERSION, updatedockerhubreadmeapp.Options{})
	}},
}

// Image is the container repository the pipe publishes to. Every pipe is
// prefixed, including the one whose command is not.
func (p pipe) Image() string {
	return "cenk1cenk2/pipe-" + p.Dir
}

// flagRef is a flag together with the command path it was reached through, so a
// failure names the subcommand a reader has to open.
type flagRef struct {
	Command string
	Flag    ucli.Flag
}

func (f flagRef) Name() string {
	return f.Flag.Names()[0]
}

func (f flagRef) Category() string {
	if c, ok := f.Flag.(ucli.CategorizableFlag); ok {
		return c.GetCategory()
	}

	return ""
}

func (f flagRef) EnvVars() []string {
	if d, ok := f.Flag.(ucli.DocGenerationFlag); ok {
		return d.GetEnvVars()
	}

	return nil
}

func (f flagRef) Visible() bool {
	if v, ok := f.Flag.(ucli.VisibleFlag); ok {
		return v.IsVisible()
	}

	return true
}

// visibleFlags walks the whole tree and returns every flag a user could be shown.
// A flag shared between subcommands is the same instance, so it comes back once
// per command it is reachable from.
func (p pipe) visibleFlags() []flagRef {
	refs := []flagRef{}

	var walk func(path string, c *ucli.Command)

	walk = func(path string, c *ucli.Command) {
		for _, f := range c.Flags {
			ref := flagRef{Command: path, Flag: f}

			if ref.Visible() {
				refs = append(refs, ref)
			}
		}

		for _, sub := range c.Commands {
			walk(path+" "+sub.Name, sub)
		}
	}

	walk(p.Name, p.New(nil))

	return refs
}
