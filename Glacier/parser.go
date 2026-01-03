package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"path/filepath" // Use standard library!

	"github.com/BurntSushi/toml"
)

// This is the ONE place this struct should be defined.
type ModMetadata struct {
	FileName    string `json:"fileName"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Loader      string `json:"loader"`
	Description string `json:"description"`
}

// Internal helper structs for JSON/TOML parsing
type FabricStruct struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type NeoForgeStruct struct {
	Mods []struct {
		ModId       string `toml:"modId"`
		DisplayName string `toml:"displayName"`
		Version     string `toml:"version"`
		Description string `toml:"description"`
	} `toml:"mods"`
}

// ParseJar is now a standalone function, not attached to 'App'
func ParseJar(path string) ModMetadata {
	result := ModMetadata{
		FileName: filepath.Base(path),
		Name:     filepath.Base(path),
		Loader:   "Unknown",
		Version:  "?.?.?",
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		return result
	}
	defer r.Close()

	// 1. Look for Fabric
	for _, f := range r.File {
		if f.Name == "fabric.mod.json" {
			parseFabric(f, &result)
			return result
		}
	}

	// 2. Look for NeoForge / Forge
	for _, f := range r.File {
		if f.Name == "META-INF/neoforge.mods.toml" || f.Name == "META-INF/mods.toml" {
			parseNeoForge(f, &result)
			return result
		}
	}
	return result
}

func parseFabric(f *zip.File, result *ModMetadata) {
	rc, _ := f.Open()
	defer rc.Close()

	var data FabricStruct
	if err := json.NewDecoder(rc).Decode(&data); err == nil {
		result.Name = data.Name
		result.Version = data.Version
		result.Loader = "Fabric"
		result.Description = data.Description
	}
}

func parseNeoForge(f *zip.File, result *ModMetadata) {
	rc, _ := f.Open()
	defer rc.Close()
	bytes, _ := io.ReadAll(rc)

	var data NeoForgeStruct
	if _, err := toml.Decode(string(bytes), &data); err == nil && len(data.Mods) > 0 {
		mod := data.Mods[0]
		result.Name = mod.DisplayName
		result.Version = mod.Version
		result.Loader = "NeoForge"
		result.Description = mod.Description
	}
}