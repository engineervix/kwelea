package nav

import (
	"os"
	"path/filepath"
	"testing"
)

// createTempDocs builds a temporary docs directory that exercises the key
// edge cases: numeric-prefix ordering, frontmatter overrides, draft skipping,
// and a one-level-deep subdirectory section.
func createTempDocs(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	files := map[string]string{
		"index.md":                                    "---\ntitle: Home\n---\nWelcome to kwelea.",
		"01-getting-started.md":                       "---\ntitle: Getting Started\n---\nStart here.",
		"02-installation.md":                          "# Installation\nHow to install.",
		"draft-page.md":                               "---\ndraft: true\n---\nThis should never appear.",
		filepath.Join("guide", "index.md"):            "---\ntitle: Guide\n---\nGuide overview.",
		filepath.Join("guide", "01-configuration.md"): "# Configuration\nConfig docs.",
	}

	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// --- frontmatter ---

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFM   Frontmatter
		wantBody string
	}{
		{
			name:     "no frontmatter",
			input:    "# Hello\nWorld",
			wantFM:   Frontmatter{},
			wantBody: "# Hello\nWorld",
		},
		{
			name:  "full frontmatter",
			input: "---\ntitle: Getting Started\ndescription: A guide\ndraft: false\n---\n# Content",
			wantFM: Frontmatter{
				Title: "Getting Started", Description: "A guide", Draft: false,
			},
			wantBody: "# Content",
		},
		{
			name:     "draft page",
			input:    "---\ndraft: true\n---\nSecret",
			wantFM:   Frontmatter{Draft: true},
			wantBody: "Secret",
		},
		{
			name:     "quoted title with colon",
			input:    "---\ntitle: \"My: Guide\"\n---\nBody",
			wantFM:   Frontmatter{Title: "My: Guide"},
			wantBody: "Body",
		},
		{
			name:     "no closing delimiter",
			input:    "---\ntitle: Broken\nBody without closing",
			wantFM:   Frontmatter{},
			wantBody: "---\ntitle: Broken\nBody without closing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fm, body, err := parseFrontmatter(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fm != tc.wantFM {
				t.Errorf("frontmatter: got %+v, want %+v", fm, tc.wantFM)
			}
			if body != tc.wantBody {
				t.Errorf("body:\n  got  %q\n  want %q", body, tc.wantBody)
			}
		})
	}
}

// --- slug helpers ---

