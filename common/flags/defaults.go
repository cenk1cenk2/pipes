package flags

// THIS EXISTS BECAUSE I AM MAINTAINING THIS LIBRARY FOR MORE THAN ONE PLACE
// The defaults for the flags are different for different use cases.
// So the logical flow would be to ignore the changes here since they are defaults for different usecases.

//revive:disable:line-length-limit

const (

	// docker.

	FLAG_DEFAULT_DOCKER_IMAGE_TAG_AS_LATEST = `[ "^tags/v?\\d.\\d.\\d$" ]`
	FLAG_DEFAULT_DOCKER_IMAGE_SANITIZE_TAGS = `[ { "condition": "([^/]*)/(.*)", "template": "{{ index $ 1 | to_upper_case }}_{{ index $ 2 }}" } ]`

	// node.

	FLAG_DEFAULT_NODE_PACKAGE_MANAGER = "yarn"

	// select-env.

	FLAG_DEFAULT_ENVIRONMENT_CONDITIONS = `[ { "condition": "^tags/v?\\d.\\d.\\d$", "environment": "production" }, { "condition": "^tags/v?\\d.\\d.\\d-.*\\.\\d$", "environment": "stage" }, { "condition" :"^heads/main$", "environment": "develop" }, { "condition": "^heads/master$", "environment": "develop" } ]`
)
