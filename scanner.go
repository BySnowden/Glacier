package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"
)

type ScanDependency struct {
	ModID        string
	VersionRange string
	Mandatory    bool
}

type ScanResult struct {
	ModID        string
	Version      string
	Loader       string // "fabric", "forge", "neoforge"
	Dependencies []ScanDependency
}


func ScanJarForIssues(path string) (*ScanResult, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	// Check for Fabric first
	if f := locateZipFile(r, "fabric.mod.json"); f != nil {
		return readFabricMetadata(f)
	}

	// Check for NeoForge
	if f := locateZipFile(r, "META-INF/neoforge.mods.toml"); f != nil {
		return readForgeMetadata(f, "neoforge")
	}

	// Also check for classic Forge
	if f := locateZipFile(r, "META-INF/mods.toml"); f != nil {
		return readForgeMetadata(f, "forge")
	}

	return nil, fmt.Errorf("no mod metadata found")
}

func locateZipFile(r *zip.ReadCloser, name string) *zip.File {
	for _, f := range r.File {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func readFabricMetadata(f *zip.File) (*ScanResult, error) {
	rc, _ := f.Open()
	defer rc.Close()

	var data struct {
		ID      string            `json:"id"`
		Version string            `json:"version"`
		Depends map[string]string `json:"depends"`
	}

	bytes, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	mod := &ScanResult{
		ModID:   data.ID,
		Version: data.Version,
		Loader:  "fabric",
	}

	for id, ver := range data.Depends {
		mod.Dependencies = append(mod.Dependencies, ScanDependency{
			ModID:        id,
			VersionRange: ver,
			Mandatory:    true,
		})
	}
	return mod, nil
}

func readForgeMetadata(f *zip.File, loaderType string) (*ScanResult, error) {
	rc, _ := f.Open()
	defer rc.Close()

	var data struct {
		Mods []struct {
			ModId   string `toml:"modId"`
			Version string `toml:"version"`
		} `toml:"mods"`
		Dependencies map[string][]struct {
			ModId        string `toml:"modId"`
			VersionRange string `toml:"versionRange"`
			Mandatory    bool   `toml:"mandatory"`
		} `toml:"dependencies"`
	}

	bytes, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	if len(data.Mods) == 0 {
		return nil, fmt.Errorf("empty mods list")
	}

	mainMod := data.Mods[0]

	mod := &ScanResult{
		ModID:   mainMod.ModId,
		Version: mainMod.Version,
		Loader:  loaderType,
	}

	if deps, ok := data.Dependencies[mainMod.ModId]; ok {
		for _, d := range deps {
			mod.Dependencies = append(mod.Dependencies, ScanDependency{
				ModID:        d.ModId,
				VersionRange: d.VersionRange,
				Mandatory:    d.Mandatory,
			})
		}
	}

	return mod, nil
}