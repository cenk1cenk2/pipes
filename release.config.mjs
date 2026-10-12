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
    [
      "@semantic-release/gitlab",
      {
        // every package of the workspace would comment on the same merge requests, and
        // a first release on all of them at once, which the instance cuts off midway.
        successCommentCondition: false
      }
    ]
  ]
}
