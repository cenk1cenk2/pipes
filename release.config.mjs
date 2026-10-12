export default {
  extends: "@cenk1cenk2/semantic-release-config",
  plugins: [
    [
      "@cenk1cenk2/semantic-release-config/presets/tag",
      {
        analysis: {
          releaseRules: {
            // a dependency bump inside a pipe changes what its image ships, so it
            // releases the pipe here, unlike the estate default.
            extend: [{ type: "chore", scope: "deps", release: "patch" }]
          }
        }
      }
    ],
    "@semantic-release/gitlab"
  ]
}
