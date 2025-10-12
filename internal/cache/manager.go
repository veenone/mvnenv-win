package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// VersionCache stores cached version information
type VersionCache struct {
	Versions  []string  `json:"versions"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Manager handles version cache operations
type Manager struct {
	cacheDir  string
	cacheFile string
}

// NewManager creates a new cache manager
func NewManager(mvnenvRoot string) *Manager {
	cacheDir := filepath.Join(mvnenvRoot, "cache")
	return &Manager{
		cacheDir:  cacheDir,
		cacheFile: filepath.Join(cacheDir, "versions.json"),
	}
}

// SaveVersions saves version list to cache
func (m *Manager) SaveVersions(versions []string) error {
	// Ensure cache directory exists
	if err := os.MkdirAll(m.cacheDir, 0755); err != nil {
		return fmt.Errorf("create cache directory: %w", err)
	}

	cache := VersionCache{
		Versions:  versions,
		UpdatedAt: time.Now(),
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}

	// Atomic write
	tempFile := m.cacheFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("write cache file: %w", err)
	}

	if err := os.Rename(tempFile, m.cacheFile); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("rename cache file: %w", err)
	}

	return nil
}

// LoadVersions loads version list from cache
func (m *Manager) LoadVersions() ([]string, error) {
	data, err := os.ReadFile(m.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cache exists yet
		}
		return nil, fmt.Errorf("read cache file: %w", err)
	}

	var cache VersionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("unmarshal cache: %w", err)
	}

	return cache.Versions, nil
}

// GetCacheAge returns how old the cache is
func (m *Manager) GetCacheAge() (time.Duration, error) {
	data, err := os.ReadFile(m.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("read cache file: %w", err)
	}

	var cache VersionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return 0, fmt.Errorf("unmarshal cache: %w", err)
	}

	return time.Since(cache.UpdatedAt), nil
}

// IsCacheStale checks if cache is older than specified duration
func (m *Manager) IsCacheStale(maxAge time.Duration) bool {
	age, err := m.GetCacheAge()
	if err != nil {
		return true // Consider stale if we can't read it
	}
	return age > maxAge
}

// CacheExists checks if cache file exists
func (m *Manager) CacheExists() bool {
	_, err := os.Stat(m.cacheFile)
	return err == nil
}
