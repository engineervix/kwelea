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

// ----- Code block title and line highlights (#6) -----

func TestCodeBlockTitleOnly(t *testing.T) {
	src := []byte("```go title=\"cmd/root.go\"\npackage main\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `<figure class="code-block" data-lang="go">`) {
		t.Errorf("expected <figure class=\"code-block\" data-lang=\"go\">, got:\n%s", s)
	}
	if !strings.Contains(s, `<div class="code-title">cmd/root.go</div>`) {
		t.Errorf("expected code-title div with filename, got:\n%s", s)
	}
	// No highlight-line should appear in a title-only block.
	if strings.Contains(s, "highlight-line") {
		t.Errorf("title-only block should not contain highlight-line, got:\n%s", s)
	}
}

func TestCodeBlockLineHighlightOnly(t *testing.T) {
	src := []byte("```go {2}\nline1\nline2\nline3\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `<figure class="code-block" data-lang="go">`) {
		t.Errorf("expected code-block wrapper even without title, got:\n%s", s)
	}
	if strings.Contains(s, "code-title") {
		t.Errorf("highlight-only block should not contain code-title, got:\n%s", s)
	}
	// Exactly one line should carry the highlight-line class (line 2).
	count := strings.Count(s, `class="line highlight-line`)
	if count != 1 {
		t.Errorf("expected exactly 1 highlight-line span, got %d:\n%s", count, s)
	}
	// And the highlighted line should contain the content of line 2. We
	// don't assert the exact class string because chroma may also have
	// added " hl" (line numbers in table mode), producing
	// "line highlight-line hl".
	if !strings.Contains(s, "line2") {
		t.Errorf("expected line2 to appear in the highlighted line, got:\n%s", s)
	}
}

func TestCodeBlockTitleAndLineHighlights(t *testing.T) {
	// Reproduces the example from issue #6.
	src := []byte("```go title=\"cmd/root.go\" {2,4-6}\nline 1\nline 2\nline 3\nline 4\nline 5\nline 6\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `<figure class="code-block" data-lang="go">`) {
		t.Errorf("expected code-block wrapper, got:\n%s", s)
	}
	if !strings.Contains(s, `<div class="code-title">cmd/root.go</div>`) {
		t.Errorf("expected code-title, got:\n%s", s)
	}
	// 4 lines should be highlighted: 2 and 4, 5, 6.
	count := strings.Count(s, `class="line highlight-line`)
	if count != 4 {
		t.Errorf("expected 4 highlight-line spans, got %d:\n%s", count, s)
	}
}

func TestCodeBlockPlainPassthrough(t *testing.T) {
	// A code block with no title and no {} should NOT be wrapped in
	// .code-block — the existing goldmark-highlighting path renders it
	// unchanged, so the common case pays zero overhead.
	src := []byte("```go\nfunc main() {}\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if strings.Contains(s, "code-block") {
		t.Errorf("plain code block should not be wrapped, got:\n%s", s)
	}
	if !strings.Contains(s, `<pre class="chroma">`) {
		t.Errorf("plain code block should still be highlighted, got:\n%s", s)
	}
}

func TestCodeBlockInvalidRangePartiallyParsed(t *testing.T) {
	// An invalid token in the range is dropped; valid ones still apply.
	// {notanumber,5} should yield just line 5.
	src := []byte("```go {notanumber,5}\nx\ny\nz\nw\nv\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	count := strings.Count(s, `class="line highlight-line`)
	if count != 1 {
		t.Errorf("expected 1 highlight-line span (only the valid '5'), got %d:\n%s", count, s)
	}
	// Line 5 is the last one ("v"). It lives between <span class="nx">
	// and </span>, so we just check for the character itself.
	if !strings.Contains(s, ">v<") {
		t.Errorf("expected the highlighted line to contain 'v', got:\n%s", s)
	}
}

func TestCodeBlockTitleEscapesHTML(t *testing.T) {
	// The title must be HTML-escaped; a raw "&" or "<" would break the
	// document. The wrap-the-attribute-in-quotes form ensures the parser
	// sees the literal string; the renderer must escape it.
	cases := []struct {
		name, title, wantSubstr string
	}{
		{"ampersand", `build & deploy.sh`, `build &amp; deploy.sh`},
		// Single-quoted title wrapping lets us embed double quotes
		// without escaping.
		{"single-quotes-around-double", `has "quotes" in it`, `has &#34;quotes&#34; in it`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Pick the quote style based on whether the title contains
			// the same kind of quote — using the other kind avoids the
			// need for backslash-escaping the inner one.
			quote := `"`
			if strings.Contains(c.title, `"`) {
				quote = `'`
			}
			src := []byte("```sh title=" + quote + c.title + quote + "\necho hi\n```\n")
			html, _, _, err := Parse("test.md", src, defaultTheme)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			s := string(html)
			if !strings.Contains(s, c.wantSubstr) {
				t.Errorf("expected escaped title %q in output, got:\n%s", c.wantSubstr, s)
			}
			// And the unescaped form must not appear in the title div.
			if idx := strings.Index(s, `class="code-title">`); idx >= 0 {
				end := strings.Index(s[idx:], `</div>`)
				if end < 0 {
					t.Fatalf("malformed output: no </div> after code-title")
				}
				titleHTML := s[idx : idx+end]
				if strings.Contains(titleHTML, c.title) {
					t.Errorf("raw title %q appeared unescaped in %q", c.title, titleHTML)
				}
			}
		})
	}
}

