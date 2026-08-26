package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func RenderOverlays(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("build").
		ShouldDisable(func(t *Task) bool {
			return len(deps.Tool.Overlays) == 0
		}).
		Set(func(t *Task) error {
			if P.Output != "" {
				if err := os.MkdirAll(P.Output, 0o755); err != nil {
					return fmt.Errorf("Can not create Kustomize output directory: %s: %w", P.Output, err)
				}
			}

			for _, overlay := range deps.Tool.Overlays {
				t.CreateSubtask(overlay).
					Set(func(t *Task) error {
						result := renderOverlay(overlay, P)

						t.Lock.Lock()
						C.Results = append(C.Results, result)
						t.Lock.Unlock()

						if result.Err != nil {
							return fmt.Errorf("Can not build Kustomize overlay: %s: %w", overlay, result.Err)
						}

						t.Log.Infof("Built Kustomize overlay: %s (%d resources)", overlay, result.DocCount)

						if P.Output != "" {
							if err := writeOverlay(overlay, result); err != nil {
								return err
							}
						}

						return nil
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

func writeOverlay(overlay string, result OverlayResult) error {
	name := strings.ReplaceAll(filepath.Clean(overlay), string(os.PathSeparator), "_")
	name = strings.Trim(name, "._")
	if name == "" {
		name = "root"
	}

	output := filepath.Join(P.Output, fmt.Sprintf("%s.yaml", name))
	if err := os.WriteFile(output, result.Yaml, 0o644); err != nil {
		return fmt.Errorf("Can not write rendered Kustomize overlay: %s: %w", output, err)
	}

	return nil
}
