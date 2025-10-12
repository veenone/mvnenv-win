package download

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Downloader handles file downloads with progress tracking
type Downloader struct {
	client *http.Client
}

// NewDownloader creates a new downloader
func NewDownloader() *Downloader {
	return &Downloader{
		client: &http.Client{},
	}
}

// ProgressCallback is called during download to report progress
type ProgressCallback func(downloaded int64, total int64)

// Download downloads a file from URL to destination path
func (d *Downloader) Download(url string, destPath string, progress ProgressCallback) error {
	// Create HTTP request
	resp, err := d.client.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Create destination file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	// Get total size
	totalSize := resp.ContentLength

	// Copy with progress
	var downloaded int64
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("write file: %w", writeErr)
			}
			downloaded += int64(n)
			if progress != nil {
				progress(downloaded, totalSize)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read response: %w", err)
		}
	}

	return nil
}

// DownloadWithChecksum downloads a file and verifies SHA-512 checksum
func (d *Downloader) DownloadWithChecksum(url string, checksumURL string, destPath string, progress ProgressCallback) error {
	// Download main file
	if err := d.Download(url, destPath, progress); err != nil {
		return err
	}

	// Download checksum
	resp, err := d.client.Get(checksumURL)
	if err != nil {
		return fmt.Errorf("download checksum: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("checksum HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	checksumData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read checksum: %w", err)
	}

	// Parse checksum (format: "checksum filename" or just "checksum")
	checksumStr := strings.TrimSpace(string(checksumData))
	parts := strings.Fields(checksumStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty checksum file")
	}
	expectedChecksum := strings.ToLower(parts[0])

	// Calculate actual checksum
	actualChecksum, err := calculateSHA512(destPath)
	if err != nil {
		return fmt.Errorf("calculate checksum: %w", err)
	}

	// Verify checksum
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// calculateSHA512 calculates SHA-512 checksum of a file
func calculateSHA512(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
