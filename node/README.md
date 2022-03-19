# gitlab-pipes-node

# Description

Login to NPM registry, installs dependencies and builds a Node application.

# Commands

## Root

```text
NAME:
   node-build - Pipe for installing node.js dependencies and building node.js applications on CI/CD.

USAGE:
   node-build [global options] command [command options] [arguments...]

VERSION:
   latest

DESCRIPTION:
   Pipe for installing node.js dependencies and building node.js applications on CI/CD.

COMMANDS:
   login
   install
   build
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --utils.ci value     Indicates this is running inside a CI/CD environment to act accordingly. [$CI]
   --utils.debug value  Set the log level debug for the application. [$DEBUG, $PLUGIN_DEBUG]
   --utils.log value    Define the log level for the application. (default: "info") [$LOG_LEVEL, $PLUGIN_LOG_LEVEL]
   --help, -h           show help (default: false)
   --version, -v        print the version (default: false)
```

## Login

```text
NAME:
   node-build login -

USAGE:
   node-build login [command options] [arguments...]

OPTIONS:
   --utils.ci value              Indicates this is running inside a CI/CD environment to act accordingly. [$CI]
   --utils.debug value           Set the log level debug for the application. [$DEBUG, $PLUGIN_DEBUG]
   --utils.log value             Define the log level for the application. (default: "info") [$LOG_LEVEL, $PLUGIN_LOG_LEVEL]
   --node.package_manager value  Preferred Package manager for nodejs. (default: "npm") [$NODE_PACKAGE_MANAGER]
   --npm.login value             npm registries to login to. Format: json({username: string, password: string, registry?: string, useHttps?: boolean}[]) [$NPM_LOGIN]
   --npm.npmrc_file value        .npmrc file to use. (default: ".npmrc") [$NPM_NPMRC_FILE]
   --help, -h                    show help (default: false)
```

## Install

```text
NAME:
   node-build install -

USAGE:
   node-build install [command options] [arguments...]

OPTIONS:
   --utils.ci value              Indicates this is running inside a CI/CD environment to act accordingly. [$CI]
   --utils.debug value           Set the log level debug for the application. [$DEBUG, $PLUGIN_DEBUG]
   --utils.log value             Define the log level for the application. (default: "info") [$LOG_LEVEL, $PLUGIN_LOG_LEVEL]
   --node.package_manager value  Preferred Package manager for nodejs. (default: "npm") [$NODE_PACKAGE_MANAGER]
   --node.install_cwd value      Install CWD for nodejs. (default: ".") [$NODE_INSTALL_CWD]
   --node.use_lock_file          Whether to use lock file or not. (default: true) [$NODE_INSTALL_USE_LOCK_FILE]
   --npm.login value             npm registries to login to. Format: json({username: string, password: string, registry?: string, useHttps?: boolean}[]) [$NPM_LOGIN]
   --npm.npmrc_file value        .npmrc file to use. (default: ".npmrc") [$NPM_NPMRC_FILE]
   --help, -h                    show help (default: false)
```

## Build

```text
NAME:
   node-build build -

USAGE:
   node-build build [command options] [arguments...]

OPTIONS:
   --utils.ci value                           Indicates this is running inside a CI/CD environment to act accordingly. [$CI]
   --utils.debug value                        Set the log level debug for the application. [$DEBUG, $PLUGIN_DEBUG]
   --utils.log value                          Define the log level for the application. (default: "info") [$LOG_LEVEL, $PLUGIN_LOG_LEVEL]
   --node.package_manager value               Preferred Package manager for nodejs. (default: "npm") [$NODE_PACKAGE_MANAGER]
   --git.branch value                         Source control management branch. [$CI_COMMIT_REF_NAME]
   --git.tag value                            Source control management tag. [$CI_COMMIT_TAG]
   --node.build_script value                  package.json script for building operation. (default: "build") [$NODE_BUILD_SCRIPT]
   --node.build_script_args value             package.json script arguments for building operation. [$NODE_BUILD_SCRIPT_ARGS]
   --node.build_cwd value                     Working directory for build operation. (default: ".") [$NODE_BUILD_CWD]
   --node.build_environment_files value       Yaml files to inject to build. [$NODE_BUILD_ENVIRONMENT_FILES]
   --node.build_environment_conditions value  Tagging regex patterns to match. json({ [name: string]: RegExp }) (default: "{ \"production\": \".*\" }") [$NODE_BUILD_ENVIRONMENT_CONDITIONS]
   --node.build_environment_fallback value    Fallback, if it does not match any conditions. Defaults to current branch name. (default: "staging") [$NODE_BUILD_ENVIRONMENT_FALLBACK]
   --help, -h                                 show help (default: false)
```
