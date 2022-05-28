<!-- clidocs -->

# NAME

gitlab-pipes-node - Pipe for installing node.js dependencies and building node.js applications on CI/CD.

# SYNOPSIS

gitlab-pipes-node

```
[--ci]
[--debug]
[--help|-h]
[--log-level]=[value]
[--version|-v]
```

# DESCRIPTION

Pipe for installing node.js dependencies and building node.js applications on CI/CD.

**Usage**:

```
gitlab-pipes-node [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--ci**: Sets whether this is running inside a CI/CD environment.

**--debug**: Enable debugging for the application.

**--help, -h**: show help

**--log-level**="": Define the log level for the application. (default: info)

**--version, -v**: print the version

# COMMANDS

## login

**--ci**: Sets whether this is running inside a CI/CD environment.

**--debug**: Enable debugging for the application.

**--log-level**="": Define the log level for the application. (default: info)

**--node.package_manager**="": Preferred Package manager for nodejs. (default: yarn)

**--npm.login**="": npm registries to login to. Format: json({username: string, password: string, registry?: string, useHttps?: boolean}[])

**--npm.npmrc**="": Pass direct contents of the NPMRC file.

**--npm.npmrc_file**="": .npmrc file to use. (default: [.npmrc])

## install

**--ci**: Sets whether this is running inside a CI/CD environment.

**--debug**: Enable debugging for the application.

**--log-level**="": Define the log level for the application. (default: info)

**--node.install_cwd**="": Install CWD for nodejs. (default: .)

**--node.package_manager**="": Preferred Package manager for nodejs. (default: yarn)

**--node.use_lock_file**: Whether to use lock file or not.

**--npm.login**="": npm registries to login to. Format: json({username: string, password: string, registry?: string, useHttps?: boolean}[])

**--npm.npmrc**="": Pass direct contents of the NPMRC file.

**--npm.npmrc_file**="": .npmrc file to use. (default: [.npmrc])

## build

**--ci**: Sets whether this is running inside a CI/CD environment.

**--debug**: Enable debugging for the application.

**--git.branch**="": Source control management branch.

**--git.tag**="": Source control management tag.

**--log-level**="": Define the log level for the application. (default: info)

**--node.build_cwd**="": Working directory for build operation. (default: .)

**--node.build_environment_conditions**="": Tagging regex patterns to match. json({ [name: string]: RegExp }) (default: { "production": "^v\\d*\\.\\d*\\.\\d*$", "stage": "^v\\d*\\.\\d*\\.\\d*-.\*$" })

**--node.build_environment_fallback**="": Fallback, if it does not match any conditions. Defaults to current branch name. (default: develop)

**--node.build_environment_files**="": Yaml files to inject to build. (default: [])

**--node.build_script**="": package.json script for building operation. (default: build)

**--node.build_script_args**="": package.json script arguments for building operation.

**--node.package_manager**="": Preferred Package manager for nodejs. (default: yarn)

## help, h

Shows a list of commands or help for one command

<!-- clidocsstop -->
