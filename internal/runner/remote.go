package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	RepoOwner  = "Lokesh-Kudipudi"
	RepoName   = "DebugLab"
	CatalogURL = "https://raw.githubusercontent.com/Lokesh-Kudipudi/DebugLab/main/catalog.json"
	TreeAPIURL = "https://api.github.com/repos/Lokesh-Kudipudi/DebugLab/git/trees/main?recursive=1"
	RawBaseURL = "https://raw.githubusercontent.com/Lokesh-Kudipudi/DebugLab/main/"
)

// CatalogEntry represents a problem available remotely.
type CatalogEntry struct {
	Name       string `json:"name"`
	Difficulty string `json:"difficulty"`
	Language   string `json:"language"`
	Path       string `json:"path"`
}

type gitTreeResponse struct {
	Tree []struct {
		Path string `json:"path"`
		Type string `json:"type"`
	} `json:"tree"`
	Truncated bool `json:"truncated"`
}

// FetchCatalog downloads the catalog from GitHub, falling back to local catalog.json if available.
func FetchCatalog() ([]CatalogEntry, error) {
	// Try remote first
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(CatalogURL)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var catalog []CatalogEntry
		if err := json.NewDecoder(resp.Body).Decode(&catalog); err == nil {
			return catalog, nil
		}
	}

	// Fallback to local
	data, err := os.ReadFile("catalog.json")
	if err != nil {
		return nil, fmt.Errorf("could not fetch remote catalog and local catalog.json not found: %w", err)
	}

	var catalog []CatalogEntry
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse local catalog: %w", err)
	}

	return catalog, nil
}

// DownloadProblem fetches only the specific problem directory using the GitHub Trees API
// and raw.githubusercontent.com downloads.
func DownloadProblem(problemPath, destDir string) error {
	client := http.Client{Timeout: 30 * time.Second}

	// 1. Fetch the full repository tree
	req, err := http.NewRequest("GET", TreeAPIURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	// Use GitHub API version header (good practice)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch repository tree: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch tree, status: %d", resp.StatusCode)
	}

	var treeResponse gitTreeResponse
	if err := json.NewDecoder(resp.Body).Decode(&treeResponse); err != nil {
		return fmt.Errorf("failed to parse tree response: %w", err)
	}

	if treeResponse.Truncated {
		return fmt.Errorf("repository tree is too large (truncated). Need authentication or structural changes.")
	}

	// problemPath should exactly match the directory structure, e.g., "problems/python/python-sort-filter"
	prefix := problemPath + "/"
	foundFiles := false

	for _, item := range treeResponse.Tree {
		if item.Type != "blob" || !strings.HasPrefix(item.Path, prefix) {
			continue
		}
		foundFiles = true

		relPath := strings.TrimPrefix(item.Path, prefix)
		targetPath := filepath.Join(destDir, relPath)

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
		}

		rawURL := RawBaseURL + item.Path
		if err := downloadFile(client, rawURL, targetPath); err != nil {
			return fmt.Errorf("failed to download file %s: %w", item.Path, err)
		}
	}

	if !foundFiles {
		return fmt.Errorf("problem path %q not found in repository", problemPath)
	}

	return nil
}

func downloadFile(client http.Client, url, targetPath string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("got status %d", resp.StatusCode)
	}

	outFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
