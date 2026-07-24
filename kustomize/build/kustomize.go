package build

import (
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func renderOverlay(root string, p *Pipe) OverlayResult {
	opts := krusty.MakeDefaultOptions()

	if p.LoadRestrictor == "none" {
		opts.LoadRestrictions = types.LoadRestrictionsNone
	}

	if p.EnableHelm {
		opts.PluginConfig = types.EnabledPluginConfig(types.BploUseStaticallyLinked)
		opts.PluginConfig.HelmConfig.Command = p.HelmCommand
		opts.PluginConfig.HelmConfig.KubeVersion = p.KubeVersion
	}

	m, err := krusty.MakeKustomizer(opts).Run(filesys.MakeFsOnDisk(), root)
	if err != nil {
		return OverlayResult{Overlay: root, Err: err}
	}

	y, err := m.AsYaml()

	return OverlayResult{Overlay: root, Yaml: y, DocCount: len(m.Resources()), Err: err}
}
