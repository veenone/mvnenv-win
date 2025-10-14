package repository

import (
	"context"
	"fmt"

	"github.com/veenone/mvnenv-win/internal/config"
	"github.com/veenone/mvnenv-win/internal/download"
	"github.com/veenone/mvnenv-win/internal/nexus"
)

// Manager manages multiple repository sources
type Manager struct {
	apache      *ApacheArchive
	nexusClient *nexus.Client
	config      *config.Manager
	mvnenvRoot  string
}

// NewManager creates a new repository manager
func NewManager(mvnenvRoot string) *Manager {
	return &Manager{
		apache:     NewApacheArchive(),
		config:     config.NewManager(mvnenvRoot),
		mvnenvRoot: mvnenvRoot,
	}
}

// initializeNexus initializes Nexus client if configured
func (m *Manager) initializeNexus() error {
	if m.nexusClient != nil {
		return nil
	}

	cfg, err := m.config.Load()
	if err != nil {
		return err
	}

	// Check new structure first (repositories.nexus), fall back to old structure (nexus)
	var nexusCfg *config.NexusConfig
	if cfg.Repositories != nil && cfg.Repositories.Nexus != nil {
		nexusCfg = cfg.Repositories.Nexus
	}

	if nexusCfg == nil || !nexusCfg.Enabled {
		return nil
	}

	// Convert config TLS to nexus TLS
	var tlsConfig *nexus.TLSConfig
	if nexusCfg.TLS != nil {
		tlsConfig = &nexus.TLSConfig{
			InsecureSkipVerify: nexusCfg.TLS.InsecureSkipVerify,
			CAFile:             nexusCfg.TLS.CAFile,
		}
	}

	client, err := nexus.NewClient(
		nexusCfg.BaseURL,
		nexusCfg.Username,
		nexusCfg.Password,
		tlsConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to create Nexus client: %w", err)
	}

	m.nexusClient = client
	return nil
}

// ListVersions returns available versions from all configured sources
func (m *Manager) ListVersions() ([]string, error) {
	var allVersions []string
	seen := make(map[string]bool)

	// Try Nexus first if configured
	if err := m.initializeNexus(); err == nil && m.nexusClient != nil {
		ctx := context.Background()
		nexusVersions, err := m.nexusClient.ListVersions(ctx)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch versions from Nexus: %v\n", err)
		} else {
			for _, v := range nexusVersions {
				if !seen[v] {
					allVersions = append(allVersions, v)
					seen[v] = true
				}
			}
		}
	}

	// Get versions from Apache archive
	apacheVersions, err := m.apache.ListVersions()
	if err != nil {
		if len(allVersions) == 0 {
			return nil, fmt.Errorf("failed to fetch versions from Apache: %w", err)
		}
		fmt.Printf("Warning: Failed to fetch versions from Apache archive: %v\n", err)
	} else {
		for _, v := range apacheVersions {
			if !seen[v] {
				allVersions = append(allVersions, v)
				seen[v] = true
			}
		}
	}

	return allVersions, nil
}

// DownloadVersion downloads a version from the first available source
func (m *Manager) DownloadVersion(version string, destPath string, progress download.ProgressCallback) error {
	// Try Nexus first if configured
	if err := m.initializeNexus(); err == nil && m.nexusClient != nil {
		fmt.Printf("Attempting to download Maven %s from Nexus...\n", version)
		ctx := context.Background()

		nexusProgress := func(downloaded, total int64) {
			if progress != nil {
				progress(downloaded, total)
			}
		}

		err := m.nexusClient.DownloadVersion(ctx, version, destPath, nexusProgress)
		if err == nil {
			return nil
		}
		fmt.Printf("Nexus download failed: %v\n", err)
		fmt.Println("Falling back to Apache archive...")
	}

	// Fall back to Apache archive
	return m.apache.DownloadVersion(version, destPath, progress)
}
