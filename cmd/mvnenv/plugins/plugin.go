package plugins

import "github.com/spf13/cobra"

// Plugin represents a pluggable feature that can be enabled/disabled at build time
type Plugin interface {
	// Name returns the plugin name
	Name() string

	// Description returns a short description
	Description() string

	// Command returns the cobra command for this plugin
	Command() *cobra.Command
}

var registeredPlugins []Plugin

// Register adds a plugin to the registry
func Register(plugin Plugin) {
	registeredPlugins = append(registeredPlugins, plugin)
}

// GetAll returns all registered plugins
func GetAll() []Plugin {
	return registeredPlugins
}
