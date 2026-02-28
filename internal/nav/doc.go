// Package nav implements filesystem walking and nav tree construction for kwelea.
//
// Phase 2 work: walk docs_dir recursively, strip numeric prefixes, parse
// frontmatter, build NavSection/NavItem trees, derive flat []Page list for
// prev/next linking.
package nav
