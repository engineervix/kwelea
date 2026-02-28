package builder

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/engineervix/kwelea/internal/nav"
	"github.com/engineervix/kwelea/internal/parser"
)

// PageData is the template context passed to every page render.
type PageData struct {
	Site         *nav.Site
	Page         *nav.Page
	SectionLabel string // nav section containing this page (empty for root pages)
	DevMode      bool   // true only during `kwelea serve`
}

// Build runs the full Phase 4 pipeline:
//  1. Create output directory
//  2. Copy embedded assets (CSS, JS, vendor)
//  3. Generate and write chroma.css
//  4. Load Go templates from embedded FS
//  5. Parse all pages (fill Page.HTML and Page.TOC)
//  6. Render all pages to HTML files
func Build(site *nav.Site, embFS fs.FS, devMode bool) error {
	outputDir := site.BuildCfg.OutputDir
	docsDir := site.BuildCfg.DocsDir

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir %q: %w", outputDir, err)
	}

	if err := copyAssets(embFS, outputDir); err != nil {
		return fmt.Errorf("copying assets: %w", err)
	}

	chromaCSS, err := parser.ChromaCSS(site.ThemeCfg)
	if err != nil {
		return fmt.Errorf("generating chroma CSS: %w", err)
	}
	chromaPath := filepath.Join(outputDir, "assets", "chroma.css")
	if err := os.WriteFile(chromaPath, []byte(chromaCSS), 0o644); err != nil {
		return fmt.Errorf("writing chroma.css: %w", err)
	}

	tmpl, err := loadTemplates(embFS)
	if err != nil {
		return err
	}

	// Pass 1: parse every page's Markdown into HTML + TOC.
	for _, page := range site.Pages {
		src, err := os.ReadFile(filepath.Join(docsDir, page.FilePath))
		if err != nil {
			return fmt.Errorf("reading %s: %w", page.FilePath, err)
		}
		html, toc, err := parser.Parse(page.FilePath, src, site.ThemeCfg)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", page.FilePath, err)
		}
		page.HTML = html
		page.TOC = toc
	}

	// Pass 2: render every page with the full template.
	for _, page := range site.Pages {
		if err := renderPage(tmpl, site, page, outputDir, devMode); err != nil {
			return fmt.Errorf("rendering %s: %w", page.Path, err)
		}
	}

	return nil
}

// loadTemplates parses all embedded Go templates into one set.
// layout.html is the entry-point template; page.html and all partials
// register themselves via {{define "name"}} blocks.
func loadTemplates(embFS fs.FS) (*template.Template, error) {
	tmpl, err := template.ParseFS(embFS,
		"templates/layout.html",
		"templates/page.html",
		"templates/partials/sidebar.html",
		"templates/partials/toc.html",
		"templates/partials/pager.html",
		"templates/partials/dev-reload.html",
	)
	if err != nil {
		return nil, fmt.Errorf("loading templates: %w", err)
	}
	return tmpl, nil
}

// renderPage executes the layout template for a single page and writes the
// result to outputDir/<page-slug>/index.html.
func renderPage(tmpl *template.Template, site *nav.Site, page *nav.Page, outputDir string, devMode bool) error {
	outPath := pageOutputPath(outputDir, page.Path)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := PageData{
		Site:         site,
		Page:         page,
		SectionLabel: sectionFor(site, page.Path),
		DevMode:      devMode,
	}
	return tmpl.ExecuteTemplate(f, "layout.html", data)
}

// pageOutputPath maps a URL path (e.g. "/getting-started/") to a filesystem
// path (e.g. "site/getting-started/index.html").
func pageOutputPath(outputDir, urlPath string) string {
	clean := strings.Trim(urlPath, "/")
	if clean == "" {
		return filepath.Join(outputDir, "index.html")
	}
	return filepath.Join(outputDir, filepath.FromSlash(clean), "index.html")
}

// sectionFor returns the NavSection label for the given URL path, or "" if the
// page is in the root (unlabelled) section.
func sectionFor(site *nav.Site, pagePath string) string {
	for _, section := range site.Nav {
		for _, item := range section.Items {
			if item.Path == pagePath {
				return section.Label
			}
		}
	}
	return ""
}
