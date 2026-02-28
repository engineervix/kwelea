package parser

import (
	"bytes"
	"context"
	"fmt"
	"html"

	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"
)

// KindD2 is the goldmark AST node kind for rendered D2 diagrams.
var KindD2 = goldmarkast.NewNodeKind("D2Diagram")

// D2Node is an AST block node that holds pre-rendered SVG HTML for a D2 diagram.
// It replaces FencedCodeBlock nodes with language "d2" during the AST transform
// pass, before goldmark-highlighting processes fenced blocks.
type D2Node struct {
	goldmarkast.BaseBlock
	SVG []byte // rendered HTML: <div class="d2-diagram"><div class="d2-light">…</div><div class="d2-dark">…</div></div>
}

func (n *D2Node) Kind() goldmarkast.NodeKind { return KindD2 }

func (n *D2Node) Dump(source []byte, level int) {
	goldmarkast.DumpHelper(n, source, level, map[string]string{"SVG": "(omitted)"}, nil)
}

// ----- AST transformer -----

// d2Transformer walks the parsed AST and replaces every FencedCodeBlock with
// language "d2" with a D2Node containing the pre-rendered SVG. This happens
// before the renderer runs, so goldmark-highlighting never sees d2 blocks.
type d2Transformer struct{}

func (t *d2Transformer) Transform(doc *goldmarkast.Document, reader text.Reader, pc parser.Context) {
	src := reader.Source()
	type replacement struct{ old, new goldmarkast.Node }
	var replacements []replacement

	goldmarkast.Walk(doc, func(n goldmarkast.Node, entering bool) (goldmarkast.WalkStatus, error) {
		if !entering || n.Kind() != goldmarkast.KindFencedCodeBlock {
			return goldmarkast.WalkContinue, nil
		}
		fcb := n.(*goldmarkast.FencedCodeBlock)
		if string(fcb.Language(src)) != "d2" {
			return goldmarkast.WalkContinue, nil
		}

		var code bytes.Buffer
		for i := 0; i < fcb.Lines().Len(); i++ {
			line := fcb.Lines().At(i)
			code.Write(line.Value(src))
		}

		svgHTML, err := renderD2Diagram(code.String())
		if err != nil {
			svgHTML = []byte(fmt.Sprintf(
				"<div class=\"d2-diagram d2-error\"><pre>D2 render error: %s</pre></div>\n",
				html.EscapeString(err.Error()),
			))
		}
		replacements = append(replacements, replacement{n, &D2Node{SVG: svgHTML}})
		return goldmarkast.WalkContinue, nil
	})

	for _, r := range replacements {
		if parent := r.old.Parent(); parent != nil {
			parent.ReplaceChild(parent, r.old, r.new)
		}
	}
}

// renderD2Diagram compiles a D2 source string and renders it as a dual-theme
// SVG pair — light (Cool Classics, theme 4) and dark (Dark Flagship, theme 201)
// — wrapped in a single <div class="d2-diagram"> container. CSS in the built
// site toggles visibility via [data-theme] on <html>.
func renderD2Diagram(source string) ([]byte, error) {
	ctx := context.Background()

	// D2 requires a Ruler to measure text dimensions for node labels.
	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return nil, fmt.Errorf("creating D2 text ruler: %w", err)
	}

	lightThemeID := d2themescatalog.CoolClassics.ID
	darkThemeID := d2themescatalog.DarkFlagshipTerrastruct.ID
	// LayoutResolver must be provided — d2lib.Compile panics on nil resolver
	// after applyDefaults sets Layout to "dagre". We use dagre directly.
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return func(ctx context.Context, g *d2graph.Graph) error {
			return d2dagrelayout.Layout(ctx, g, nil)
		}, nil
	}
	compileOpts := &d2lib.CompileOptions{
		Ruler:          ruler,
		LayoutResolver: layoutResolver,
	}

	lightDiagram, _, err := d2lib.Compile(ctx, source, compileOpts, &d2svg.RenderOpts{
		ThemeID: &lightThemeID,
	})
	if err != nil {
		return nil, fmt.Errorf("compiling D2 diagram (light theme): %w", err)
	}
	lightSVG, err := d2svg.Render(lightDiagram, &d2svg.RenderOpts{
		ThemeID: &lightThemeID,
	})
	if err != nil {
		return nil, fmt.Errorf("rendering D2 SVG (light theme): %w", err)
	}

	darkDiagram, _, err := d2lib.Compile(ctx, source, compileOpts, &d2svg.RenderOpts{
		ThemeID: &darkThemeID,
	})
	if err != nil {
		return nil, fmt.Errorf("compiling D2 diagram (dark theme): %w", err)
	}
	darkSVG, err := d2svg.Render(darkDiagram, &d2svg.RenderOpts{
		ThemeID: &darkThemeID,
	})
	if err != nil {
		return nil, fmt.Errorf("rendering D2 SVG (dark theme): %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString("<div class=\"d2-diagram\">\n")
	buf.WriteString("<div class=\"d2-light\">")
	buf.Write(lightSVG)
	buf.WriteString("</div>\n")
	buf.WriteString("<div class=\"d2-dark\">")
	buf.Write(darkSVG)
	buf.WriteString("</div>\n")
	buf.WriteString("</div>\n")
	return buf.Bytes(), nil
}

// ----- node renderer -----

type d2NodeRenderer struct{}

func (r *d2NodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindD2, r.renderD2Node)
}

func (r *d2NodeRenderer) renderD2Node(w util.BufWriter, _ []byte, node goldmarkast.Node, entering bool) (goldmarkast.WalkStatus, error) {
	if !entering {
		return goldmarkast.WalkContinue, nil
	}
	_, _ = w.Write(node.(*D2Node).SVG)
	return goldmarkast.WalkContinue, nil
}

// ----- extension -----

type d2Extension struct{}

// NewD2Extension returns a goldmark.Extender that intercepts ```d2 fenced
// code blocks at AST-transform time and replaces them with pre-rendered
// dual-theme inline SVG.
func NewD2Extension() goldmark.Extender {
	return &d2Extension{}
}

func (e *d2Extension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&d2Transformer{}, 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&d2NodeRenderer{}, 100),
		),
	)
}
