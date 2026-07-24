package build

import (
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// renderOverlay builds a single overlay through the Kustomize library. The
// enable-helm and load-restrictor defaults mirror ArgoCD's global
// kustomize.buildOptions (--enable-helm --load-restrictor LoadRestrictionsNone),
// so overlays validate here the same way ArgoCD renders them in the cluster.
func renderOverlay(root string, p *Pipe) OverlayResult {
	opts := krusty.MakeDefaultOptions()

	switch LoadRestrictor(p.LoadRestrictor) {
	case LoadRestrictorNone:
		opts.LoadRestrictions = types.LoadRestrictionsNone
	case LoadRestrictorRootOnly:
		opts.LoadRestrictions = types.LoadRestrictionsRootOnly
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
