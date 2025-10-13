package nexus

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client represents a Nexus repository client
type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

// TLSConfig holds TLS configuration options
type TLSConfig struct {
	InsecureSkipVerify bool
	CAFile             string
}

// MavenMetadata represents the maven-metadata.xml structure
type MavenMetadata struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Versioning struct {
		Latest   string   `xml:"latest"`
		Release  string   `xml:"release"`
		Versions []string `xml:"versions>version"`
	} `xml:"versioning"`
}

// NewClient creates a new Nexus client
func NewClient(baseURL, username, password string, tlsConfig *TLSConfig) (*Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	// Configure TLS
	if tlsConfig != nil {
		if tlsConfig.InsecureSkipVerify {
			transport.TLSClientConfig.InsecureSkipVerify = true
		}

		if tlsConfig.CAFile != "" {
			caCert, err := os.ReadFile(tlsConfig.CAFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA certificate: %w", err)
			}

			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				return nil, fmt.Errorf("failed to parse CA certificate")
			}

			transport.TLSClientConfig.RootCAs = caCertPool
		}
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return &Client{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		httpClient: client,
	}, nil
}

// ListVersions retrieves available Maven versions from Nexus metadata
func (c *Client) ListVersions(ctx context.Context) ([]string, error) {
	// Construct maven-metadata.xml URL
	// Format: {baseURL}/org/apache/maven/apache-maven/maven-metadata.xml
	metadataURL := fmt.Sprintf("%s/org/apache/maven/apache-maven/maven-metadata.xml", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", metadataURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication if configured
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var metadata MavenMetadata
	if err := xml.Unmarshal(body, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return metadata.Versioning.Versions, nil
}

// DownloadVersion downloads a Maven distribution from Nexus
func (c *Client) DownloadVersion(ctx context.Context, version, destPath string, progress func(downloaded, total int64)) error {
	// Construct artifact URL
	// Format: {baseURL}/org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip
	artifactURL := fmt.Sprintf("%s/org/apache/maven/apache-maven/%s/apache-maven-%s-bin.zip",
		c.baseURL, version, version)

	req, err := http.NewRequestWithContext(ctx, "GET", artifactURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication if configured
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create destination file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download with progress tracking
	totalSize := resp.ContentLength
	downloaded := int64(0)

	buf := make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("failed to write file: %w", writeErr)
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
			return fmt.Errorf("failed to read response: %w", err)
		}
	}

	return nil
}
