// Package builder implements the HTML renderer, search index writer, and asset
// copier for kwelea.
//
// The primary entry point is [Build], which performs a two-pass pipeline:
//
//  1. Parse pass — every non-draft Markdown page is processed by [parser.Parse].
//     The resulting HTML, table-of-contents, and optional H1 title are stored
//     on each [nav.Page].
//
//  2. Render pass — each page is executed through the html/template pipeline
//     (layout.html → page.html + partials) and written to
//     <output_dir>/<path>/index.html.
//
// After both passes, [Build] copies embedded CSS/JS/font assets to the output
// directory and writes search-index.json for client-side FlexSearch.
package builder
