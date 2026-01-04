package main

import (
	"github.com/blang/semver/v4"
	mvn "github.com/masahiro331/go-mvn-version"
)

// DependencyIssue is what we send to the Frontend to display
type DependencyIssue struct {
	ModName    string // The mod that is broken (e.g. "Sodium")
	MissingDep string // The mod it needs (e.g. "Fabric API")
	Reason     string // "Missing" or "Version Mismatch"
	Required   string // ">=0.40.0"
	Found      string // "0.30.0" or "Not Installed"
}

// ValidateDependencies checks all mods against the installed list
func ValidateDependencies(mods []ScanResult) []DependencyIssue {
	issues := []DependencyIssue{}

	// 1. Create a "Lookup Table" for fast checking
	// Map: ModID -> Version String
	installed := make(map[string]string)
	for _, mod := range mods {
		installed[mod.ModID] = mod.Version
	}

	// 2. Loop through every mod's dependencies
	for _, mod := range mods {
		for _, dep := range mod.Dependencies {
			if !dep.Mandatory {
				continue
			}

            if dep.ModID == "minecraft" || dep.ModID == "java" {
				continue
			}

			foundVersion, exists := installed[dep.ModID]
			if !exists {
				issues = append(issues, DependencyIssue{
					ModName:    mod.ModID,
					MissingDep: dep.ModID,
					Reason:     "Missing Dependency",
					Required:   dep.VersionRange,
					Found:      "Not Installed",
				})
				continue
			}

			if !isVersionCompatible(mod.Loader, foundVersion, dep.VersionRange) {
				issues = append(issues, DependencyIssue{
					ModName:    mod.ModID,
					MissingDep: dep.ModID,
					Reason:     "Version Mismatch",
					Required:   dep.VersionRange,
					Found:      foundVersion,
				})
			}
		}
	}

	return issues
}

// Helper to switch between Fabric and Forge logic
func isVersionCompatible(loader, currentVer, rangeStr string) bool {
	// If no version range is specified, or it's a wildcard, it's fine
	if rangeStr == "" || rangeStr == "*" {
		return true
	}

	if loader == "fabric" {
		// Fabric uses SemVer
		v, err := semver.Parse(currentVer)
		if err != nil {
			return true // If we can't parse the version, assume it's fine to avoid false alarms
		}
		expectedRange, err := semver.ParseRange(rangeStr)
		if err != nil {
			return true
		}
		return expectedRange(v)
	}

	// Forge/NeoForge uses Maven Ranges
	v, err := mvn.NewVersion(currentVer)
	if err != nil {
		return true
	}
	c, err := mvn.NewConstraints(rangeStr)
	if err != nil {
		return true
	}
	return c.Check(v)
}