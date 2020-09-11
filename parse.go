package lionrouter

import (
	"strings"
)

type pathKey struct {
	name         string
	wildcard     bool
	wildcardPath bool
}

// splitPath splits a given data path with '/' as path separator
// splitPath ignores trailing slashes
// returns split path or error
func splitPath(path string) []string {
	// check for root
	if path == "" || path == "/" {
		return []string{""}
	}

	// split by path separator '/'
	split := strings.Split(path, "/")

	// ignore leading slash
	if len(split) > 1 && split[0] == "" {
		split = split[1:]
	}

	// ignore trailing slash
	if len(split) > 1 && split[len(split)-1] == "" {
		split = split[:len(split)-1]
	}

	return split
}

// parsePath parses a web path
// returns an array of pathKey or error
func parsePath(path string) []*pathKey {
	// split path
	split := splitPath(path)

	// parse each element
	keys := make([]*pathKey, len(split))

	for i, k := range split {
		keys[i] = parseKey(k)
	}

	return keys
}

// parseKey parses a single key
// it returns the parsed key as pathKey
func parseKey(key string) *pathKey {
	parsed := &pathKey{key, false, false}

	// return immediately, trivial key found
	if key == "" || len(key) == 1 {
		return parsed
	}

	// check for wildcard options
	leadingChar := key[0]

	if leadingChar == ':' {
		parsed.name = key[1:]
		parsed.wildcard = true
	} else if leadingChar == '*' {
		parsed.name = key[1:]
		parsed.wildcardPath = true
	}

	return parsed
}
