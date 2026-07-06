package parser

import (
	"bytes"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/engineervix/kwelea/internal/config"
)

// ChromaCSS generates a combined Chroma syntax-highlighting CSS string for
// both the light and dark themes specified in themeCfg.
//
// The light theme uses standard .chroma selectors. Dark-theme rules are
// prefixed with [data-theme="dark"] so they activate only when that attribute
// is set on <html>, matching the toggle logic in the template.
func ChromaCSS(themeCfg config.ThemeConfig) (string, error) {
	formatter := chromahtml.New(chromahtml.WithClasses(true))

	lightStyle := styles.Get(themeCfg.LightCodeTheme)
	if lightStyle == nil {
		lightStyle = styles.Fallback
	}
	darkStyle := styles.Get(themeCfg.DarkCodeTheme)
	if darkStyle == nil {
		darkStyle = styles.Fallback
	}

	var buf bytes.Buffer

	// Light theme — standard selectors.
	if err := formatter.WriteCSS(&buf, lightStyle); err != nil {
		return "", fmt.Errorf("generating light Chroma CSS (%s): %w", themeCfg.LightCodeTheme, err)
	}

	buf.WriteString("\n")

	// Dark theme — prefix every selector with [data-theme="dark"].
	var darkBuf bytes.Buffer
	if err := formatter.WriteCSS(&darkBuf, darkStyle); err != nil {
		return "", fmt.Errorf("generating dark Chroma CSS (%s): %w", themeCfg.DarkCodeTheme, err)
	}
	buf.WriteString(prefixCSSSelectors(darkBuf.String(), `[data-theme="dark"]`))

	return buf.String(), nil
}

