// Package nav implements the filesystem walker, nav tree builder, and core
// data types used throughout kwelea.
//
// The primary entry point is [NewSite], which reads a [config.Config], walks
// the docs directory (or follows the manual [[nav]] order defined in
// kwelea.toml), parses YAML-style frontmatter from every Markdown file, strips
// numeric filename prefixes, derives URL slugs, and returns a fully populated
// [Site] whose pages are linked in prev/next order.
//
// Exported types — [Site], [NavSection], [NavItem], [Page], [TocItem] — are
// the data model that the builder and server packages depend on.
package nav
