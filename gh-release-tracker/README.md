# gh-release-tracker

## Description

## Usage

```bash
NAME:
   gh-release-tracker - Tracks the given Github repositories latest tags and generates tag file accordingly for further processing.

USAGE:
   gh-release-tracker [global options] command [command options] [arguments...]

VERSION:
   latest

DESCRIPTION:
   Tracks the given Github repositories latest tags and generates tag file accordingly for further processing.

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --utils.ci value       Indicates this is running inside a CI/CD environment to act accordingly. [$CI]
   --utils.debug value    Set the log level debug for the application. [$DEBUG, $PLUGIN_DEBUG]
   --utils.log value      Define the log level for the application. (default: "info") [$LOG_LEVEL, $PLUGIN_LOG_LEVEL]
   --gh.token value       Github token for the API requests. [$GH_TOKEN, $GITHUB_TOKEN]
   --gh.repository value  Target repository to fetch the latest tag. [$GH_REPOSITORY, $GH_REPOSITORY]
   --gh.tags_file value   Read tags from a file. [$TAGS_FILE]
   --help, -h             show help (default: false)
   --version, -v          print the version (default: false)
```
