package parser

import (
	"bytes"
	"html"
	"io"
	"strings"

	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// KindAdmonition is the goldmark AST node kind for ::: container blocks.
var KindAdmonition = goldmarkast.NewNodeKind("Admonition")

// AdmonitionNode is an AST block node wrapping a ::: container block.
// Its children are the parsed Markdown contents of the block body.
type AdmonitionNode struct {
	goldmarkast.BaseBlock
	AdmonitionType string // "info" | "tip" | "warning" | "danger" | "details"
	Title          string // for details: the <summary> text; otherwise empty
}

// NewAdmonitionNode allocates an AdmonitionNode with the given type and title.
func NewAdmonitionNode(admonType, title string) *AdmonitionNode {
	n := &AdmonitionNode{AdmonitionType: admonType, Title: title}
	n.SetBlankPreviousLines(true)
	return n
}

// Kind returns KindAdmonition, satisfying the goldmark ast.Node interface.
func (n *AdmonitionNode) Kind() goldmarkast.NodeKind { return KindAdmonition }

// Dump writes a debug representation of the node to standard output, satisfying
// the goldmark ast.Node interface.
func (n *AdmonitionNode) Dump(source []byte, level int) {
	goldmarkast.DumpHelper(n, source, level, map[string]string{
		"Type":  n.AdmonitionType,
		"Title": n.Title,
	}, nil)
}

// ----- block parser -----

var validAdmonitionTypes = map[string]bool{
	"info": true, "tip": true, "warning": true, "danger": true, "details": true,
}

type admonitionParser struct{}

func (p *admonitionParser) Trigger() []byte { return []byte{':'} }

func (p *admonitionParser) Open(parent goldmarkast.Node, reader text.Reader, pc parser.Context) (goldmarkast.Node, parser.State) {
	line, seg := reader.PeekLine()
	line = bytes.TrimRight(line, "\r\n")
	if !bytes.HasPrefix(line, []byte(":::")) {
		return nil, parser.NoChildren
	}
	rest := bytes.TrimSpace(line[3:])
	if len(rest) == 0 {
		return nil, parser.NoChildren // bare ::: is not an opener
	}

	var admonType, title string
	if idx := bytes.IndexByte(rest, ' '); idx >= 0 {
		admonType = strings.ToLower(string(rest[:idx]))
		title = string(bytes.TrimSpace(rest[idx+1:]))
	} else {
		admonType = strings.ToLower(string(rest))
	}
	if !validAdmonitionTypes[admonType] {
		return nil, parser.NoChildren
	}

	reader.Advance(seg.Len())
	return NewAdmonitionNode(admonType, title), parser.HasChildren
}

func (p *admonitionParser) Continue(node goldmarkast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, seg := reader.PeekLine()
	if bytes.Equal(bytes.TrimSpace(line), []byte(":::")) {
		reader.Advance(seg.Len())
		return parser.Close
	}
	return parser.Continue | parser.HasChildren
}

func (p *admonitionParser) Close(node goldmarkast.Node, reader text.Reader, pc parser.Context) {}

func (p *admonitionParser) CanInterruptParagraph() bool { return true }
func (p *admonitionParser) CanAcceptIndentedLine() bool { return false }

// ----- HTML renderer -----

var admonitionLabels = map[string]string{
	"info":    "Info",
	"tip":     "Tip",
	"warning": "Warning",
	"danger":  "Danger",
}

type admonitionHTMLRenderer struct{}

func (r *admonitionHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindAdmonition, r.renderAdmonition)
}

func (r *admonitionHTMLRenderer) renderAdmonition(w util.BufWriter, source []byte, node goldmarkast.Node, entering bool) (goldmarkast.WalkStatus, error) {
	n := node.(*AdmonitionNode)
	if entering {
		if n.AdmonitionType == "details" {
			title := n.Title
			if title == "" {
				title = "Details"
			}
			_, _ = io.WriteString(w, "<details class=\"admonition admonition-details\">\n")
			_, _ = io.WriteString(w, "<summary>"+html.EscapeString(title)+"</summary>\n")
		} else {
			label := admonitionLabels[n.AdmonitionType]
			_, _ = io.WriteString(w, "<div class=\"admonition admonition-"+n.AdmonitionType+"\">\n")
			_, _ = io.WriteString(w, "<div class=\"admonition-label\">"+label+"</div>\n")
			_, _ = io.WriteString(w, "<div class=\"admonition-body\">\n")
		}
	} else {
		if n.AdmonitionType == "details" {
			_, _ = io.WriteString(w, "</details>\n")
		} else {
			_, _ = io.WriteString(w, "</div>\n</div>\n")
		}
	}
	return goldmarkast.WalkContinue, nil
}

// ----- extension -----

type admonitionsExtension struct{}

// Admonitions is the goldmark.Extender that adds ::: container block support.
var Admonitions goldmark.Extender = &admonitionsExtension{}

func (e *admonitionsExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(&admonitionParser{}, 800),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&admonitionHTMLRenderer{}, 500),
		),
	)
}
