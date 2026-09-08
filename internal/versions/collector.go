package versions

import (
	"fmt"
	"path"
	"regexp"
	"slices"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/sirupsen/logrus"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"
)

// Collector gathers the tags or the versions a pipe publishes under. Every source is
// optional: the task of a source left at its zero value disables itself. The pipe
// composes the sources itself, since the parent task that reports the result is also
// where it hangs whatever else has to run once the values are in.
type Collector struct {
	// Name prefixes every task the collector creates.
	Name string

	FromUser []string

	File       string
	FileStrict bool
	FileDir    string

	// LatestWhen are the patterns matched against References; the first one that
	// matches adds LatestValue to the result. A nil slice leaves the task out.
	LatestWhen  []string
	LatestValue string
	References  []string

	// Templates rewrite a value and Sanitize strips patterns out of it. The first
	// entry that matches wins and the rest are skipped.
	Templates []Match
	Sanitize  []Match

	// Format is the last step, where a pipe turns a processed value into whatever
	// it actually publishes under.
	Format func(string) string
}

// UserTask collects what the pipeline named on the command line.
func (c *Collector) UserTask(tl *TaskList, out *[]string) *Task {
	return tl.CreateTask(c.Name, "user").
		ShouldDisable(func(_ *Task) bool {
			return len(c.FromUser) == 0
		}).
		Set(func(t *Task) error {
			for _, v := range slices.Compact(slices.Clone(c.FromUser)) {
				if err := c.add(t, out, v); err != nil {
					return err
				}
			}

			return nil
		})
}

// FileTask collects what the file the pipeline pointed at names, which is how a
// job passes values to the one after it.
func (c *Collector) FileTask(tl *TaskList, out *[]string) *Task {
	return tl.CreateTask(c.Name, "file").
		ShouldDisable(func(_ *Task) bool {
			return c.File == ""
		}).
		Set(func(t *Task) error {
			values, err := tagsfile.Parse(t.Log, path.Join(c.FileDir, c.File), c.FileStrict)

			if err != nil {
				return err
			}

			for _, v := range values {
				t.CreateSubtask(v).
					Set(func(t *Task) error {
						return c.add(t, out, v)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

// LatestTask adds the latest value when a source control reference matches, and
// leaves itself out for a pipe that has no notion of one.
func (c *Collector) LatestTask(tl *TaskList, out *[]string) *Task {
	return tl.CreateTask(c.Name, "latest").
		ShouldDisable(func(_ *Task) bool {
			return c.LatestWhen == nil
		}).
		Set(func(t *Task) error {
			matched, err := git.MatchAny(c.LatestWhen, c.References)
			if err != nil {
				return fmt.Errorf("Can not process regular expression for latest tag: %w", err)
			}

			if matched < 0 {
				return nil
			}

			if err := c.add(t, out, c.LatestValue); err != nil {
				return err
			}

			t.Log.Infof(
				"Will tag as latest since the reference matches: %s -> %s",
				c.LatestValue,
				c.LatestWhen[matched],
			)

			return nil
		})
}

// The file source appends from parallel subtasks, which is why the append is
// taken under the task list lock rather than a lock of the collector's own.
func (c *Collector) add(t *Task, out *[]string, value string) error {
	processed, err := c.Process(t.Log, value)

	if err != nil {
		return err
	}

	t.Lock.Lock()
	*out = append(*out, processed)
	t.Lock.Unlock()

	return nil
}

// Process runs one raw value through the templates, the sanitizers and the pipe's
// own formatting, in that order.
func (c *Collector) Process(log *logrus.Entry, value string) (string, error) {
	value, err := c.template(log, value)

	if err != nil {
		return "", err
	}

	value, err = c.sanitize(log, value)

	if err != nil {
		return "", err
	}

	if value == "" {
		return value, fmt.Errorf("Can not add empty tag to list.")
	}

	if c.Format != nil {
		value = c.Format(value)
	}

	return value, nil
}

// A value that matches nothing is still rendered as a template of its own, which
// is what lets a pipeline pass a template in directly instead of a condition.
func (c *Collector) template(log *logrus.Entry, value string) (string, error) {
	for _, m := range c.Templates {
		re, err := regexp.Compile(m.Match)

		if err != nil {
			return "", fmt.Errorf("Can not compile tag template regular expression: %w", err)
		}

		log.Tracef("Trying to apply template to tag: %s with %v", value, re.String())

		matches := re.FindStringSubmatch(value)

		if matches == nil {
			continue
		}

		log.Debugf("Applying template since condition matched for given tag: %s -> %s with %v", value, re.String(), matches)

		return InlineTemplate(m.Template, matches)
	}

	return InlineTemplate[any](value, nil)
}

func (c *Collector) sanitize(log *logrus.Entry, value string) (string, error) {
	for _, m := range c.Sanitize {
		re, err := regexp.Compile(m.Match)

		if err != nil {
			return "", fmt.Errorf("Can not compile sanitize regular expression: %w", err)
		}

		log.Debugf("Trying to sanitize tag: %s with %v", value, re.String())

		matches := re.FindStringSubmatch(value)

		if matches == nil {
			continue
		}

		log.Debugf("Sanitizing since condition matched for given tag: %s -> %s with %v", value, re.String(), matches)

		return InlineTemplate(m.Template, matches)
	}

	return value, nil
}
