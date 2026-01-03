/*
 * Glacier - Mod Manager & Dependency Check
 * Copyright (C) 2026 BySnowden
 * License: GPLv3
 */

package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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