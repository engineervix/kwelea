package parser

import (
	"bytes"

	goldmarkast "github.com/yuin/goldmark/ast"

	"github.com/engineervix/kwelea/internal/nav"
)

// extractTOC walks the goldmark AST and collects all h2 and h3 headings as
// TocItem values. parser.WithAutoHeadingID() must be active so that each
// heading node carries an "id" attribute by the time this is called (between
// Parse and Render).
func extractTOC(doc goldmarkast.Node, src []byte) []nav.TocItem {
	var items []nav.TocItem
	goldmarkast.Walk(doc, func(n goldmarkast.Node, entering bool) (goldmarkast.WalkStatus, error) {
		if !entering || n.Kind() != goldmarkast.KindHeading {
			return goldmarkast.WalkContinue, nil
		}
		h := n.(*goldmarkast.Heading)
		if h.Level < 2 || h.Level > 3 {
			return goldmarkast.WalkContinue, nil
		}

		id := ""
		if idAttr, ok := h.AttributeString("id"); ok {
			if idBytes, ok := idAttr.([]byte); ok {
				id = string(idBytes)
			}
		}

		label := headingPlainText(h, src)
		if id != "" || label != "" {
			items = append(items, nav.TocItem{
				ID:    id,
				Text:  label,
				Level: h.Level,
			})
		}
		return goldmarkast.WalkContinue, nil
	})
	return items
}

// headingPlainText returns the plain-text content of a heading node by
// recursively stripping all inline formatting and collecting leaf text.
func headingPlainText(n goldmarkast.Node, src []byte) string {
	var buf bytes.Buffer
	collectInlineText(n, src, &buf)
	return buf.String()
}

func collectInlineText(n goldmarkast.Node, src []byte, buf *bytes.Buffer) {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		switch t := c.(type) {
		case *goldmarkast.Text:
			buf.Write(t.Segment.Value(src))
			if t.SoftLineBreak() {
				buf.WriteByte(' ')
			}
		case *goldmarkast.String:
			buf.Write(t.Value)
		default:
			collectInlineText(c, src, buf)
		}
	}
}
