<!-- clidocs -->

# NAME

markdown-toc - Finds the markdown files and adds TOC to them.

# SYNOPSIS

markdown-toc

```
[--ci]
[--debug]
[--help|-h]
[--log-level]=[value]
[--markdown-toc.arguments]=[value]
[--markdown-toc.pattern]=[value]
[--version|-v]
```

# DESCRIPTION

Finds the markdown files and adds TOC to them.

**Usage**:

```
markdown-toc [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--ci**: Sets whether this is running inside a CI/CD environment.

**--debug**: Enable debugging for the application.

**--help, -h**: show help

**--log-level**="": Define the log level for the application. (default: info)

**--markdown-toc.arguments**="": Pass in the arguments for markdown-toc. (default: --bullets='-')

**--markdown-toc.pattern**="": Pattern for markdown. Supports multiple patterns with comma-separated values. (default: [README.md])

**--version, -v**: print the version


# COMMANDS

## help, h

Shows a list of commands or help for one command

<!-- clidocsstop -->
