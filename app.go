/*
 * Glacier - Mod Manager & Dependency Check
 * Copyright (C) 2026 BySnowden
 * License: GPLv3
 */

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
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

const AppVersion = "1.1.1"

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// CheckForUpdates hits the GitHub API to see if a newer release exists
func (a *App) CheckForUpdates() UpdateResponse {
	// 1. Request latest release metadata from GitHub API (Not the website HTML)
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

	// 3. Extract the tag name (e.g., "v1.2.0")
	tagName, ok := result["tag_name"].(string)
	if !ok {
		return UpdateResponse{false, "", ""}
	}

	// 4. Compare versions (Remove 'v' prefix if present)
	cleanTag := strings.TrimPrefix(tagName, "v")

	if cleanTag != AppVersion {
		// Get the HTML URL (the page where users can download the exe)
		downloadUrl, _ := result["html_url"].(string)

		return UpdateResponse{
			UpdateAvailable: true,
			LatestVersion:   cleanTag,
			DownloadUrl:     downloadUrl,
		}
	}

	// No update available
	return UpdateResponse{false, AppVersion, ""}
}

// TriggerUpdate opens the download URL in the user's default browser
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
	// gets all jar files
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var scannedResults []ScanResult

	// scans every jar for metadata
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jar") {
			fullPath := filepath.Join(folderPath, entry.Name())

			// calls the function from scanner.go
			res, err := ScanJarForIssues(fullPath)
			if err == nil {
				scannedResults = append(scannedResults, *res)
			}
		}
	}
	return ValidateDependencies(scannedResults), nil
}
