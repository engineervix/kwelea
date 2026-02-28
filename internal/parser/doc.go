// Package parser implements the Markdown processing pipeline for kwelea.
//
// Phase 3 work: goldmark pipeline, frontmatter extraction, ToC walker,
// admonitions extension (:::), and D2 diagram extension (```d2 fences).
//
// Dependencies (added in Phase 3):
//   - github.com/yuin/goldmark
//   - github.com/yuin/goldmark-highlighting/v2
//   - github.com/alecthomas/chroma/v2
//   - oss.terrastruct.com/d2
package parser
