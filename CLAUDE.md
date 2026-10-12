## Gotchas

- Do not write flag source-chain (env alias resolution) specs: resolving env vars through the flag chain is urfave/cli's behavior, already covered upstream. Tests cover pipe task logic only.
- A pipe image publishes only from its `<pipe>@x.y.z` tag pipeline, and the default branch builds no pipe. A change ships only when a releasing commit touches the pipe's directory, so a build matrix edit in `<pipe>/.gitlab-ci.yml` committed as `ci:` releases nothing.
- Global variables of the root pipeline reach a child pipeline expanded in the parent, where `$CI_PROJECT_DIR` is still empty, so every trigger job sets `inherit: variables: false` and the child defines its own paths.
- `vendor/` is never cached, only handed on as the artifact of `install`: a tree restored from an older dependency graph breaks go.
