/*
 * Glacier - Mod Manager & Dependency Check
 * Copyright (C) 2026 BySnowden
 * License: GPLv3
 */

package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

type UpdateResponse struct {
	UpdateAvailable bool   `json:"updateAvailable"`
	LatestVersion   string `json:"latestVersion"`
	DownloadUrl     string `json:"downloadUrl"`
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
		runtime.EventsEmit(pw.ctx, "update-progress", int(percentage))
	}

	return n, nil
}

const AppVersion = "1.1.1"

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// CheckForUpdates hits the GitHub API to find the specific .exe asset
func (a *App) CheckForUpdates() UpdateResponse {
	// 1. Request latest release metadata from GitHub API
	resp, err := http.Get("https://api.github.com/repos/BySnowden/Glacier/releases/latest")
	if err != nil {
		return UpdateResponse{false, "", ""}
	}
	defer resp.Body.Close()

	// 2. Parse JSON response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return UpdateResponse{false, "", ""}
	}

	tagName, ok := result["tag_name"].(string)
	if !ok {
		return UpdateResponse{false, "", ""}
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
		}

		return UpdateResponse{
			UpdateAvailable: true,
			LatestVersion:   cleanTag,
			DownloadUrl:     downloadUrl,
		}
	}

	return UpdateResponse{false, AppVersion, ""}
}

// DownloadAndInstall downloads the update and runs the installer
func (a *App) DownloadAndInstall(url string) error {
	// 1. Create a temp file
	tempFile, err := os.CreateTemp("", "Glacier-Update-*.exe")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// 2. Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 3. Set up progress tracking
	pw := &progressWriter{
		total: uint64(resp.ContentLength),
		ctx:   a.ctx,
	}

	// 4. Copy data from internet to file (with progress)
	if _, err = io.Copy(io.MultiWriter(tempFile, pw), resp.Body); err != nil {
		return err
	}

	// 5. Run the installer silently
	// "cmd /C start ..." allows it to run detached so we can close this app safely
	cmd := exec.Command("cmd", "/C", "start", tempFile.Name(), "/SILENT")
	if err := cmd.Start(); err != nil {
		return err
	}

	// 6. Kill this app so the installer can overwrite files
	runtime.Quit(a.ctx)
	return nil
}

// TriggerUpdate opens the download URL in the user's default browser (Manual Fallback)
func (a *App) TriggerUpdate(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
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