func TestCodeBlockHighlightOutOfRange(t *testing.T) {
	// Asking for line 99 on a 1-line block: the wrapper still renders,
	// the block is still highlighted, but no line gets the class. The
	// behaviour matches what a user would expect from a typo — silent
	// ignore rather than a broken page.
	src := []byte("```go {99}\nonly-one-line\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `<figure class="code-block"`) {
		t.Errorf("expected code-block wrapper, got:\n%s", s)
	}
	if strings.Contains(s, "highlight-line") {
		t.Errorf("out-of-range highlight should produce no highlight-line, got:\n%s", s)
	}
}

func TestCodeBlockEmptyTitleFallsThrough(t *testing.T) {
	// title="" is a no-op — treat the block as plain (no wrapper).
	src := []byte("```go title=\"\"\nfoo\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if strings.Contains(string(html), "code-block") {
		t.Errorf("empty title should not wrap the block, got:\n%s", html)
	}
}

func TestCodeBlockD2StillWorks(t *testing.T) {
	// D2 blocks are caught by the d2 transformer before ours sees them.
	// Make sure adding title=… or {…} to a d2 block does not break
	// diagram rendering.
	src := []byte("```d2\ndirection: right\na -> b\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `class="d2-diagram"`) {
		t.Errorf("expected d2-diagram wrapper, got:\n%s", s)
	}
	if strings.Contains(s, "code-block") {
		t.Errorf("d2 block should not be wrapped in .code-block, got:\n%s", s)
	}
}

func TestCodeBlockRangeWithLeadingSpace(t *testing.T) {
	// The {n,m-p} expression is preceded by whitespace; the parser must
	// accept that.
	src := []byte("```go   {2,4-5}\n1\n2\n3\n4\n5\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	count := strings.Count(s, `class="line highlight-line`)
	if count != 3 {
		t.Errorf("expected 3 highlight-line spans (lines 2, 4, 5), got %d:\n%s", count, s)
	}
}

func TestCodeBlockTitleBeforeBraces(t *testing.T) {
	// Order shouldn't matter.
	src := []byte("```go {1} title=\"x.go\"\nfoo\nbar\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `<div class="code-title">x.go</div>`) {
		t.Errorf("expected code-title regardless of attribute order, got:\n%s", s)
	}
	if !strings.Contains(s, "highlight-line") {
		t.Errorf("expected highlight-line regardless of attribute order, got:\n%s", s)
	}
}

func TestCodeBlockReverseRangeIgnored(t *testing.T) {
	// {5-2} is an empty range; the whole expression should be dropped
	// silently rather than raising an error.
	src := []byte("```go {5-2}\n1\n2\n3\n4\n5\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if strings.Contains(string(html), "highlight-line") {
		t.Errorf("reverse range should produce no highlights, got:\n%s", html)
	}
}

func TestCodeBlockNoLanguageStillWorks(t *testing.T) {
	// A block with no language but with a title should still render the
	// wrapper. Chroma falls back to the plain-text lexer.
	src := []byte("```title=\"snippet.txt\"\nplain text\n```\n")
	html, _, _, err := Parse("test.md", src, defaultTheme)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `data-lang=""`) && !strings.Contains(s, `<figure class="code-block">`) {
		t.Errorf("expected code-block wrapper, got:\n%s", s)
	}
	if !strings.Contains(s, `code-title">snippet.txt`) {
		t.Errorf("expected code-title for language-less block, got:\n%s", s)
	}
}

// min is a small helper kept local to avoid Go version concerns.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
