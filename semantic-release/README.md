<!-- clidocs -->

# NAME

semantic-release - Releases applications through semantic-release library.

# SYNOPSIS

semantic-release

```
[--ci]
[--debug]
[--help|-h]
[--log-level]=[value]
[--packages.apk]=[value]
[--packages.node]=[value]
[--semantic_release.dry_run]
[--semantic_release.run_multi]
[--version|-v]
```

# DESCRIPTION

Releases applications through semantic-release library.

**Usage**:

```
semantic-release [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--ci**: Sets whether this is running inside a CI/CD environment.

**--debug**: Enable debugging for the application.

**--help, -h**: show help

**--log-level**="": Define the log level for the application. (default: info)

**--packages.apk**="": APK applications to install before running semantic-release. (default: [])

**--packages.node**="": Node packages to install before running semantic-release. (default: [])

**--semantic_release.dry_run**: Node packages to install before running semantic-release.

**--semantic_release.run_multi**: Uses @qiwi/multi-semantic-release package to do a workspace release.

**--version, -v**: print the version

# COMMANDS

## help, h

Shows a list of commands or help for one command

<!-- clidocsstop -->
