package nav

import (
	"html/template"

	"github.com/engineervix/kwelea/internal/config"
)

// Site is the root context passed to every template render. It is built once
// per build and then treated as read-only.
type Site struct {
	Title        string
	Version      string
	BaseURL      string
	BasePath     string // URL path prefix derived from BaseURL, e.g. "/kwelea" or ""
	Repo         string
	RepoPlatform string        // "github", "gitlab", or "" for unknown/not set
	ExtraHead    template.HTML // verbatim HTML injected into <head>; from kwelea.toml [site].extra_head
	ExtraFooter  template.HTML // verbatim HTML injected into <footer>; from kwelea.toml [site].extra_footer
	BuildCfg     config.BuildConfig
	ServeCfg     config.ServeConfig
	ThemeCfg     config.ThemeConfig
	Nav          []NavSection // full nav tree, built by this package
	Pages        []*Page      // flat ordered list for prev/next linking
}

// NavSection is a labelled group of pages in the sidebar.
// Label is empty ("") for root-level pages that have no section heading.
type NavSection struct {
	Label string
	Items []NavItem
}

// NavItem is a single entry in the sidebar nav.
// Active is set by the renderer when this item matches the current page.
type NavItem struct {
	Title  string
	Path   string
	Active bool
}

// Page holds all data needed to render a single documentation page.
// HTML and TOC are zero until Phase 3 fills them in via the parser.
type Page struct {
	Title                string
	TitleFromFrontmatter bool // true when Title came from frontmatter; false when derived from filename
	Description          string
	Path                 string        // canonical URL path, e.g. /getting-started/
	FilePath             string        // source .md file, relative to docs_dir
	HTML                 template.HTML // filled by parser (Phase 3)
	TOC                  []TocItem     // filled by parser (Phase 3)
	Prev                 *Page
	Next                 *Page
	Draft                bool
}

// TocItem is a heading entry in the in-page table of contents.
type TocItem struct {
	ID    string
	Text  string
	Level int // 2 (h2) or 3 (h3)
}
