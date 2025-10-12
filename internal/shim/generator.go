package shim

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ShimGenerator creates and manages Maven command shims
type ShimGenerator struct {
	shimsDir    string
	shimBinary  string
	versionsDir string
}

// NewShimGenerator creates a shim generator
func NewShimGenerator(mvnenvRoot string) *ShimGenerator {
	return &ShimGenerator{
		shimsDir:    filepath.Join(mvnenvRoot, "shims"),
		shimBinary:  filepath.Join(mvnenvRoot, "bin", "shim.exe"),
		versionsDir: filepath.Join(mvnenvRoot, "versions"),
	}
}

// GenerateShims creates shim executables for all Maven commands
func (g *ShimGenerator) GenerateShims() ([]string, error) {
	// Ensure shims directory exists
	if err := os.MkdirAll(g.shimsDir, 0755); err != nil {
		return nil, fmt.Errorf("create shims directory: %w", err)
	}

	// Core Maven commands to shim
	commands := []string{"mvn", "mvnDebug"}

	// Scan installed versions for additional commands
	additionalCmds, err := g.discoverAdditionalCommands()
	if err == nil {
		commands = append(commands, additionalCmds...)
	}

	var generatedPaths []string

	for _, cmd := range commands {
		// Generate .exe shim
		exePath, err := g.generateShimFile(cmd, ".exe")
		if err != nil {
			return nil, fmt.Errorf("generate %s.exe: %w", cmd, err)
		}
		generatedPaths = append(generatedPaths, exePath)

		// Generate .cmd shim
		cmdPath, err := g.generateBatchShim(cmd)
		if err != nil {
			return nil, fmt.Errorf("generate %s.cmd: %w", cmd, err)
		}
		generatedPaths = append(generatedPaths, cmdPath)
	}

	return generatedPaths, nil
}

// generateShimFile creates executable shim by copying shim.exe
func (g *ShimGenerator) generateShimFile(command string, ext string) (string, error) {
	destPath := filepath.Join(g.shimsDir, command+ext)

	// Check if shim.exe exists
	if _, err := os.Stat(g.shimBinary); err != nil {
		return "", fmt.Errorf("shim.exe not found at %s: %w", g.shimBinary, err)
	}

	// Copy shim.exe to destination
	if err := copyFile(g.shimBinary, destPath); err != nil {
		return "", fmt.Errorf("copy shim binary: %w", err)
	}

	return destPath, nil
}

// generateBatchShim creates .cmd shim that calls .exe shim
func (g *ShimGenerator) generateBatchShim(command string) (string, error) {
	destPath := filepath.Join(g.shimsDir, command+".cmd")

	script := fmt.Sprintf(`@echo off
"%%~dp0%s.exe" %%*
exit /b %%ERRORLEVEL%%
`, command)

	if err := os.WriteFile(destPath, []byte(script), 0755); err != nil {
		return "", fmt.Errorf("write batch shim: %w", err)
	}

	return destPath, nil
}

// discoverAdditionalCommands scans installed versions for commands like mvnyjp
func (g *ShimGenerator) discoverAdditionalCommands() ([]string, error) {
	entries, err := os.ReadDir(g.versionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	cmdSet := make(map[string]bool)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		binDir := filepath.Join(g.versionsDir, entry.Name(), "bin")
		binEntries, err := os.ReadDir(binDir)
		if err != nil {
			continue
		}

		for _, binEntry := range binEntries {
			name := binEntry.Name()
			if strings.HasSuffix(name, ".cmd") {
				cmd := strings.TrimSuffix(name, ".cmd")
				if cmd != "mvn" && cmd != "mvnDebug" && !cmdSet[cmd] {
					cmdSet[cmd] = true
				}
			}
		}
	}

	var additionalCmds []string
	for cmd := range cmdSet {
		additionalCmds = append(additionalCmds, cmd)
	}

	return additionalCmds, nil
}

// copyFile copies file from src to dst atomically
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	tmpDst := dst + ".tmp"
	destination, err := os.Create(tmpDst)
	if err != nil {
		return err
	}

	_, err = io.Copy(destination, source)
	destination.Close()
	if err != nil {
		os.Remove(tmpDst)
		return err
	}

	if err := os.Rename(tmpDst, dst); err != nil {
		os.Remove(tmpDst)
		return err
	}

	return nil
}