// prefixCSSSelectors prepends prefix to every CSS selector in chroma formatter
// output. Chroma emits one rule per line in the form:
//
//	/* Token Name */ .chroma { ... }
//	/* Token Name */ .chroma .tok { ... }
//
// so finding the first '.' on each non-blank line reliably locates the selector.
func prefixCSSSelectors(css, prefix string) string {
	lines := strings.Split(css, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if i := strings.Index(line, "."); i >= 0 {
			result = append(result, line[:i]+prefix+" "+line[i:])
		} else {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// =============================================================================
// Code-block title and line highlighting (closes #6)
// =============================================================================
//
// Extends fenced code blocks with two extra info-string attributes:
//
//	go title="cmd/root.go"          // filename label above the block
//	go {2,4-6}                      // highlight lines 2 and 4-6
//	go title="x.go" {2,4-6}         // both
//
// Implementation summary:
//
//   - parseCodeBlockInfo() reads the info string after the language and
//     returns a codeBlockAttrs struct with the parsed title and highlight
//     ranges. Invalid expressions are silently ignored so that a stray
//     character does not break the rest of the block.
//
//   - codeAttrsTransformer runs as an AST transformer and replaces any
//     FencedCodeBlock that has at least one custom attribute with a new
//     CodeBlockNode. Plain blocks (no title, no highlights) are left
//     untouched so goldmark-highlighting still renders them with zero
//     overhead — the common case is unchanged.
//
//   - codeAttrsNodeRenderer renders CodeBlockNode by formatting the source
//     with Chroma in line-numbered table mode, post-processing the output
//     to add class="highlight-line" to the matching <span class="line">
//     rows, and wrapping the whole thing in <figure class="code-block">
//     with an optional <div class="code-title"> above the <pre>.

// codeBlockAttrs holds the non-language attributes parsed from a fenced
// code block's info string.
type codeBlockAttrs struct {
	title        string
	highlight    []int // 1-indexed line numbers to highlight
	highlightSet map[int]bool
}

// hasTitle reports whether a non-empty title was set.
func (a codeBlockAttrs) hasTitle() bool { return a.title != "" }

// hasHighlights reports whether at least one line highlight was set.
func (a codeBlockAttrs) hasHighlights() bool { return len(a.highlight) > 0 }

// parseCodeBlockInfo extracts title="…" and {n,m-p} attributes from a fenced
// code block info string. The language token (first whitespace-delimited word)
// is skipped. The function is tolerant: malformed expressions are silently
// dropped so a stray character never breaks the rest of the block. Pass "" if
// the input is empty.
//
// Recognised forms (after the language):
//
//	title="some text"      // title; quotes may be single or double
//	{3,5-8}                // highlight lines 3 and 5 through 8 inclusive
//	title="x.go" {2,4-6}   // both, in any order
//
// Whitespace separates attributes; the opening "{" must follow a space or
// be at the start of the info string so it cannot collide with braces
// embedded in titles.
func parseCodeBlockInfo(info string) codeBlockAttrs {
	var attrs codeBlockAttrs

	rest := strings.TrimSpace(info)
	if rest == "" {
		return attrs
	}

	// Strip the leading language token (first whitespace-delimited word)
	// if the info string contains whitespace. If the whole string is a
	// single token (no whitespace) and starts with an attribute like
	// title=… or {…}, treat the entire string as the attribute list —
	// there is no real language in that case.
	hasWS := strings.ContainsAny(rest, " 	")
	if hasWS {
		rest = stripLanguageToken(rest)
	}

	// Split on whitespace, but keep quoted title together. We scan by hand
	// rather than using strings.Fields because we need to preserve the quoted
	// string as a single token.
	tokens := tokenizeInfo(rest)
	for _, tok := range tokens {
		switch {
		case strings.HasPrefix(tok, "title="):
			v, ok := unquoteAttr(tok[len("title="):])
			if ok {
				attrs.title = v
			}
		case strings.HasPrefix(tok, "{"):
			attrs.highlight = append(attrs.highlight, parseHighlightRanges(tok)...)
		}
	}
	if n := len(attrs.highlight); n > 0 {
		attrs.highlightSet = make(map[int]bool, n)
		for _, ln := range attrs.highlight {
			attrs.highlightSet[ln] = true
		}
	}
	return attrs
}

// tokenizeInfo splits an info string on whitespace while keeping title="…"
// (with internal spaces) as a single token. Other than title=, no attribute
// is allowed to contain spaces, so plain strings.Fields would split the rest.

// stripLanguageToken removes the first whitespace-delimited token from s and
// returns what remains, with surrounding whitespace trimmed. If s has no
// whitespace, the entire input is treated as the "language" and the returned
// string is empty — the caller is expected to fall back to using the full
// input as attribute syntax in that case.
func stripLanguageToken(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '	' {
			return strings.TrimSpace(s[i+1:])
		}
	}
	return ""
}

// looksLikeSingleAttr reports whether the info string contains no
// whitespace and starts with a kwelea attribute (title= or {…}). In that
// case the whole string is treated as the attribute list, with no
// language.
func looksLikeSingleAttr(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if strings.ContainsAny(s, " 	") {
		return false
	}
	return strings.HasPrefix(s, "title=") || strings.HasPrefix(s, "{")
}

// joinFencedLines concatenates the raw bytes of every line in a fenced
// code block, preserving the newlines between them.
func joinFencedLines(fcb *goldmarkast.FencedCodeBlock, src []byte) []byte {
	var buf bytes.Buffer
	n := fcb.Lines().Len()
	for i := 0; i < n; i++ {
		seg := fcb.Lines().At(i)
		buf.Write(seg.Value(src))
	}
	return buf.Bytes()
}

func tokenizeInfo(s string) []string {
	var out []string
	i := 0
	for i < len(s) {
		// skip whitespace
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
		if i >= len(s) {
			break
		}
		// if this token starts with title=, read until the matching closing quote
		if strings.HasPrefix(s[i:], "title=") {
			j := i + len("title=")
			// find opening quote
			if j < len(s) && (s[j] == '"' || s[j] == '\'') {
				quote := s[j]
				j++
				start := j
				for j < len(s) && s[j] != quote {
					j++
				}
				out = append(out, s[i:start]) // includes opening quote + content
				if j < len(s) {
					out[len(out)-1] = s[i : j+1] // include closing quote
					j++
				}
				i = j
				continue
			}
		}
		// generic token: read until next whitespace
		start := i
		for i < len(s) && s[i] != ' ' && s[i] != '\t' {
			i++
		}
		out = append(out, s[start:i])
	}
	return out
}

// unquoteAttr strips a single or double pair of quotes from around s. If the
// surrounding quotes don't match (or are missing), the value is returned
// verbatim and ok is false so the caller can decide to drop the attribute.
func unquoteAttr(s string) (string, bool) {
	if len(s) < 2 {
		return "", false
	}
	first, last := s[0], s[len(s)-1]
	if (first == '"' || first == '\'') && first == last {
		return s[1 : len(s)-1], true
	}
	return "", false
}

// parseHighlightRanges parses a "{n,m-p,…}" expression and returns the
// 1-indexed line numbers it contains. Invalid numbers and empty ranges
// (e.g. {3-2}) are skipped; the rest of the expression still parses.
func parseHighlightRanges(s string) []int {
	var out []int
	if !strings.HasPrefix(s, "{") || !strings.HasSuffix(s, "}") {
		return out
	}
	inner := s[1 : len(s)-1]
	if inner == "" {
		return out
	}
	for _, part := range strings.Split(inner, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if idx := strings.Index(part, "-"); idx >= 0 {
			lo, errLo := strconv.Atoi(strings.TrimSpace(part[:idx]))
			hi, errHi := strconv.Atoi(strings.TrimSpace(part[idx+1:]))
			if errLo != nil || errHi != nil || lo < 1 || hi < lo {
				continue
			}
			for ln := lo; ln <= hi; ln++ {
				out = append(out, ln)
			}
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil || n < 1 {
			continue
		}
		out = append(out, n)
	}
	return out
}

// ----- AST node -----

// KindCodeBlock is the goldmark AST node kind for a fenced code block that
// carries kwelea-specific attributes (title and/or line highlights).
var KindCodeBlock = goldmarkast.NewNodeKind("CodeBlock")

// CodeBlockNode is an AST block node that replaces a FencedCodeBlock during
// the codeAttrsTransformer pass when the block's info string contains
// title="…" or {n,m-p} attributes. The source bytes are stored verbatim
// because Chroma tokenises from the raw source.
type CodeBlockNode struct {
	goldmarkast.BaseBlock
	Language     string       // "" when no language was specified
	Title        string       // "" when no title was specified
	Highlight    []int        // 1-indexed lines to highlight (nil/empty when none)
	HighlightSet map[int]bool // set form for O(1) membership
	Source       []byte       // raw source bytes (lines joined with '\n')
}

// Kind returns KindCodeBlock, satisfying the goldmark ast.Node interface.
func (n *CodeBlockNode) Kind() goldmarkast.NodeKind { return KindCodeBlock }

// Dump writes a debug representation of the node to standard output,
// satisfying the goldmark ast.Node interface.
func (n *CodeBlockNode) Dump(source []byte, level int) {
	goldmarkast.DumpHelper(n, source, level, map[string]string{
		"Language":  n.Language,
		"Title":     n.Title,
		"Highlight": strconv.Itoa(len(n.Highlight)),
	}, nil)
}

// ----- AST transformer -----

// codeAttrsTransformer walks the parsed AST and replaces every
// FencedCodeBlock that has at least one custom attribute with a
// CodeBlockNode. Plain blocks are left untouched so goldmark-highlighting
// still renders them with zero overhead.
type codeAttrsTransformer struct{}

// Transform walks the document once, collecting any FencedCodeBlock whose
// info string contains kwelea attributes (title="…" or {n,m-p}). The walk
// is read-only; the actual node replacement happens in a second pass so we
// do not mutate the tree while iterating it.
func (t *codeAttrsTransformer) Transform(doc *goldmarkast.Document, reader text.Reader, _ parser.Context) {
	src := reader.Source()

	type replacement struct{ old, newNode goldmarkast.Node }
	var replacements []replacement

	_ = goldmarkast.Walk(doc, func(n goldmarkast.Node, entering bool) (goldmarkast.WalkStatus, error) {
		if !entering || n.Kind() != goldmarkast.KindFencedCodeBlock {
			return goldmarkast.WalkContinue, nil
		}
		fcb := n.(*goldmarkast.FencedCodeBlock)

		// Read the raw info string (the part after the opening fence).
		var info []byte
		if fcb.Info != nil {
			info = fcb.Info.Segment.Value(src)
		}
		infoStr := string(info)

		// Goldmark treats the first whitespace-delimited token as the
		// language. If the entire info string is a single token (no
		// whitespace) and that token is itself a kwelea attribute, there
		// is no real language — the whole string is the attribute list.
		// In that case, parse it directly.
		if looksLikeSingleAttr(infoStr) {
			attrs := parseCodeBlockInfo(infoStr)
			if !attrs.hasTitle() && !attrs.hasHighlights() {
				return goldmarkast.WalkContinue, nil
			}
			replacements = append(replacements, replacement{
				old: fcb,
				newNode: &CodeBlockNode{
					Language:     "",
					Title:        attrs.title,
					Highlight:    attrs.highlight,
					HighlightSet: attrs.highlightSet,
					Source:       joinFencedLines(fcb, src),
				},
			})
			return goldmarkast.WalkContinue, nil
		}

		attrs := parseCodeBlockInfo(infoStr)
		if !attrs.hasTitle() && !attrs.hasHighlights() {
			return goldmarkast.WalkContinue, nil
		}

		// Collect the raw source lines.
		var srcBuf bytes.Buffer
		lineCount := fcb.Lines().Len()
		for i := 0; i < lineCount; i++ {
			line := fcb.Lines().At(i)
			srcBuf.Write(line.Value(src))
		}

		lang := ""
		if fcb.Language(src) != nil {
			lang = string(fcb.Language(src))
		}
		// If the goldmark-reported language is itself a kwelea attribute
		// (single-token info with no whitespace), it's not a real
		// language — clear it.
		if strings.ContainsAny(lang, "={") {
			lang = ""
		}

		replacements = append(replacements, replacement{
			old: fcb,
			newNode: &CodeBlockNode{
				Language:     lang,
				Title:        attrs.title,
				Highlight:    attrs.highlight,
				HighlightSet: attrs.highlightSet,
				Source:       srcBuf.Bytes(),
			},
		})
		return goldmarkast.WalkContinue, nil
	})

	for _, r := range replacements {
		if parent := r.old.Parent(); parent != nil {
			parent.ReplaceChild(parent, r.old, r.newNode)
		}
	}
}

// ----- node renderer -----

// codeAttrsRenderer renders CodeBlockNode entries. It uses Chroma directly
// (rather than delegating to goldmark-highlighting) so it has full control
// over the wrapper markup and can post-process the per-line classes to add
// highlight-line where requested.
type codeAttrsRenderer struct {
	lightStyle string
	darkStyle  string
}

// renderCodeBlock emits:
//
//	<figure class="code-block" data-lang="go">
//	  <div class="code-title">cmd/root.go</div>
//	  <div class="code-body"><Chroma HTML></div>
//	</figure>
//
// If there is no title, the code-title div is omitted and the body sits
// directly inside the figure so existing block styles keep applying.
func (r *codeAttrsRenderer) renderCodeBlock(w util.BufWriter, _ []byte, node goldmarkast.Node, entering bool) (goldmarkast.WalkStatus, error) {
	if !entering {
		return goldmarkast.WalkContinue, nil
	}
	n := node.(*CodeBlockNode)

	body, err := renderCodeBlockBody(n, r.lightStyle, r.darkStyle)
	if err != nil {
		// Fall back to a plain escaped <pre> so a Chroma failure never
		// breaks the whole page. The error is swallowed; the
		// non-decorated block still renders.
		body = fmt.Sprintf(`<pre class="code-fallback"><code>%s</code></pre>`, html.EscapeString(string(n.Source)))
	}

	langAttr := ""
	if n.Language != "" {
		langAttr = ` data-lang="` + html.EscapeString(n.Language) + `"`
	}

	_, _ = io_WriteString(w, `<figure class="code-block"`+langAttr+`>`)
	if n.Title != "" {
		_, _ = io_WriteString(w, `<div class="code-title">`+html.EscapeString(n.Title)+`</div>`)
	}
	_, _ = io_WriteString(w, `<div class="code-body">`)
	_, _ = io_WriteString(w, body)
	_, _ = io_WriteString(w, `</div></figure>`)
	return goldmarkast.WalkContinue, nil
}

// renderCodeBlockBody formats source through Chroma in line-numbered table
// mode, then post-processes the output to add the highlight-line class to
// the <span class="line"> rows matching n.Lines. The result is the inner
// HTML of the <figure class="code-block"> — i.e. the chroma <div class=
// "chroma">…</div> tree, with one class edit per highlighted line.
func renderCodeBlockBody(n *CodeBlockNode, lightStyleName, darkStyleName string) (string, error) {
	style := styles.Get(lightStyleName)
	if style == nil {
		style = styles.Fallback
	}

	var lexer chroma.Lexer
	if n.Language != "" {
		lexer = lexers.Get(n.Language)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// Build chroma HighlightLines ranges. Chroma's HighlightLines takes
	// 1-based, inclusive line numbers (see the function's godoc), so we
	// pass n.Highlight through unchanged. The set form n.HighlightSet
	// uses the same 1-based numbering.
	var hlRanges [][2]int
	for _, ln := range n.Highlight {
		hlRanges = append(hlRanges, [2]int{ln, ln})
	}

	opts := []chromahtml.Option{
		chromahtml.WithClasses(true),
		chromahtml.WithLineNumbers(true),
		chromahtml.LineNumbersInTable(true),
	}
	if len(hlRanges) > 0 {
		opts = append(opts, chromahtml.HighlightLines(hlRanges))
	}
	formatter := chromahtml.New(opts...)

	iter, err := lexer.Tokenise(nil, string(n.Source))
	if err != nil {
		return "", fmt.Errorf("tokenising code: %w", err)
	}
	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iter); err != nil {
		return "", fmt.Errorf("formatting chroma output: %w", err)
	}
	return injectHighlightClass(buf.String(), n.HighlightSet), nil
}

// injectHighlightClass rewrites the per-line <span class="line …"> tags
// emitted by Chroma's table-mode line-numbering formatter to add the
// project-level highlight-line class. It does this by tracking the current
// 1-indexed line number and, for every <span class="line …"> it sees,
// checking membership in the highlight set.
//
// We do the counting manually rather than relying on the line number
// rendered in the gutter: the gutter markup wraps each line number in its
// own <span> and is awkward to parse, while <span class="line"> appears
// exactly once per code line and is the canonical anchor for per-line
// decoration.
func injectHighlightClass(html string, lines map[int]bool) string {
	if len(lines) == 0 {
		return html
	}
	marker := []byte(`<span class="line`)
	idx := 0
	ln := 0
	var out bytes.Buffer
	out.Grow(len(html))
	for {
		next := bytes.Index([]byte(html[idx:]), marker)
		if next < 0 {
			out.WriteString(html[idx:])
			break
		}
		out.WriteString(html[idx : idx+next])
		// advance past <span class="line
		spanStart := idx + next
		spanEnd := spanStart + len(marker)
		// find the end of the opening tag
		close := bytes.IndexByte([]byte(html[spanEnd:]), '>')
		if close < 0 {
			out.WriteString(html[spanStart:])
			break
		}
		ln++
		openTag := html[spanStart : spanEnd+close+1]
		if lines[ln] {
			out.WriteString(rewriteLineOpenTag(openTag))
		} else {
			out.WriteString(openTag)
		}
		idx = spanEnd + close + 1
	}
	return out.String()
}

// rewriteLineOpenTag inserts " highlight-line" into the class attribute of
// a <span class="line …"> tag emitted by Chroma. The two forms we see
// from Chroma are:
//
//	<span class="line">             (bare)
//	<span class="line hl">          (chroma added "hl" for HighlightLines)
//
// The class value is whatever sits between the opening quote of class=
// and the closing `">` of the tag. We rebuild the tag so the rewritten
// version starts with "line highlight-line" plus any original classes.
func rewriteLineOpenTag(tag string) string {
	if !strings.HasPrefix(tag, `<span class="line`) {
		return tag
	}
	if strings.Contains(tag, "highlight-line") {
		return tag
	}
	if !strings.HasSuffix(tag, `">`) {
		return tag
	}
	// Strip the surrounding <span class="…"> wrapper and re-emit it with
	// highlight-line inserted as the second class. The interior is
	// exactly the class attribute value (no other attributes are
	// expected in the per-line span).
	interior := tag[len(`<span class="`):]
	// interior now looks like: line">, line hl">, or line highlight-line">.
	// Drop the trailing `">` to get the class value.
	if !strings.HasSuffix(interior, `">`) {
		return tag
	}
	classValue := interior[:len(interior)-2]
	// classValue now looks like: "line" or "line hl" or "line highlight-line".
	// We rebuild so highlight-line comes right after "line".
	if classValue == "line" {
		return `<span class="line highlight-line">`
	}
	if strings.HasPrefix(classValue, "line ") {
		return `<span class="line highlight-line ` + classValue[len("line "):] + `">`
	}
	// Unknown form — leave it alone rather than risk corruption.
	return tag
}

// ----- extension -----

// codeAttrsExtension bundles the AST transformer and the node renderer.
type codeAttrsExtension struct {
	lightStyle string
	darkStyle  string
}

// NewCodeAttrsExtension returns a goldmark.Extender that adds title and
// line-highlight support to fenced code blocks. The theme names match the
// ones used to build the chroma stylesheet (see ChromaCSS).
func NewCodeAttrsExtension(themeCfg config.ThemeConfig) goldmark.Extender {
	return &codeAttrsExtension{
		lightStyle: themeCfg.LightCodeTheme,
		darkStyle:  themeCfg.DarkCodeTheme,
	}
}

func (e *codeAttrsExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&codeAttrsTransformer{}, 200),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&codeAttrsRenderer{
				lightStyle: e.lightStyle,
				darkStyle:  e.darkStyle,
			}, 200),
		),
	)
}

// RegisterFuncs satisfies renderer.NodeRenderer.
func (r *codeAttrsRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCodeBlock, r.renderCodeBlock)
}

// io_WriteString is a tiny wrapper around w.WriteString that returns the
// error from the underlying BufWriter. The renderer only ignores errors
// anyway, but we keep the same shape as other goldmark renderers in the
// project.
func io_WriteString(w util.BufWriter, s string) (int, error) {
	return w.WriteString(s)
}