func TestStripNumericPrefix(t *testing.T) {
	tests := []struct{ in, want string }{
		{"01-getting-started", "getting-started"},
		{"10-api", "api"},
		{"1-intro", "intro"},
		{"123-deep", "deep"},
		{"guide", "guide"},                   // no prefix — unchanged
		{"no-prefix-here", "no-prefix-here"}, // starts with letters
	}
	for _, tc := range tests {
		if got := stripNumericPrefix(tc.in); got != tc.want {
			t.Errorf("stripNumericPrefix(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSlugToTitle(t *testing.T) {
	tests := []struct{ in, want string }{
		{"getting-started", "Getting Started"},
		{"configuration", "Configuration"},
		{"api-reference", "Api Reference"},
		{"faq", "Faq"},
		{"my_guide", "My Guide"},
	}
	for _, tc := range tests {
		if got := slugToTitle(tc.in); got != tc.want {
			t.Errorf("slugToTitle(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFileToURLPath(t *testing.T) {
	tests := []struct{ in, want string }{
		{"index.md", "/"},
		{"01-getting-started.md", "/getting-started/"},
		{"02-installation.md", "/installation/"},
		{filepath.Join("guide", "index.md"), "/guide/"},
		{filepath.Join("guide", "01-configuration.md"), "/guide/configuration/"},
	}
	for _, tc := range tests {
		if got := fileToURLPath(tc.in); got != tc.want {
			t.Errorf("fileToURLPath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// --- walker ---

func TestWalkDocs(t *testing.T) {
	dir := createTempDocs(t)

	pages, err := WalkDocs(dir)
	if err != nil {
		t.Fatalf("WalkDocs: %v", err)
	}

	// Draft page must never appear.
	for _, p := range pages {
		if p.FilePath == "draft-page.md" {
			t.Error("draft page appeared in results")
		}
	}

	// Expect exactly 5 non-draft pages.
	if len(pages) != 5 {
		t.Errorf("got %d pages, want 5", len(pages))
	}

	// Frontmatter title takes precedence over derived title.
	byFile := map[string]*Page{}
	for _, p := range pages {
		byFile[p.FilePath] = p
	}
	if p := byFile["index.md"]; p == nil || p.Title != "Home" {
		t.Errorf("index.md: got title %q, want %q", byFile["index.md"].Title, "Home")
	}
	if p := byFile["01-getting-started.md"]; p == nil || p.Title != "Getting Started" {
		t.Errorf("01-getting-started.md: got title %q, want %q", p.Title, "Getting Started")
	}

	// Derived title from filename (no frontmatter title).
	if p := byFile["02-installation.md"]; p == nil || p.Title != "Installation" {
		t.Errorf("02-installation.md: got derived title %q, want %q", p.Title, "Installation")
	}

	// URL paths are correct.
	byPath := map[string]*Page{}
	for _, p := range pages {
		byPath[p.Path] = p
	}
	for _, want := range []string{"/", "/getting-started/", "/installation/", "/guide/", "/guide/configuration/"} {
		if byPath[want] == nil {
			t.Errorf("no page at path %q", want)
		}
	}
}

// --- auto nav ---

func TestBuildAutoNav(t *testing.T) {
	dir := createTempDocs(t)
	pages, err := WalkDocs(dir)
	if err != nil {
		t.Fatalf("WalkDocs: %v", err)
	}

	sections := buildAutoNav(pages)

	if len(sections) != 2 {
		t.Fatalf("got %d sections, want 2", len(sections))
	}

	// Root section has an empty label.
	if sections[0].Label != "" {
		t.Errorf("root section label: got %q, want empty string", sections[0].Label)
	}

	// Second section is labelled "Guide".
	if sections[1].Label != "Guide" {
		t.Errorf("guide section label: got %q, want %q", sections[1].Label, "Guide")
	}

	// index.md must be the first item in each section.
	if len(sections[0].Items) == 0 || sections[0].Items[0].Path != "/" {
		t.Errorf("root section: first item path = %q, want /", sections[0].Items[0].Path)
	}
	if len(sections[1].Items) == 0 || sections[1].Items[0].Path != "/guide/" {
		t.Errorf("guide section: first item path = %q, want /guide/", sections[1].Items[0].Path)
	}
}

// --- prev/next ---

func TestLinkPrevNext(t *testing.T) {
	pages := []*Page{
		{Title: "A", Path: "/a/"},
		{Title: "B", Path: "/b/"},
		{Title: "C", Path: "/c/"},
	}
	LinkPrevNext(pages)

	if pages[0].Prev != nil {
		t.Error("first page: Prev should be nil")
	}
	if pages[0].Next != pages[1] {
		t.Error("first page: Next should be second page")
	}
	if pages[1].Prev != pages[0] {
		t.Error("middle page: Prev should be first page")
	}
	if pages[1].Next != pages[2] {
		t.Error("middle page: Next should be last page")
	}
	if pages[2].Next != nil {
		t.Error("last page: Next should be nil")
	}
}

// --- lookupPageByRef (implements the user-contributed function) ---

// TestLookupPageByRef validates the three-tier matching logic and the
// ambiguous-slug error.
func TestLookupPageByRef(t *testing.T) {
	pages := []*Page{
		{FilePath: "index.md", Path: "/"},
		{FilePath: "01-getting-started.md", Path: "/getting-started/"},
		{FilePath: filepath.Join("guide", "01-configuration.md"), Path: "/guide/configuration/"},
	}

	tests := []struct {
		ref      string
		wantPath string // empty string means expect nil (not found)
	}{
		// 1. Exact FilePath matches
		{"index.md", "/"},
		{filepath.Join("guide", "01-configuration.md"), "/guide/configuration/"},
		// 2. Bare filename matches (first match wins)
		{"01-getting-started.md", "/getting-started/"},
		// 3. Slug-only convenience form
		{"getting-started", "/getting-started/"},
		{"configuration", "/guide/configuration/"},
		// 4. Not found
		{"nonexistent.md", ""},
	}

	for _, tc := range tests {
		t.Run(tc.ref, func(t *testing.T) {
			p, err := lookupPageByRef(pages, tc.ref)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantPath == "" {
				if p != nil {
					t.Errorf("expected nil, got page at %q", p.Path)
				}
				return
			}
			if p == nil {
				t.Errorf("got nil, want page at %q", tc.wantPath)
				return
			}
			if p.Path != tc.wantPath {
				t.Errorf("got path %q, want %q", p.Path, tc.wantPath)
			}
		})
	}

	// Ambiguous slug — two pages share the same slug in different directories.
	// The build must fail with an error naming both files.
	t.Run("ambiguous_slug", func(t *testing.T) {
		ambiguous := []*Page{
			{FilePath: filepath.Join("guide", "01-configuration.md"), Path: "/guide/configuration/"},
			{FilePath: filepath.Join("api", "01-configuration.md"), Path: "/api/configuration/"},
		}
		_, err := lookupPageByRef(ambiguous, "configuration")
		if err == nil {
			t.Error("expected an error for ambiguous slug, got nil")
		}
	})
}
