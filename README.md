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
- **Plugin System**: Optional plugins for extended functionality (Nexus mirroring, etc.)

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

# Install multiple Maven versions at once
mvnenv install 3.8.8 3.9.6 4.0.0-alpha-13

# Install latest Maven version
mvnenv install latest

# List available versions from Apache archive
mvnenv install -l

# Install with options
mvnenv install -f 3.9.6              # Force reinstall if already exists
mvnenv install -s 3.9.6              # Skip if already exists (no error)
mvnenv install -c 3.9.6              # Clear cache before installing
mvnenv install -q 3.9.6              # Quiet mode (suppress output)
mvnenv install --offline 3.9.6       # Offline mode (Nexus only, no Apache fallback)

# Combine flags
mvnenv install -f -q 3.9.6           # Force + quiet
mvnenv install -c -s 3.8.8 3.9.6     # Clear + skip-existing + multiple versions

# Uninstall Maven version
mvnenv uninstall <version>

# List installed versions
mvnenv versions

# Check if current version is the latest
mvnenv status
```

#### Installation Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--list` | `-l` | List all available Maven versions |
| `--quiet` | `-q` | Suppress progress output (errors still shown) |
| `--force` | `-f` | Force reinstall even if version already exists |
| `--skip-existing` | `-s` | Skip installation if version exists (no error) |
| `--clear` | `-c` | Clear download cache before installing |
| `--offline` | | Offline mode: only use Nexus (fail if unavailable) |

### Version Selection

mvnenv uses a three-tier hierarchy to resolve which Maven version to use:

1. **Shell**: Set via `MVNENV_MAVEN_VERSION` environment variable (highest priority)
2. **Local**: Set via `.maven-version` file in current or parent directories
3. **Global**: Set via `mvnenv global` command (lowest priority)

#### Global Version

The global version provides a system-wide default Maven version used when no shell or local version is specified.

```bash
# Set global version
mvnenv global 3.9.4

# Show current global version
mvnenv global

# Unset global version
mvnenv global --unset
```

#### Local Version

Set a project-specific Maven version using a `.maven-version` file:

```bash
# Set local version (project-specific)
mvnenv local 3.8.6
```

#### Shell Version

Override version for the current shell session:

```bash
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

## Plugins

mvnenv-win supports optional plugins that can be enabled at build time for extended functionality.

### Mirror Plugin

Create internal Nexus mirrors by downloading Maven distributions from Apache and uploading to Nexus:

```bash
# Mirror all Maven versions to configured Nexus repository
mvnenv mirror

# Dry-run to see what would be mirrored
mvnenv mirror --dry-run

# Mirror only the 10 most recent versions
mvnenv mirror --max 10
```

**Note:** The mirror plugin requires:
- Nexus repository configured in config.yaml
- Nexus user with write permissions
- Plugin must be enabled at build time

### Building with Plugins

```bash
# Build without plugins (standard)
go build -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go

# Build with all plugins enabled
go build -tags "mirror" -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go

# Or use Makefile
make build          # Without plugins
make build-plugins  # With all plugins
make dist           # Create production distribution with plugins
```

See [PLUGINS.md](PLUGINS.md) for complete plugin documentation.

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
# Build standard binary (without plugins)
make build

# Build with all plugins enabled
make build-plugins

# Build shim executable
make build-shim

# Build everything (main + plugins + shim)
make build-all

# Create production distribution package
make dist

# Direct Go commands
go build -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go
go build -tags "mirror" -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go
```

### Running Tests

```bash
make test
```

### Makefile Targets

- `make build` - Build standard binary without plugins
- `make build-plugins` - Build binary with all plugins enabled
- `make build-shim` - Build shim executable
- `make build-all` - Build everything
- `make dist` - Create production distribution package with plugins
- `make dist-noplugin` - Create production distribution package without plugins
- `make clean` - Remove build artifacts
- `make test` - Run tests
- `make help` - Display available targets

## License

Apache License 2.0

## Acknowledgments

Inspired by pyenv-win and the pyenv project for Unix-like systems.
