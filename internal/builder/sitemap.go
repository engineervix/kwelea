package builder

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/engineervix/kwelea/internal/nav"
)

// sitemapURLSet is the top-level XML element for a sitemap.
// See https://www.sitemaps.org/schemas/sitemap/0.9
type sitemapURLSet struct {
	XMLName xml.Name      `xml:"urlset"`
	XMLNS   string        `xml:"xmlns,attr"`
	URLs    []sitemapURL  `xml:"url"`
}

// sitemapURL is a single <url> entry inside the sitemap.
type sitemapURL struct {
	Loc string `xml:"loc"`
}

// writeSitemap generates a standard sitemap.xml containing all non-draft pages
// and writes it to outputDir/sitemap.xml. It must be called after the site has
// been fully built (pages walked and nav resolved).
//
// If the site has no base_url configured the sitemap is skipped silently —
// without a canonical origin URL the <loc> entries would be meaningless.
func writeSitemap(site *nav.Site, outputDir string) error {
	if site.BaseURL == "" {
		return nil
	}

	base := strings.TrimRight(site.BaseURL, "/")

	urlSet := sitemapURLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]sitemapURL, 0, len(site.Pages)),
	}

	for _, page := range site.Pages {
		// Draft pages are already filtered out by WalkDocs, but guard anyway.
		if page.Draft {
			continue
		}
		loc := base + page.Path
		urlSet.URLs = append(urlSet.URLs, sitemapURL{Loc: loc})
	}

	data, err := xml.MarshalIndent(urlSet, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling sitemap: %w", err)
	}

	// xml.MarshalIndent does not append a trailing newline.
	data = append(data, '\n')

	// Prepend the XML declaration manually — encoding/xml does not add one
	// with MarshalIndent, and a leading <?xml ...?> is required by the spec.
	out := append([]byte(xml.Header), data...)

	outPath := filepath.Join(outputDir, "sitemap.xml")
	return os.WriteFile(outPath, out, 0o644)
}