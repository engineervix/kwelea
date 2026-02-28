package nav

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/engineervix/kwelea/internal/config"
)

// Build creates the navigation tree from the discovered pages.
// If navEntries is non-empty (manual [[nav]] sections from kwelea.toml) it
// uses those; otherwise it builds the nav automatically from the filesystem.
// An error is returned only when a manual nav ref is ambiguous.
func Build(pages []*Page, navEntries []config.NavEntry) ([]NavSection, error) {
	if len(navEntries) > 0 {
		return buildManualNav(pages, navEntries)
	}
	return buildAutoNav(pages), nil
}

// buildAutoNav groups pages by their first-level directory and creates one
// NavSection per group. Root-level files form a section with an empty label.
// Within each section, index.md-based pages are sorted to the front.
func buildAutoNav(pages []*Page) []NavSection {
	type group struct {
		dir   string
		label string
		items []*Page
	}

	var groups []group
	groupIdx := map[string]int{} // dir → index into groups slice

	for _, p := range pages {
		dir := firstDir(p.FilePath)
		idx, seen := groupIdx[dir]
		if !seen {
			label := ""
			if dir != "" {
				label = slugToTitle(stripNumericPrefix(dir))
			}
			idx = len(groups)
			groups = append(groups, group{dir: dir, label: label})
			groupIdx[dir] = idx
		}
		groups[idx].items = append(groups[idx].items, p)
	}

	sections := make([]NavSection, len(groups))
	for i, g := range groups {
		// Sort so that index.md-based pages always appear first in the section.
		sort.SliceStable(g.items, func(a, b int) bool {
			aIsIndex := filepath.Base(g.items[a].FilePath) == "index.md"
			bIsIndex := filepath.Base(g.items[b].FilePath) == "index.md"
			return aIsIndex && !bIsIndex
		})
		items := make([]NavItem, len(g.items))
		for j, p := range g.items {
			items[j] = NavItem{Title: p.Title, Path: p.Path}
		}
		sections[i] = NavSection{Label: g.label, Items: items}
	}

	return sections
}

// buildManualNav creates the nav from explicit [[nav]] sections in kwelea.toml.
func buildManualNav(pages []*Page, entries []config.NavEntry) ([]NavSection, error) {
	sections := make([]NavSection, 0, len(entries))
	for _, e := range entries {
		items := make([]NavItem, 0, len(e.Pages))
		for _, ref := range e.Pages {
			p, err := lookupPageByRef(pages, ref)
			if err != nil {
				return nil, err
			}
			if p == nil {
				continue // referenced page not found — skip gracefully
			}
			items = append(items, NavItem{Title: p.Title, Path: p.Path})
		}
		if len(items) > 0 {
			sections = append(sections, NavSection{Label: e.Section, Items: items})
		}
	}
	return sections, nil
}

// lookupPageByRef finds the Page that matches a nav reference string from a
// [[nav]] pages list in kwelea.toml.
//
// Matching is tried in priority order:
//
//  1. Exact FilePath match ("guide/01-configuration.md") — globally unique,
//     first match wins.
//  2. Base filename match ("01-getting-started.md") — first match wins.
//  3. Slug match ("getting-started") — numeric prefix and .md extension are
//     stripped from both the ref and each page's filename. If more than one
//     page produces the same slug, the build fails with a clear error naming
//     the conflicting files and instructing the author to use the full path.
func lookupPageByRef(pages []*Page, ref string) (*Page, error) {
	// 1. Exact FilePath match.
	for _, p := range pages {
		if p.FilePath == ref {
			return p, nil
		}
	}

	// 2. Base filename match.
	refBase := filepath.Base(ref)
	for _, p := range pages {
		if filepath.Base(p.FilePath) == refBase {
			return p, nil
		}
	}

	// 3. Slug match — strip numeric prefix and .md extension from both sides.
	refSlug := stripNumericPrefix(strings.TrimSuffix(refBase, ".md"))
	var matches []*Page
	for _, p := range pages {
		pageSlug := stripNumericPrefix(strings.TrimSuffix(filepath.Base(p.FilePath), ".md"))
		if pageSlug == refSlug {
			matches = append(matches, p)
		}
	}

	switch len(matches) {
	case 0:
		return nil, nil
	case 1:
		return matches[0], nil
	default:
		names := make([]string, len(matches))
		for i, m := range matches {
			names[i] = fmt.Sprintf("%q", m.FilePath)
		}
		return nil, fmt.Errorf(
			"nav ref %q is ambiguous — slug matches %s; use the full relative path instead",
			ref, strings.Join(names, ", "),
		)
	}
}

// FlattenNav returns all unique Pages from the nav sections in order.
// The result is used to assign Prev/Next links.
func FlattenNav(sections []NavSection, pages []*Page) []*Page {
	byPath := make(map[string]*Page, len(pages))
	for _, p := range pages {
		byPath[p.Path] = p
	}

	var flat []*Page
	seen := make(map[string]bool)
	for _, section := range sections {
		for _, item := range section.Items {
			if seen[item.Path] {
				continue
			}
			seen[item.Path] = true
			if p, ok := byPath[item.Path]; ok {
				flat = append(flat, p)
			}
		}
	}
	return flat
}

// LinkPrevNext sets Prev and Next on each page in the slice, which must
// already be in nav order (output of FlattenNav).
func LinkPrevNext(pages []*Page) {
	for i, p := range pages {
		if i > 0 {
			p.Prev = pages[i-1]
		}
		if i < len(pages)-1 {
			p.Next = pages[i+1]
		}
	}
}

// NewSite runs the full Phase 2 pipeline: walk docs_dir, build the nav tree,
// flatten into a page list, and link prev/next.
func NewSite(cfg *config.Config) (*Site, error) {
	pages, err := WalkDocs(cfg.Build.DocsDir)
	if err != nil {
		return nil, fmt.Errorf("scanning %q: %w", cfg.Build.DocsDir, err)
	}

	sections, err := Build(pages, cfg.Nav)
	if err != nil {
		return nil, fmt.Errorf("building nav: %w", err)
	}
	flat := FlattenNav(sections, pages)
	LinkPrevNext(flat)

	return &Site{
		Title:        cfg.Site.Title,
		Version:      cfg.Site.Version,
		BaseURL:      cfg.Site.BaseURL,
		Repo:         cfg.Site.Repo,
		RepoPlatform: repoplatform(cfg.Site.Repo),
		BuildCfg:     cfg.Build,
		ServeCfg:     cfg.Serve,
		ThemeCfg:     cfg.Theme,
		Nav:          sections,
		Pages:        flat,
	}, nil
}

// repoplatform returns "github", "gitlab", or "" for the given repo URL.
func repoplatform(repo string) string {
	switch {
	case strings.Contains(repo, "github.com"):
		return "github"
	case strings.Contains(repo, "gitlab.com"):
		return "gitlab"
	default:
		return ""
	}
}
