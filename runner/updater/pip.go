package updater

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/harryzcy/snuuze/types"
)

func upgradePip(cache *Cache, info *types.UpgradeInfo) error {
	file, err := cache.Get(info.Dependency.File)
	if err != nil {
		return fmt.Errorf("failed to get file %s from cache: %s", info.Dependency.File, err)
	}

	lines := bytes.Split(file, []byte("\n"))
	lineIdx := info.Dependency.Position.Line

	oldVersion := extractPipVersionConstrains(info.Dependency.Version)
	newVersion := info.ToVersion
	lines[lineIdx] = bytes.Replace(lines[lineIdx], []byte(oldVersion), []byte(newVersion), 1)

	file = bytes.Join(lines, []byte("\n"))
	cache.Set(info.Dependency.File, file)

	return nil
}

func extractPipVersionConstrains(version string) string {
	version = strings.TrimPrefix(version, "==")
	version = strings.TrimPrefix(version, ">=")
	version = strings.TrimPrefix(version, ">")
	version = strings.TrimSpace(version)
	return version
}
