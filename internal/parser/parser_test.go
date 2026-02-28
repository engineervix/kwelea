package parser

import (
	"strings"
	"testing"

	"github.com/engineervix/kwelea/internal/config"
)

var defaultTheme = config.ThemeConfig{
	LightCodeTheme: "github",
	DarkCodeTheme:  "github-dark",
}

// ----- Parse() integration tests -----

func TestParseBasicMarkdown(t *testing.T) {
	src := []byte("# Hello\n\nThis is **bold** and _italic_.\n")
	html, toc, h1, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	// H1 is stripped from the body; the template renders it via Page.Title.
	if strings.Contains(string(html), "<h1") {
		t.Errorf("expected <h1> to be stripped from output, got:\n%s", html)
	}
	if h1 != "Hello" {
		t.Errorf("expected h1Title %q, got %q", "Hello", h1)
	}
	if !strings.Contains(string(html), "<strong>bold</strong>") {
		t.Errorf("expected <strong>bold</strong>, got:\n%s", html)
	}
	if len(toc) != 0 {
		t.Errorf("expected no ToC items for h1-only doc, got %d", len(toc))
	}
}

func TestParseH1Stripped(t *testing.T) {
	// H1 stripped; h2 onward untouched.
	src := []byte("# Page Title\n\n## Section\n\nContent.\n")
	html, toc, h1, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if strings.Contains(string(html), "<h1") {
		t.Errorf("h1 should be stripped from output, got:\n%s", html)
	}
	if h1 != "Page Title" {
		t.Errorf("expected h1Title %q, got %q", "Page Title", h1)
	}
	if !strings.Contains(string(html), "<h2") {
		t.Errorf("h2 should still be present, got:\n%s", html)
	}
	if len(toc) != 1 || toc[0].Text != "Section" {
		t.Errorf("expected one ToC item 'Section', got %+v", toc)
	}
}

func TestParseNoH1ReturnsEmptyTitle(t *testing.T) {
	src := []byte("## Section\n\nContent.\n")
	_, _, h1, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if h1 != "" {
		t.Errorf("expected empty h1Title for doc with no H1, got %q", h1)
	}
}

func TestParseFrontmatterStripped(t *testing.T) {
	src := []byte("---\ntitle: My Page\n---\n\n# My Page\n\nContent.\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if strings.Contains(string(html), "---") {
		t.Errorf("frontmatter delimiters should be stripped, got:\n%s", html)
	}
	if strings.Contains(string(html), "title: My Page") {
		t.Errorf("frontmatter content should be stripped, got:\n%s", html)
	}
}

// ----- Admonitions -----

func TestAdmonitionInfo(t *testing.T) {
	src := []byte("::: info\nThis is an info block.\n:::\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `class="admonition admonition-info"`) {
		t.Errorf("expected admonition-info class, got:\n%s", s)
	}
	if !strings.Contains(s, `class="admonition-label"`) {
		t.Errorf("expected admonition-label div, got:\n%s", s)
	}
	if !strings.Contains(s, "Info") {
		t.Errorf("expected label text 'Info', got:\n%s", s)
	}
	if !strings.Contains(s, "This is an info block.") {
		t.Errorf("expected body text in output, got:\n%s", s)
	}
}

func TestAdmonitionWarning(t *testing.T) {
	src := []byte("::: warning\nWatch out!\n:::\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, "admonition-warning") {
		t.Errorf("expected admonition-warning class, got:\n%s", s)
	}
	if !strings.Contains(s, "Warning") {
		t.Errorf("expected label 'Warning', got:\n%s", s)
	}
}

func TestAdmonitionDetails(t *testing.T) {
	src := []byte("::: details Click to expand\nHidden content.\n:::\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, "<details") {
		t.Errorf("expected <details> element, got:\n%s", s)
	}
	if !strings.Contains(s, "<summary>Click to expand</summary>") {
		t.Errorf("expected <summary> with title text, got:\n%s", s)
	}
	if !strings.Contains(s, "Hidden content.") {
		t.Errorf("expected hidden content in output, got:\n%s", s)
	}
}

func TestAdmonitionDetailsDefaultTitle(t *testing.T) {
	src := []byte("::: details\nContent here.\n:::\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !strings.Contains(string(html), "<summary>Details</summary>") {
		t.Errorf("expected default summary 'Details', got:\n%s", html)
	}
}

func TestAdmonitionUnknownTypeIgnored(t *testing.T) {
	// Unknown types should not be parsed as admonitions.
	src := []byte("::: unknown\nContent.\n:::\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if strings.Contains(string(html), "admonition-unknown") {
		t.Errorf("unknown admonition type should not produce admonition HTML, got:\n%s", html)
	}
}

