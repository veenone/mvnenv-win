# mvnenv-win

Maven version manager for Windows - inspired by pyenv-win

## Overview

mvnenv-win is a command-line tool for managing multiple Apache Maven installations on Windows. It allows you to easily switch between different Maven versions for various projects without manual PATH updates or system-wide configuration changes.

## Features

- **Multi-Version Management**: Install and maintain multiple Maven versions simultaneously
- **Smart Version Selection**: Automatic version resolution using shell > local > global hierarchy
- **Simple CLI**: Familiar pyenv-style commands
- **Native Windows**: Pure Go implementation, no WSL or Cygwin required
- **Apache Maven Integration**: Direct downloads from official Apache Maven archives
- **Checksum Verification**: SHA-512 checksum verification for downloaded files

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/veenone/mvnenv-win.git
cd mvnenv-win

# Build
go build -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go

# Add to PATH
# Add C:\path\to\mvnenv-win\bin to your PATH environment variable
```

## Quick Start

```bash
# List available Maven versions
mvnenv install -l

# Install a specific Maven version
mvnenv install 3.9.4

# Set as global default
mvnenv global 3.9.4

# Check current version
mvnenv version

# List installed versions
mvnenv versions
```

## Usage

### Version Management

```bash
# Install Maven version
mvnenv install <version>

# List available versions from Apache archive
mvnenv install -l

# Uninstall Maven version
mvnenv uninstall <version>

# List installed versions
mvnenv versions
```

### Version Selection

mvnenv uses a three-tier hierarchy to resolve which Maven version to use:

1. **Shell**: Set via `MVNENV_MAVEN_VERSION` environment variable (highest priority)
2. **Local**: Set via `.maven-version` file in current or parent directories
3. **Global**: Set in `%USERPROFILE%\.mvnenv\config\config.yaml` (lowest priority)

```bash
# Set global version (system-wide default)
mvnenv global 3.9.4

# Set local version (project-specific)
mvnenv local 3.8.6

# Set shell version (current session only)
mvnenv shell 3.9.4
# Then set environment variable:
#   PowerShell: $env:MVNENV_MAVEN_VERSION = "3.9.4"
#   cmd.exe: set MVNENV_MAVEN_VERSION=3.9.4
```

### Utility Commands

```bash
# Show current Maven version and source
mvnenv version

# Show path to Maven executable
mvnenv which mvn

# List all available commands
mvnenv commands

# Show help
mvnenv help
mvnenv <command> --help
```

## Directory Structure

mvnenv-win uses the following directory structure in `%USERPROFILE%\.mvnenv\`:

```
.mvnenv/
├── bin/            # mvnenv executable
├── cache/          # Downloaded Maven archives
├── config/         # Configuration files
│   └── config.yaml # Global configuration
└── versions/       # Installed Maven versions
    ├── 3.8.6/
    ├── 3.9.4/
    └── ...
```

## Version Resolution

When you run a Maven command, mvnenv resolves the version in this order:

1. Check `MVNENV_MAVEN_VERSION` environment variable
2. Look for `.maven-version` file in current directory and parent directories
3. Use global version from config.yaml
4. Error if no version is set

## Requirements

- Windows 10 or later
- Go 1.21+ (for building from source)

## Examples

### Per-Project Maven Versions

```bash
# Project A requires Maven 3.8.6
cd /path/to/project-a
mvnenv local 3.8.6

# Project B requires Maven 3.9.4
cd /path/to/project-b
mvnenv local 3.9.4

# Each project now uses its specified version automatically
```

### Testing Different Maven Versions

```bash
# Temporarily test with Maven 3.9.4
mvnenv shell 3.9.4
$env:MVNENV_MAVEN_VERSION = "3.9.4"  # PowerShell
# Run your tests...

# Unset to return to normal resolution
$env:MVNENV_MAVEN_VERSION = ""
```

## Troubleshooting

### Version not found

```bash
# List available versions
mvnenv install -l

# Install the version
mvnenv install 3.9.4
```

### Check current version

```bash
# Show which version is active and where it's set from
mvnenv version

# Show path to Maven executable
mvnenv which mvn
```

## Development

### Building

```bash
# Build binary
go build -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go

# Or using Make
make build
```

### Running Tests

```bash
make test
```

## License

Apache License 2.0

## Acknowledgments

Inspired by pyenv-win and the pyenv project for Unix-like systems.
