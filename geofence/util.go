package geofence

import (
	"path/filepath"
	"regexp"
	"strings"
)

var slugger = regexp.MustCompile("[^a-z0-9]+")

// Slugs the basename of the path, removing the path and extension
// "/path/to/file_2.gz " -> "file-2"
// yoinked from diglet/util
func slug(path string) string {
	s := filepath.Base(path)
	s = strings.TrimSuffix(s, filepath.Ext(s))
	return slugged(s, "-")
}

func slugged(s, delim string) string {
	return strings.Trim(slugger.ReplaceAllString(strings.ToLower(s), delim), delim)
}
