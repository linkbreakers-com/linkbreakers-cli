package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	latestURL     = "https://cli.linkbreakers.com/latest.json"
	installURL    = "https://cli.linkbreakers.com/install.sh"
	cacheFilename = "update-check.json"
	cacheTTL      = 24 * time.Hour
)

type ReleaseInfo struct {
	Version    string `json:"version"`
	ReleaseURL string `json:"release_url"`
}

type checkCache struct {
	LastCheckedAt string      `json:"last_checked_at"`
	Release       ReleaseInfo `json:"release"`
}

func LatestRelease(client *http.Client) (ReleaseInfo, error) {
	req, err := http.NewRequest(http.MethodGet, latestURL, nil)
	if err != nil {
		return ReleaseInfo{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "linkbreakers-cli")

	resp, err := client.Do(req)
	if err != nil {
		return ReleaseInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return ReleaseInfo{}, fmt.Errorf("version endpoint returned %s", resp.Status)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return ReleaseInfo{}, err
	}
	return release, nil
}

func MaybeNotify(currentVersion string, stderr io.Writer) {
	if currentVersion == "" || currentVersion == "dev" {
		return
	}

	release, err := cachedRelease()
	if err != nil || release.Version == "" {
		return
	}

	if compareVersions(release.Version, currentVersion) <= 0 {
		return
	}

	_, _ = fmt.Fprintf(stderr,
		"A new version of linkbreakers is available: %s (current: %s)\nUpdate with: linkbreakers self-update\nInstaller: curl -fsSL %s | bash\n",
		release.Version,
		currentVersion,
		installURL,
	)
}

func SelfUpdate(currentVersion string) (ReleaseInfo, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	release, err := LatestRelease(client)
	if err != nil {
		return ReleaseInfo{}, err
	}

	if currentVersion != "" && currentVersion != "dev" && compareVersions(release.Version, currentVersion) <= 0 {
		return release, nil
	}

	if runtime.GOOS == "windows" {
		return release, fmt.Errorf("self-update is not yet supported on Windows; download the latest release from %s", release.ReleaseURL)
	}

	execPath, err := os.Executable()
	if err != nil {
		return ReleaseInfo{}, err
	}

	assetURL, ext, err := assetURLFor(release.Version)
	if err != nil {
		return ReleaseInfo{}, err
	}

	req, err := http.NewRequest(http.MethodGet, assetURL, nil)
	if err != nil {
		return ReleaseInfo{}, err
	}
	req.Header.Set("User-Agent", "linkbreakers-cli")

	resp, err := client.Do(req)
	if err != nil {
		return ReleaseInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return ReleaseInfo{}, fmt.Errorf("download failed: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReleaseInfo{}, err
	}

	bin, err := extractBinary(data, ext)
	if err != nil {
		return ReleaseInfo{}, err
	}

	targetDir := filepath.Dir(execPath)
	tmpPath := filepath.Join(targetDir, ".linkbreakers.new")
	if err := os.WriteFile(tmpPath, bin, 0o755); err != nil {
		return ReleaseInfo{}, err
	}
	if err := os.Rename(tmpPath, execPath); err != nil {
		_ = os.Remove(tmpPath)
		return ReleaseInfo{}, err
	}

	_ = saveCache(checkCache{
		LastCheckedAt: time.Now().UTC().Format(time.RFC3339),
		Release:       release,
	})

	return release, nil
}

func cachedRelease() (ReleaseInfo, error) {
	cache, err := readCache()
	if err == nil && cache.Release.Version != "" {
		checkedAt, parseErr := time.Parse(time.RFC3339, cache.LastCheckedAt)
		if parseErr == nil && time.Since(checkedAt) < cacheTTL {
			return cache.Release, nil
		}
	}

	client := &http.Client{Timeout: 3 * time.Second}
	release, err := LatestRelease(client)
	if err != nil {
		if cache.Release.Version != "" {
			return cache.Release, nil
		}
		return ReleaseInfo{}, err
	}

	_ = saveCache(checkCache{
		LastCheckedAt: time.Now().UTC().Format(time.RFC3339),
		Release:       release,
	})

	return release, nil
}

func cachePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "linkbreakers", cacheFilename), nil
}

func readCache() (checkCache, error) {
	path, err := cachePath()
	if err != nil {
		return checkCache{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return checkCache{}, err
	}

	var cache checkCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return checkCache{}, err
	}
	return cache, nil
}

func saveCache(cache checkCache) error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func assetURLFor(version string) (string, string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goarch {
	case "amd64", "arm64":
	default:
		return "", "", fmt.Errorf("unsupported architecture %s", goarch)
	}

	switch goos {
	case "darwin", "linux":
		return fmt.Sprintf(
			"https://github.com/linkbreakers-com/linkbreakers-cli/releases/download/v%s/linkbreakers-cli_%s_%s_%s.tar.gz",
			version, version, goos, goarch,
		), ".tar.gz", nil
	case "windows":
		return fmt.Sprintf(
			"https://github.com/linkbreakers-com/linkbreakers-cli/releases/download/v%s/linkbreakers-cli_%s_windows_%s.zip",
			version, version, goarch,
		), ".zip", nil
	default:
		return "", "", fmt.Errorf("unsupported OS %s", goos)
	}
}

func extractBinary(data []byte, ext string) ([]byte, error) {
	switch ext {
	case ".tar.gz":
		gzr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer gzr.Close()

		tr := tar.NewReader(gzr)
		for {
			hdr, err := tr.Next()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return nil, err
			}
			if filepath.Base(hdr.Name) == "linkbreakers" {
				return io.ReadAll(tr)
			}
		}
	case ".zip":
		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}
		for _, f := range zr.File {
			if filepath.Base(f.Name) == "linkbreakers.exe" {
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}
				defer rc.Close()
				return io.ReadAll(rc)
			}
		}
	}

	return nil, errors.New("linkbreakers binary not found in archive")
}

func compareVersions(a, b string) int {
	split := func(v string) []string {
		return strings.Split(strings.TrimPrefix(v, "v"), ".")
	}

	aa := split(a)
	bb := split(b)
	n := len(aa)
	if len(bb) > n {
		n = len(bb)
	}

	for i := 0; i < n; i++ {
		var av, bv int
		if i < len(aa) {
			fmt.Sscanf(aa[i], "%d", &av)
		}
		if i < len(bb) {
			fmt.Sscanf(bb[i], "%d", &bv)
		}
		if av > bv {
			return 1
		}
		if av < bv {
			return -1
		}
	}
	return 0
}
