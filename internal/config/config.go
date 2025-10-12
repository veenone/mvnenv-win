package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration file structure
type Config struct {
	Version       string `yaml:"version"`
	GlobalVersion string `yaml:"global_version,omitempty"`
	AutoRehash    bool   `yaml:"auto_rehash"`
	mu            sync.RWMutex
}

// Manager handles configuration file operations
type Manager struct {
	configPath string
	config     *Config
	mu         sync.RWMutex
}

// NewManager creates a new configuration manager
func NewManager(mvnenvRoot string) *Manager {
	configDir := filepath.Join(mvnenvRoot, "config")
	configPath := filepath.Join(configDir, "config.yaml")

	return &Manager{
		configPath: configPath,
	}
}

// Load loads configuration from disk
func (m *Manager) Load() (*Config, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if config file exists
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		// Return default config
		m.config = &Config{
			Version:    "1.0",
			AutoRehash: true,
		}
		return m.config, nil
	}

	// Read config file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	m.config = &config
	return m.config, nil
}

// Save saves configuration to disk
func (m *Manager) Save(config *Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Ensure config directory exists
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	// Write atomically (temp file + rename)
	tmpPath := m.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write temp config: %w", err)
	}

	if err := os.Rename(tmpPath, m.configPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename config file: %w", err)
	}

	m.config = config
	return nil
}

// GetGlobalVersion returns the global Maven version
func (m *Manager) GetGlobalVersion() (string, error) {
	config, err := m.Load()
	if err != nil {
		return "", err
	}

	return config.GlobalVersion, nil
}

// SetGlobalVersion sets the global Maven version
func (m *Manager) SetGlobalVersion(version string) error {
	config, err := m.Load()
	if err != nil {
		config = &Config{
			Version:    "1.0",
			AutoRehash: true,
		}
	}

	config.GlobalVersion = version
	return m.Save(config)
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() (*Config, error) {
	return m.Load()
}
