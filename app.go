/*
 * Glacier - Mod Manager & Dependency Check
 * Copyright (C) 2026 BySnowden
 * License: GPLv3
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

type UpdateResponse struct {
	UpdateAvailable bool   `json:"updateAvailable"`
	LatestVersion   string `json:"latestVersion"`
	DownloadUrl     string `json:"downloadUrl"`
	Error           string `json:"error,omitempty"`
}

type UpdateProgress struct {
	Stage      string `json:"stage"`
	Percentage int    `json:"percentage"`
	Message    string `json:"message"`
}

// progressWriter tracks download progress
type progressWriter struct {
	total      uint64
	downloaded uint64
	ctx        context.Context
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.downloaded += uint64(n)

	if pw.total > 0 {
		percentage := float64(pw.downloaded) / float64(pw.total) * 100
		runtime.EventsEmit(pw.ctx, "update-progress", UpdateProgress{
			Stage:      "downloading",
			Percentage: int(percentage),
			Message:    fmt.Sprintf("Downloaded %d%% of update", int(percentage)),
		})
	}

	return n, nil
}

const AppVersion = "1.1.4"

// Update configuration
const (
	GITHUB_OWNER         = "BySnowden"
	GITHUB_REPO          = "Glacier"
	UPDATE_CHECK_TIMEOUT = 10 * time.Second
	DOWNLOAD_TIMEOUT     = 30 * time.Minute
)

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// CheckForUpdates hits the GitHub API to find the specific .exe asset
func (a *App) CheckForUpdates() UpdateResponse {
	// 1. Request latest release metadata from GitHub API with timeout
	client := &http.Client{Timeout: UPDATE_CHECK_TIMEOUT}
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GITHUB_OWNER, GITHUB_REPO)
	resp, err := client.Get(apiURL)
	if err != nil {
		return UpdateResponse{false, "", "", fmt.Sprintf("Failed to check for updates: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return UpdateResponse{false, "", "", fmt.Sprintf("GitHub API returned status: %d", resp.StatusCode)}
	}

	// 2. Parse JSON response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return UpdateResponse{false, "", "", fmt.Sprintf("Failed to parse update response: %v", err)}
	}

	tagName, ok := result["tag_name"].(string)
	if !ok {
		return UpdateResponse{false, "", "", "Invalid release data from GitHub"}
	}
	cleanTag := strings.TrimPrefix(tagName, "v")

	if cleanTag != AppVersion {
		assets, _ := result["assets"].([]interface{})
		downloadUrl := ""

		for _, asset := range assets {
			a := asset.(map[string]interface{})
			name := a["name"].(string)

			if strings.HasSuffix(strings.ToLower(name), ".exe") {
				downloadUrl = a["browser_download_url"].(string)
				// Prefer the amd64 one if running on Intel/AMD, or arm64 if on Surface
				// For now, this just grabs the first executable it sees.
				break
			}
		}

		// Fallback: If no exe found in assets, link to the release page
		if downloadUrl == "" {
			downloadUrl, _ = result["html_url"].(string)
			return UpdateResponse{false, "", "", "No executable found in latest release"}
		}

		return UpdateResponse{
			UpdateAvailable: true,
			LatestVersion:   cleanTag,
			DownloadUrl:     downloadUrl,
		}
	}

	return UpdateResponse{false, AppVersion, "", ""}
}

// DownloadAndInstall downloads the update and runs the installer
func (a *App) DownloadAndInstall(url string) error {
	runtime.EventsEmit(a.ctx, "update-progress", UpdateProgress{
		Stage:      "preparing",
		Percentage: 0,
		Message:    "Preparing to download update...",
	})

	// 1. Create a temp file
	tempFile, err := os.CreateTemp("", "Glacier-Update-*.exe")
	if err != nil {
		runtime.EventsEmit(a.ctx, "update-error", fmt.Sprintf("Failed to create temp file: %v", err))
		return err
	}
	tempPath := tempFile.Name()
	tempFile.Close() // Close immediately, we'll reopen for writing

	// 2. Download with timeout
	runtime.EventsEmit(a.ctx, "update-progress", UpdateProgress{
		Stage:      "downloading",
		Percentage: 0,
		Message:    "Starting download...",
	})

	client := &http.Client{Timeout: DOWNLOAD_TIMEOUT}
	resp, err := client.Get(url)
	if err != nil {
		os.Remove(tempPath)
		runtime.EventsEmit(a.ctx, "update-error", fmt.Sprintf("Failed to start download: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		os.Remove(tempPath)
		runtime.EventsEmit(a.ctx, "update-error", fmt.Sprintf("Download failed with status: %d", resp.StatusCode))
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// 3. Reopen temp file for writing
	tempFile, err = os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		os.Remove(tempPath)
		runtime.EventsEmit(a.ctx, "update-error", fmt.Sprintf("Failed to open temp file: %v", err))
		return err
	}

	// 4. Copy with progress tracking
	pw := &progressWriter{
		total: uint64(resp.ContentLength),
		ctx:   a.ctx,
	}

	if _, err = io.Copy(io.MultiWriter(tempFile, pw), resp.Body); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		runtime.EventsEmit(a.ctx, "update-error", fmt.Sprintf("Download failed: %v", err))
		return err
	}

	// 5. IMPORTANT: Close the file now so Windows releases the lock
	tempFile.Close()

	runtime.EventsEmit(a.ctx, "update-progress", UpdateProgress{
		Stage:      "installing",
		Percentage: 100,
		Message:    "Starting installation...",
	})

	// 6. Run the installer with proper flags
	cmd := exec.Command(tempPath, "/SILENT", "/SP-", "/NORESTART")

	if err := cmd.Start(); err != nil {
		os.Remove(tempPath)
		runtime.EventsEmit(a.ctx, "update-error", fmt.Sprintf("Failed to start installer: %v", err))
		return err
	}

	// 7. Give installer a moment to start, then quit
	go func() {
		time.Sleep(2 * time.Second)
		runtime.Quit(a.ctx)
	}()

	return nil
}

// TriggerUpdate opens the download URL in the user's default browser (Manual Fallback)
func (a *App) TriggerUpdate(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

// GetAppVersion returns the current application version
func (a *App) GetAppVersion() string {
	return AppVersion
}

// RestartApp restarts the application (useful after updates)
func (a *App) RestartApp() error {
	// Get the current executable path
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// Start new instance
	cmd := exec.Command(executable)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Quit current instance
	runtime.Quit(a.ctx)
	return nil
}

func (a *App) SelectModFolder() string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Mods Folder",
	})
	if err != nil {
		return ""
	}
	return selection
}

func (a *App) ScanMods(folderPath string) ([]ModMetadata, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var jars []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".jar") {
			jars = append(jars, file.Name())
		}
	}

	results := make([]ModMetadata, 0, len(jars))
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, jarName := range jars {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			fullPath := filepath.Join(folderPath, name)
			meta := ParseJar(fullPath)

			mutex.Lock()
			results = append(results, meta)
			mutex.Unlock()
		}(jarName)
	}

	wg.Wait()
	return results, nil
}

func (a *App) ValidateMods(folderPath string) ([]DependencyIssue, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var scannedResults []ScanResult

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jar") {
			fullPath := filepath.Join(folderPath, entry.Name())

			res, err := ScanJarForIssues(fullPath)
			if err == nil {
				scannedResults = append(scannedResults, *res)
			}
		}
	}
	return ValidateDependencies(scannedResults), nil
}
