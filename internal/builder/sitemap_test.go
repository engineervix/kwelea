package builder

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/engineervix/kwelea/internal/nav"
)

func TestHasScheme(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"https://example.com", true},
		{"http://localhost:4000", true},
		{"https://engineervix.github.io/kwelea", true},
		{"", false},
		{"/kwelea", false},
		{"kwelea", false},
		{"ftp://files.example.com", false},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			if got := hasScheme(tc.in); got != tc.want {
				t.Errorf("hasScheme(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestWriteSitemap(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		wantCount  int
		wantWrite  bool // whether sitemap.xml should exist
		wantLocs   []string
	}{
		{
			name:     "https base_url with pages",
			baseURL:  "https://example.com",
			wantCount: 3,
			wantWrite: true,
			wantLocs: []string{
				"https://example.com/",
				"https://example.com/guide/",
				"https://example.com/getting-started/",
			},
		},
		{
			name:     "base_url with trailing slash",
			baseURL:  "https://example.com/kwelea/",
			wantCount: 2,
			wantWrite: true,
			wantLocs: []string{
				"https://example.com/kwelea/",
				"https://example.com/kwelea/guide/",
			},
		},
		{
			name:     "no base_url — skip silently",
			baseURL:  "",
			wantCount: 0,
			wantWrite: false,
		},
		{
			name:     "relative base_url — skip silently",
			baseURL:  "/kwelea",
			wantCount: 0,
			wantWrite: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pages := []*nav.Page{
				{Title: "Home", Path: "/"},
				{Title: "Guide", Path: "/guide/"},
				{Title: "Getting Started", Path: "/getting-started/"},
			}

			site := &nav.Site{
				BaseURL: tc.baseURL,
				Pages:   pages,
			}

			if tc.wantCount > 0 {
				// Truncate the page list to match wantCount for the sub-tests.
				site.Pages = pages[:tc.wantCount]
			}

			outDir := t.TempDir()

			if err := writeSitemap(site, outDir); err != nil {
				t.Fatalf("writeSitemap: %v", err)
			}

			sitemapPath := filepath.Join(outDir, "sitemap.xml")
			_, err := os.Stat(sitemapPath)
			if tc.wantWrite {
				if err != nil {
					t.Fatalf("expected sitemap.xml to exist: %v", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected sitemap.xml NOT to exist")
				}
				return // nothing more to check
			}

			// Parse and validate the generated XML.
			data, err := os.ReadFile(sitemapPath)
			if err != nil {
				t.Fatalf("reading sitemap.xml: %v", err)
			}

			var urlSet sitemapURLSet
			if err := xml.Unmarshal(data, &urlSet); err != nil {
				t.Fatalf("unmarshalling sitemap: %v", err)
			}

			// Validate xmlns.
			if urlSet.XMLNS != "http://www.sitemaps.org/schemas/sitemap/0.9" {
				t.Errorf("xmlns: got %q, want %q", urlSet.XMLNS, "http://www.sitemaps.org/schemas/sitemap/0.9")
			}

			// Validate URL count.
			if len(urlSet.URLs) != tc.wantCount {
				t.Fatalf("got %d <url> entries, want %d", len(urlSet.URLs), tc.wantCount)
			}

			// Validate <loc> values.
			gotLocs := make([]string, len(urlSet.URLs))
			for i, u := range urlSet.URLs {
				gotLocs[i] = u.Loc
			}
			for i, want := range tc.wantLocs {
				if i >= len(gotLocs) {
					t.Errorf("missing <loc> entry %d: want %q", i, want)
					continue
				}
				if gotLocs[i] != want {
					t.Errorf("loc[%d]: got %q, want %q", i, gotLocs[i], want)
				}
			}
		})
	}
}

func TestWriteSitemapExcludesDrafts(t *testing.T) {
	site := &nav.Site{
		BaseURL: "https://example.com",
		Pages: []*nav.Page{
			{Title: "Home", Path: "/", Draft: false},
			{Title: "Secret", Path: "/secret/", Draft: true},
			{Title: "Guide", Path: "/guide/", Draft: false},
		},
	}

	outDir := t.TempDir()

	if err := writeSitemap(site, outDir); err != nil {
		t.Fatalf("writeSitemap: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "sitemap.xml"))
	if err != nil {
		t.Fatalf("reading sitemap.xml: %v", err)
	}

	var urlSet sitemapURLSet
	if err := xml.Unmarshal(data, &urlSet); err != nil {
		t.Fatalf("unmarshalling sitemap: %v", err)
	}

	if len(urlSet.URLs) != 2 {
		t.Fatalf("got %d <url> entries, want 2 (draft must be excluded)", len(urlSet.URLs))
	}

	for _, u := range urlSet.URLs {
		if u.Loc == "https://example.com/secret/" {
			t.Error("draft page appeared in sitemap")
		}
	}
}