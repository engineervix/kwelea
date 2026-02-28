package parser

import (
	"bytes"
	"fmt"
	"html/template"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/engineervix/kwelea/internal/config"
	"github.com/engineervix/kwelea/internal/nav"
)

// Parse processes a Markdown source file and returns the rendered HTML body,
// a table of contents extracted from h2/h3 headings, and any error.
//
// filePath is used for error messages only. themeCfg selects the Chroma style
// names written into CSS classes (see ChromaCSS for the matching stylesheet).
func Parse(filePath string, src []byte, themeCfg config.ThemeConfig) (template.HTML, []nav.TocItem, error) {
	// Strip YAML frontmatter so goldmark does not misparse "---" as a setext rule.
	body := stripFrontmatter(src)

	md := newMarkdown(themeCfg)
	reader := text.NewReader(body)
	doc := md.Parser().Parse(reader)

	// ToC must be extracted between Parse and Render because AutoHeadingID
	// sets id attributes on the AST nodes; they are accessible here but are
	// written to the HTML output — not re-readable — after rendering.
	toc := extractTOC(doc, body)

	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, body, doc); err != nil {
		return "", nil, fmt.Errorf("rendering %s: %w", filePath, err)
	}

	return template.HTML(buf.String()), toc, nil
}

// newMarkdown returns a goldmark.Markdown configured with all kwelea extensions:
//   - GFM (tables, strikethrough, linkify, task lists)
//   - Syntax highlighting using Chroma CSS classes (dual-theme via ChromaCSS)
//   - Admonitions (:::)
//   - D2 diagrams (```d2 fenced blocks)
//   - Auto-heading IDs for ToC extraction
func newMarkdown(themeCfg config.ThemeConfig) goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				// WithStyle names the fallback style for inline rendering; since we
				// use WithClasses(true) the actual colours come from ChromaCSS output.
				highlighting.WithStyle(themeCfg.LightCodeTheme),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
			Admonitions,
			NewD2Extension(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			// Allow raw HTML in Markdown source — authors of documentation sites
			// are trusted; this mirrors the behaviour of most doc generators.
			goldmarkhtml.WithUnsafe(),
		),
	)
}

// stripFrontmatter removes the leading YAML frontmatter block (if any) from
// the Markdown source so goldmark does not misparse "---" as a setext rule.
func stripFrontmatter(src []byte) []byte {
	if !bytes.HasPrefix(src, []byte("---\n")) {
		return src
	}
	rest := src[4:] // skip opening "---\n"
	idx := bytes.Index(rest, []byte("\n---"))
	if idx < 0 {
		return src // no closing delimiter — treat whole file as body
	}
	body := rest[idx+4:] // skip "\n---"
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}
	return body
}