func TestAdmonitionNestedMarkdown(t *testing.T) {
	// Content inside admonitions is parsed as Markdown.
	src := []byte("::: tip\n**bold** and `code`\n:::\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, "<strong>bold</strong>") {
		t.Errorf("expected nested <strong>, got:\n%s", s)
	}
	if !strings.Contains(s, "<code>code</code>") {
		t.Errorf("expected nested <code>, got:\n%s", s)
	}
}

// ----- ToC extraction -----

func TestToCExtractionH2H3(t *testing.T) {
	src := []byte("## Section One\n\n### Sub-section\n\n## Section Two\n")
	_, toc, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(toc) != 3 {
		t.Fatalf("expected 3 ToC items, got %d: %+v", len(toc), toc)
	}
	if toc[0].Level != 2 || toc[0].Text != "Section One" {
		t.Errorf("unexpected toc[0]: %+v", toc[0])
	}
	if toc[1].Level != 3 || toc[1].Text != "Sub-section" {
		t.Errorf("unexpected toc[1]: %+v", toc[1])
	}
	if toc[2].Level != 2 || toc[2].Text != "Section Two" {
		t.Errorf("unexpected toc[2]: %+v", toc[2])
	}
}

func TestToCH1Excluded(t *testing.T) {
	src := []byte("# Title\n\n## Section\n")
	_, toc, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	for _, item := range toc {
		if item.Level == 1 {
			t.Errorf("h1 headings should not appear in ToC: %+v", item)
		}
	}
	if len(toc) != 1 || toc[0].Text != "Section" {
		t.Errorf("expected exactly the h2 'Section', got %+v", toc)
	}
}

func TestToCIDs(t *testing.T) {
	src := []byte("## Getting Started\n\n## FAQ\n")
	_, toc, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	for _, item := range toc {
		if item.ID == "" {
			t.Errorf("ToC item %q has empty ID (AutoHeadingID not active?)", item.Text)
		}
	}
}

// ----- Syntax highlighting -----

func TestHighlightingUsesClasses(t *testing.T) {
	src := []byte("```go\nfunc main() {}\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	// CSS-class mode: spans have class attributes, not style attributes.
	if strings.Contains(s, `style="color`) {
		t.Errorf("expected CSS classes not inline styles, but found style attr in:\n%s", s)
	}
	if !strings.Contains(s, `class="`) {
		t.Errorf("expected class attributes on highlighted spans, got:\n%s", s)
	}
}

// ----- ChromaCSS -----

func TestChromaCSSContainsBothThemes(t *testing.T) {
	css, err := ChromaCSS(defaultTheme)
	if err != nil {
		t.Fatalf("ChromaCSS error: %v", err)
	}
	if !strings.Contains(css, ".chroma") {
		t.Errorf("expected .chroma selector in light CSS, got:\n%s", css[:min(200, len(css))])
	}
	if !strings.Contains(css, `[data-theme="dark"]`) {
		t.Errorf("expected [data-theme=\"dark\"] scoped rules in CSS, got:\n%s", css[:min(200, len(css))])
	}
}

func TestChromaCSSFallbackOnUnknownStyle(t *testing.T) {
	cfg := config.ThemeConfig{
		LightCodeTheme: "nonexistent-style",
		DarkCodeTheme:  "another-nonexistent-style",
	}
	css, err := ChromaCSS(cfg)
	if err != nil {
		t.Fatalf("ChromaCSS should fall back gracefully, got error: %v", err)
	}
	if css == "" {
		t.Error("ChromaCSS should return non-empty CSS even for unknown styles (fallback)")
	}
}

// ----- D2 diagrams -----

func TestD2DiagramRenders(t *testing.T) {
	src := []byte("```d2\ndirection: right\nuser -> server\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `class="d2-diagram"`) {
		t.Errorf("expected d2-diagram wrapper, got:\n%s", s)
	}
	if !strings.Contains(s, `class="d2-light"`) {
		t.Errorf("expected d2-light div, got:\n%s", s)
	}
	if !strings.Contains(s, `class="d2-dark"`) {
		t.Errorf("expected d2-dark div, got:\n%s", s)
	}
	if !strings.Contains(s, "<svg") {
		t.Errorf("expected inline SVG in output, got:\n%s", s)
	}
}

func TestD2InvalidSourceShowsError(t *testing.T) {
	src := []byte("```d2\n{{{invalid d2 syntax\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse should not propagate D2 errors, got: %v", err)
	}
	if !strings.Contains(string(html), "d2-error") {
		t.Errorf("expected d2-error class for invalid D2, got:\n%s", html)
	}
}

// min is a small helper kept local to avoid Go version concerns.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
