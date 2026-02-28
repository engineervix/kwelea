package nav

import (
	"path/filepath"
	"strings"
	"unicode"
)

// stripNumericPrefix removes a leading "NN-" ordering prefix from a path segment.
//
//	"01-getting-started" → "getting-started"
//	"10-api"             → "api"
//	"guide"              → "guide"  (unchanged — no numeric prefix)
func stripNumericPrefix(s string) string {
	for i, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '-' && i > 0 {
			return s[i+1:]
		}
		break
	}
	return s
}

// slugToTitle converts a URL slug into a human-readable page title.
// It splits on hyphens and underscores and title-cases the first letter of
// each word.
//
//	"getting-started" → "Getting Started"
//	"api-reference"   → "Api Reference"
//	"configuration"   → "Configuration"
func slugToTitle(slug string) string {
	words := strings.FieldsFunc(slug, func(r rune) bool {
		return r == '-' || r == '_'
	})
	for i, w := range words {
		if w == "" {
			continue
		}
		runes := []rune(w)
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// fileToURLPath derives the canonical URL path for a .md source file.
// relPath is the file's path relative to docs_dir, using OS path separators.
//
//	"index.md"                  → "/"
//	"01-getting-started.md"     → "/getting-started/"
//	"guide/index.md"            → "/guide/"
//	"guide/01-configuration.md" → "/guide/configuration/"
func fileToURLPath(relPath string) string {
	slashed := filepath.ToSlash(relPath)
	parts := strings.Split(slashed, "/")
	base := parts[len(parts)-1]
	dirs := parts[:len(parts)-1]

	// Strip numeric ordering prefixes from every directory segment.
	urlDirs := make([]string, 0, len(dirs))
	for _, d := range dirs {
		urlDirs = append(urlDirs, stripNumericPrefix(d))
	}

	if base == "index.md" {
		if len(urlDirs) == 0 {
			return "/"
		}
		return "/" + strings.Join(urlDirs, "/") + "/"
	}

	slug := stripNumericPrefix(strings.TrimSuffix(base, ".md"))
	if len(urlDirs) == 0 {
		return "/" + slug + "/"
	}
	return "/" + strings.Join(urlDirs, "/") + "/" + slug + "/"
}

// derivedTitle returns a display title for a file whose frontmatter has no
// explicit title set.
func derivedTitle(relPath string) string {
	slashed := filepath.ToSlash(relPath)
	parts := strings.Split(slashed, "/")
	base := parts[len(parts)-1]

	if base == "index.md" {
		if len(parts) == 1 {
			return "Home" // root index.md
		}
		// Use the parent directory name as the section title.
		return slugToTitle(stripNumericPrefix(parts[len(parts)-2]))
	}
	return slugToTitle(stripNumericPrefix(strings.TrimSuffix(base, ".md")))
}

// firstDir returns the first directory component of a docs-relative file path.
// Files at the root of docs_dir return "". Deeply nested files return only
// their top-level directory ("guide/advanced/deep.md" → "guide").
func firstDir(relPath string) string {
	d := filepath.ToSlash(filepath.Dir(relPath))
	if d == "." || d == "" {
		return ""
	}
	if idx := strings.Index(d, "/"); idx >= 0 {
		return d[:idx]
	}
	return d
}
