# mvnenv-win Plugins

mvnenv-win supports an optional plugin system that allows extending functionality through build-time enabled features.

## Available Plugins

### Mirror Plugin

The mirror plugin provides functionality to create and maintain an internal Nexus mirror of Maven distributions.

**Features:**
- Download all available Maven versions from Apache Maven archive
- Upload distributions to a configured Nexus repository
- Skip already mirrored versions for efficiency
- Dry-run mode for planning
- Limit number of versions to mirror

**Usage:**

```bash
# Mirror all Maven versions to Nexus
mvnenv mirror

# Dry-run to see what would be mirrored
mvnenv mirror --dry-run

# Mirror only the 10 most recent versions
mvnenv mirror --max 10

# Force re-upload of existing versions
mvnenv mirror --skip-existing=false
```

**Requirements:**
- Nexus repository must be configured in config.yaml
- Nexus user must have write permissions
- Sufficient disk space for temporary downloads

## Building with Plugins

Plugins are enabled at build time using Go build tags.

### Build without Plugins (Standard)

```bash
# Using Makefile
make build

# Using go directly
go build -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go
```

### Build with All Plugins Enabled

```bash
# Using Makefile
make build-plugins

# Using go directly
go build -tags "mirror" -ldflags "-X main.Version=$(cat VERSION)" -o bin/mvnenv.exe cmd/mvnenv/main.go
```

### Build Everything (Plugins + Shim)

```bash
make build-all
```

## Creating a Production Distribution

The Makefile includes a `dist` target that creates a complete production-ready distribution:

```bash
make dist
```

This will:
1. Clean previous builds
2. Build mvnenv.exe with all plugins enabled
3. Build shim.exe
4. Create a versioned distribution directory
5. Copy all necessary files (executables, docs, config examples)
6. Display installation instructions

The distribution will be created in `dist/mvnenv-<version>/` and includes:
- `bin/mvnenv.exe` (with plugins)
- `bin/shim.exe`
- `VERSION`
- `README.md`
- `SETUP.md`
- `NEXUS.md`
- `config/config.example.yaml`

## Plugin Architecture

mvnenv-win uses Go build tags to conditionally include plugin code:

1. **Plugin Interface**: Defined in `cmd/mvnenv/plugins/plugin.go`
2. **Plugin Implementation**: Each plugin is in its own subdirectory under `cmd/mvnenv/plugins/`
3. **Build Tag**: Plugin files use `// +build <tag>` to enable/disable compilation
4. **Stub Files**: Plugins include stub files with `// +build !<tag>` to prevent build errors when disabled
5. **Registration**: Plugins register themselves via `init()` functions when enabled

### Plugin Structure

```
cmd/mvnenv/plugins/
├── plugin.go           # Plugin interface and registry
└── mirror/             # Mirror plugin
    ├── mirror.go       # Plugin implementation (enabled with "mirror" tag)
    └── stub.go         # Empty stub (enabled without "mirror" tag)
```

### Adding a New Plugin

1. Create a new subdirectory under `cmd/mvnenv/plugins/`
2. Create the main plugin file with build tag: `// +build myplugin`
3. Create a stub file with inverse build tag: `// +build !myplugin`
4. Implement the `Plugin` interface
5. Register the plugin in `init()` function
6. Update the Makefile's `PLUGIN_TAGS` variable
7. Add blank import in `cmd/mvnenv/main.go`

## Why Plugins?

The plugin system allows:

1. **Smaller Binary Size**: Users who don't need plugin features get a smaller binary
2. **Enterprise Features**: Advanced features can be enabled only for enterprise deployments
3. **Maintainability**: Plugin code is isolated and can evolve independently
4. **Flexibility**: Organizations can choose which features to enable
5. **Testing**: Easier to test core functionality separately from plugins

## Plugin Status

Current plugins:

| Plugin | Tag | Status | Description |
|--------|-----|--------|-------------|
| Mirror | mirror | Stable | Create Nexus mirrors of Maven distributions |

## Troubleshooting

### Plugin Command Not Found

If a plugin command is not available, the binary was likely built without that plugin:

```bash
# Check available commands
mvnenv commands

# Rebuild with plugins
make build-plugins
```

### Build Errors

If you encounter build errors with plugins:

1. Ensure Go version is 1.21+
2. Verify build tags are correctly specified
3. Check that stub files exist for disabled plugins
4. Run `go mod tidy` to ensure dependencies are correct
