package repository

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/veenone/mvnenv-win/internal/download"
)

// ApacheArchive handles Maven downloads from Apache archive
type ApacheArchive struct {
	baseURL    string
	downloader *download.Downloader
}

// NewApacheArchive creates a new Apache archive client
func NewApacheArchive() *ApacheArchive {
	return &ApacheArchive{
		baseURL:    "https://archive.apache.org/dist/maven/maven-3/",
		downloader: download.NewDownloader(),
	}
}

// ListVersions scrapes Apache archive to find available Maven versions
func (a *ApacheArchive) ListVersions() ([]string, error) {
	resp, err := http.Get(a.baseURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse directory listing for version links
	// Pattern: <a href="3.9.4/">3.9.4/</a>
	re := regexp.MustCompile(`<a href="(\d+\.\d+\.\d+(?:-[^"/]+)?)/">`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var versions []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			version := match[1]
			if !seen[version] {
				versions = append(versions, version)
				seen[version] = true
			}
		}
	}

	return versions, nil
}

// DownloadVersion downloads a Maven version from Apache archive
func (a *ApacheArchive) DownloadVersion(version string, destPath string, progress download.ProgressCallback) error {
	// Construct URLs
	// https://archive.apache.org/dist/maven/maven-3/3.9.4/binaries/apache-maven-3.9.4-bin.zip
	url := fmt.Sprintf("%s%s/binaries/apache-maven-%s-bin.zip", a.baseURL, version, version)
	checksumURL := url + ".sha512"

	fmt.Printf("Downloading Maven %s from Apache archive...\n", version)

	// Download with checksum verification
	if err := a.downloader.DownloadWithChecksum(url, checksumURL, destPath, progress); err != nil {
		// If checksum download fails, try without verification
		if strings.Contains(err.Error(), "checksum") {
			fmt.Println("Warning: Checksum verification failed, downloading without verification")
			return a.downloader.Download(url, destPath, progress)
		}
		return err
	}

	return nil
}
