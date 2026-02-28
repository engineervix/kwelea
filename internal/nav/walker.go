package nav

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// WalkDocs discovers all non-draft .md files under docsDir and returns them
// as Page values with metadata populated. HTML and TOC are left empty until
// Phase 3 fills them in via the Markdown parser.
//
// Files are returned in the order filepath.WalkDir visits them — lexical
// order within each directory. The nav builder re-sorts within sections to
// put index.md pages first.
func WalkDocs(docsDir string) ([]*Page, error) {
	var pages []*Page

	err := filepath.WalkDir(docsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		fm, _, err := parseFrontmatter(string(data))
		if err != nil {
			return fmt.Errorf("parsing frontmatter in %s: %w", path, err)
		}
		if fm.Draft {
			return nil // silently skip draft pages
		}

		relPath, err := filepath.Rel(docsDir, path)
		if err != nil {
			return fmt.Errorf("resolving relative path for %s: %w", path, err)
		}

		title := fm.Title
		if title == "" {
			title = derivedTitle(relPath)
		}

		pages = append(pages, &Page{
			Title:                title,
			TitleFromFrontmatter: fm.Title != "",
			Description:          fm.Description,
			Path:                 fileToURLPath(relPath),
			FilePath:             relPath,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking docs directory %q: %w", docsDir, err)
	}
	return pages, nil
}
