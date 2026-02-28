package builder

import (
	"encoding/json"
	"fmt"
	"html"
	html_template "html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/engineervix/kwelea/internal/nav"
)

// SearchEntry is one document in the search index written to search-index.json.
// The JS client loads this file and builds an in-memory FlexSearch index from it.
type SearchEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
	Body  string `json:"body"`
}

var (
	htmlTagRE = regexp.MustCompile(`<[^>]+>`)
	spaceRE   = regexp.MustCompile(`\s+`)
)

// writeSearchIndex generates search-index.json from all parsed pages and
// writes it to outputDir. It must be called after Pass 1 (parse) so that
// every Page.HTML is populated.
func writeSearchIndex(pages []*nav.Page, outputDir string) error {
	entries := make([]SearchEntry, 0, len(pages))
	for _, p := range pages {
		entries = append(entries, SearchEntry{
			ID:    pathToSearchID(p.Path),
			Title: p.Title,
			Path:  p.Path,
			Body:  plainText(p.HTML),
		})
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("marshalling search index: %w", err)
	}
	outPath := filepath.Join(outputDir, "search-index.json")
	return os.WriteFile(outPath, data, 0o644)
}

// plainText strips HTML tags from rendered page HTML and returns clean,
// collapsed plain text suitable for full-text indexing.
func plainText(h html_template.HTML) string {
	s := htmlTagRE.ReplaceAllString(string(h), " ")
	s = html.UnescapeString(s)
	return strings.TrimSpace(spaceRE.ReplaceAllString(s, " "))
}

// pathToSearchID converts a URL path to a stable, URL-safe document ID.
// "/"                       → "index"
// "/getting-started/"       → "getting-started"
// "/guide/configuration/"   → "guide-configuration"
func pathToSearchID(path string) string {
	s := strings.Trim(path, "/")
	if s == "" {
		return "index"
	}
	return strings.ReplaceAll(s, "/", "-")
}
