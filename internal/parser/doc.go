// Package parser implements the Markdown processing pipeline for kwelea.
//
// The two exported entry points are:
//
//   - [Parse] converts a Markdown source file into rendered HTML, an in-page
//     table-of-contents, and the plain-text H1 title (if present). It runs the
//     full goldmark pipeline including syntax highlighting, admonitions, and D2
//     diagram rendering.
//
//   - [ChromaCSS] generates a combined Chroma stylesheet containing both a
//     light-mode and a dark-mode rule set, written once per build to
//     assets/chroma.css in the output directory.
//
// Markdown extensions provided by this package:
//
//   - Admonitions: ::: info / ::: tip / ::: warning / ::: danger / ::: details
//     blocks rendered as styled <div> or <details> elements.
//
//   - D2 diagrams: fenced code blocks with language "d2" are compiled to
//     inline SVG pairs (light + dark) using the D2 Go library.
package parser
