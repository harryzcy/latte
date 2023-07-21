package gomajor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// ModPrefix returns the module path with no SIV
func ModPrefix(modpath string) string {
	prefix, _, ok := module.SplitPathVersion(modpath)
	if !ok {
		prefix = modpath
	}
	return prefix
}

// ModMajor returns the major version in vN format
func ModMajor(modpath string) (string, bool) {
	_, major, ok := module.SplitPathVersion(modpath)
	if ok {
		major = strings.TrimPrefix(major, "/")
		major = strings.TrimPrefix(major, ".")
	}
	return major, ok
}

// SplitPath splits the package path into the module path and the package subdirectory.
// It requires the a module path prefix to figure this out.
func SplitPath(modprefix, pkgpath string) (modpath, pkgdir string, ok bool) {
	if !strings.HasPrefix(pkgpath, modprefix) {
		return "", "", false
	}
	modpathlen := len(modprefix)
	if rest := pkgpath[modpathlen:]; len(rest) > 0 && rest[0] != '/' && rest[0] != '.' {
		return "", "", false
	}
	if strings.HasPrefix(pkgpath[modpathlen:], "/") {
		modpathlen++
	}
	if idx := strings.Index(pkgpath[modpathlen:], "/"); idx >= 0 {
		modpathlen += idx
	} else {
		modpathlen = len(pkgpath)
	}
	modpath = modprefix
	if major, ok := ModMajor(pkgpath[:modpathlen]); ok {
		modpath = JoinPath(modprefix, major, "")
	}
	pkgdir = strings.TrimPrefix(pkgpath[len(modpath):], "/")
	return modpath, pkgdir, true
}

// SplitSpec splits the path/to/package@query format strings
func SplitSpec(spec string) (path, query string) {
	parts := strings.SplitN(spec, "@", 2)
	if len(parts) == 2 {
		path = parts[0]
		query = parts[1]
	} else {
		path = spec
	}
	return
}

// JoinPath creates a full package path given a module prefix, version, and package directory.
func JoinPath(modprefix, version, pkgdir string) string {
	version = strings.TrimPrefix(version, ".")
	version = strings.TrimPrefix(version, "/")
	major := semver.Major(version)
	pkgpath := modprefix
	switch {
	case strings.HasPrefix(pkgpath, "gopkg.in"):
		pkgpath += "." + major
	case major != "" && major != "v0" && major != "v1" && !strings.Contains(version, "+incompatible"):
		if !strings.HasSuffix(pkgpath, "/") {
			pkgpath += "/"
		}
		pkgpath += major
	}
	if pkgdir != "" {
		pkgpath += "/" + pkgdir
	}
	return pkgpath
}

// FindModFile recursively searches up the directory structure until it
// finds the go.mod, reaches the root of the directory tree, or encounters
// an error.
func FindModFile(dir string) (string, error) {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		name := filepath.Join(dir, "go.mod")
		_, err := os.Stat(name)
		if err == nil {
			return name, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cannot find go.mod")
		}
		dir = parent
	}
}

// Direct returns a list of all modules that are direct dependencies
func Direct(dir string) ([]module.Version, error) {
	name, err := FindModFile(dir)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	file, err := modfile.ParseLax(name, data, nil)
	if err != nil {
		return nil, err
	}
	var mods []module.Version
	for _, req := range file.Require {
		if !req.Indirect {
			mods = append(mods, req.Mod)
		}
	}
	sort.Slice(mods, func(i, j int) bool {
		return mods[i].Path < mods[j].Path
	})
	return mods, nil
}
