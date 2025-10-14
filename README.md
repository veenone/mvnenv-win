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
- **Nexus Repository Support**: Download from private Nexus repositories with authentication and custom SSL/TLS
- **Checksum Verification**: SHA-512 checksum verification for downloaded files

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/veenone/mvnenv-win.git
cd mvnenv-win

# Build both executables
go build -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go
go build -o bin/shim.exe cmd/shim/main.go

# Create mvnenv directory structure
mkdir -p %USERPROFILE%\.mvnenv\bin

# Copy executables to mvnenv directory
copy bin\mvnenv.exe %USERPROFILE%\.mvnenv\bin\
copy bin\shim.exe %USERPROFILE%\.mvnenv\bin\

# Add to PATH (PowerShell - run as Administrator or use User PATH)
$env:Path = "$env:USERPROFILE\.mvnenv\shims;$env:USERPROFILE\.mvnenv\bin;" + $env:Path
[Environment]::SetEnvironmentVariable("Path", "$env:USERPROFILE\.mvnenv\shims;$env:USERPROFILE\.mvnenv\bin;" + [Environment]::GetEnvironmentVariable("Path", "User"), "User")

# Or add manually:
# Add these directories to your PATH environment variable (in order):
#   1. %USERPROFILE%\.mvnenv\shims  (highest priority)
#   2. %USERPROFILE%\.mvnenv\bin
```

### Initial Setup

After installation, generate the shims for Maven commands:

```bash
# This creates shim executables that intercept Maven commands
mvnenv rehash
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

# Install latest Maven version
mvnenv install latest

# List available versions from Apache archive
mvnenv install -l

# Uninstall Maven version
mvnenv uninstall <version>

# List installed versions
mvnenv versions

# Check if current version is the latest
mvnenv status
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

# Find latest installed version
mvnenv latest
mvnenv latest 3.9        # Latest 3.9.x version

# Find latest available version
mvnenv latest --remote
mvnenv latest --remote 3.8

# Update version cache
mvnenv update

# Regenerate shims
mvnenv rehash

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
├── bin/            # mvnenv and shim executables
│   ├── mvnenv.exe
│   └── shim.exe
├── shims/          # Command shims (mvn.exe, mvnDebug.exe, etc.)
│   ├── mvn.exe
│   ├── mvn.cmd
│   ├── mvnDebug.exe
│   └── mvnDebug.cmd
├── cache/          # Downloaded Maven archives and version cache
│   ├── apache-maven-3.9.4-bin.zip
│   └── versions.json           # Cached list of available versions
├── config/         # Configuration files
│   └── config.yaml             # Global configuration
└── versions/       # Installed Maven versions
    ├── 3.8.6/
    ├── 3.9.4/
    └── ...
```

**Important:** The `shims` directory must be first in your PATH to intercept Maven commands.

## Version Resolution

When you run a Maven command, mvnenv resolves the version in this order:

1. Check `MVNENV_MAVEN_VERSION` environment variable
2. Look for `.maven-version` file in current directory and parent directories
3. Use global version from config.yaml
4. Error if no version is set

## Version Cache

mvnenv caches the list of available Maven versions from Apache archive to improve performance:

- Cache is automatically created when you run `mvnenv install -l`
- Cache is valid for 24 hours before being automatically refreshed
- Use `mvnenv update` to manually refresh the cache
- The `mvnenv latest --remote` command also uses the cache

This reduces network calls and speeds up version listing operations.

## Nexus Repository Integration

mvnenv-win supports downloading Maven distributions from private Nexus Repository Manager instances. This is useful for enterprise environments that require using internal repositories.

### Quick Setup

Edit your configuration file at `%USERPROFILE%\.mvnenv\config\config.yaml`:

```yaml
nexus:
  enabled: true
  base_url: "https://nexus.example.com/repository/maven-public"
  username: "your-username"
  password: "your-password"
```

### Self-Signed Certificates

For Nexus servers with self-signed certificates:

```yaml
nexus:
  enabled: true
  base_url: "https://nexus.internal.company.com/repository/maven-central"
  username: "myuser"
  password: "mypassword"
  tls:
    insecure_skip_verify: true
```

### Custom CA Certificate

For corporate CA certificates:

```yaml
nexus:
  enabled: true
  base_url: "https://nexus.internal.company.com/repository/maven-central"
  username: "myuser"
  password: "mypassword"
  tls:
    ca_file: "C:\\company\\ca\\root-ca.pem"
```

**Note:** When Nexus is configured, mvnenv will try to download from Nexus first and automatically fall back to Apache Maven archives if Nexus is unavailable.

See [NEXUS.md](NEXUS.md) for complete documentation on Nexus integration.

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

### Maven commands still use system Maven

If `mvn -version` shows the wrong version, your PATH may not be configured correctly:

```bash
# Check if shims directory is in PATH and has priority
echo %PATH%

# The shims directory should appear BEFORE any other Maven installation
# Example: C:\Users\YourName\.mvnenv\shims;C:\Program Files\Maven\...

# Regenerate shims
mvnenv rehash

# Verify shim is being used
where mvn
# Should show: C:\Users\YourName\.mvnenv\shims\mvn.exe
```

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

### Shims not working

```bash
# Regenerate shims after installing/uninstalling versions
mvnenv rehash

# Check if shim.exe exists
dir %USERPROFILE%\.mvnenv\bin\shim.exe

# Check if shims were created
dir %USERPROFILE%\.mvnenv\shims\
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
